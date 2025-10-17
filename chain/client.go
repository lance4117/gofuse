package chain

import (
	"context"
	"os"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

// Client 区块链客户端
type Client struct {
	// 资源句柄
	grpcConn  *grpc.ClientConn
	rpcClient *rpchttp.HTTP

	// 核心对象
	Encoding  moduletestutil.TestEncodingConfig
	ClientCtx client.Context
	KeyRing   keyring.Keyring
	Config    Config
}

// New 创建新的区块链客户端实例（多实例安全）
func New(cfg Config) (*Client, error) {
	// 如果是空的，默认注册bank和auth模块
	if cfg.RegisterModules == nil {
		cfg.RegisterModules = func() moduletestutil.TestEncodingConfig {
			logger.Warn(errs.WarnRegisterModulesEmpty)
			return moduletestutil.MakeTestEncodingConfig(
				bank.AppModule{},
				auth.AppModule{},
			)
		}
	}
	encoding := cfg.RegisterModules()

	// 手动注册账户接口（很多查询/解码依赖）
	authtypes.RegisterInterfaces(encoding.InterfaceRegistry)

	// keyring
	kr, err := keyring.New(
		sdk.KeyringServiceName(),
		cfg.KeyringBackend,
		cfg.HomeDir,
		os.Stdin,
		encoding.Codec,
	)
	if err != nil {
		return nil, err
	}

	// gRPC 连接（单连接、可复用）
	grpcConn, err := grpc.NewClient(
		cfg.NodeGRPC,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(cfg.GRPCMaxRecvBytes),
			grpc.MaxCallSendMsgSize(cfg.GRPCMaxSendBytes),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                cfg.KeepaliveTime,
			Timeout:             cfg.KeepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, err
	}

	// Tendermint RPC（用于高度/状态/ABCI 查询等；广播建议走 gRPC TxService）
	rpcClient, err := rpchttp.New(cfg.NodeRPC, "/websocket")
	if err != nil {
		_ = grpcConn.Close()
		return nil, err
	}

	// 构建 client.Context
	cctx := client.Context{}.
		WithCodec(encoding.Codec).
		WithInterfaceRegistry(encoding.InterfaceRegistry).
		WithTxConfig(encoding.TxConfig).
		WithLegacyAmino(encoding.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync). // gRPC TxService 下保留无妨
		WithHomeDir(cfg.HomeDir).
		WithChainID(cfg.ChainID).
		WithKeyring(kr).
		WithNodeURI(cfg.NodeRPC).
		WithGRPCClient(grpcConn).
		WithSkipConfirmation(cfg.SkipConfirmation).
		WithClient(rpcClient)

	c := &Client{
		grpcConn:  grpcConn,
		rpcClient: rpcClient,
		Encoding:  encoding,
		ClientCtx: cctx,
		KeyRing:   kr,
		Config:    cfg,
	}

	// 健康检查
	if err = c.Ping(5 * time.Second); err != nil {
		return nil, err
	}

	return c, nil
}

// Address 通过 keyring 名称获取地址
func (c *Client) Address(name string) (sdk.AccAddress, error) {
	key, err := c.KeyRing.Key(name)
	if err != nil {
		return nil, err
	}
	return key.GetAddress()
}

// AccountNumberSequence 通过地址查询账户号/序列号（支持各种账户类型）
func (c *Client) AccountNumberSequence(addr sdk.AccAddress) (uint64, uint64, error) {
	return c.ClientCtx.AccountRetriever.GetAccountNumberSequence(c.ClientCtx, addr)
}

func (c *Client) Account(address sdk.AccAddress) (client.Account, error) {
	return c.ClientCtx.AccountRetriever.GetAccount(c.ClientCtx, address)
}

// TxFactory 返回一个带链配置的 *tx.Factory，便于构建/签名/广播
func (c *Client) TxFactory() (sdktx.Factory, error) {
	// 解析手续费
	fees, err := sdk.ParseCoinsNormalized(c.Config.Fee)
	if err != nil {
		return sdktx.Factory{}, err
	}
	// 创建交易工厂
	txf := sdktx.Factory{}.
		WithChainID(c.Config.ChainID).
		WithKeybase(c.KeyRing).
		WithTxConfig(c.ClientCtx.TxConfig).
		WithGas(c.Config.GasLimit).
		WithFees(fees.String()).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
	return txf, nil
}

// TxBuilder 构建 client.TxBuilder 并设置消息
func (c *Client) TxBuilder(msgs ...sdk.Msg) (client.TxBuilder, error) {
	txBuilder := c.ClientCtx.TxConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msgs...); err != nil {
		return nil, err
	}
	return txBuilder, nil
}

// SimulateGas 使用 gRPC TxService.Simulate 估算 gas（返回 GasUsed）
func (c *Client) SimulateGas(ctx context.Context, txBuilder client.TxBuilder) (uint64, error) {
	enc := c.ClientCtx.TxConfig.TxEncoder()
	txBytes, err := enc(txBuilder.GetTx())
	if err != nil {
		return 0, err
	}
	svc := txtypes.NewServiceClient(c.grpcConn)
	rsp, err := svc.Simulate(ctx, &txtypes.SimulateRequest{TxBytes: txBytes})
	if err != nil {
		return 0, err
	}
	return rsp.GasInfo.GetGasUsed(), nil
}

// SignTxWith 显式给定 accNum/seq 的签名（并发压测推荐用）
func (c *Client) SignTxWith(signerName string, txBuilder client.TxBuilder, accNum, seq uint64) error {
	factory, err := c.TxFactory()
	if err != nil {
		return err
	}
	txf := factory.
		WithAccountNumber(accNum).
		WithSequence(seq)
	return sdktx.Sign(context.Background(), txf, signerName, txBuilder, true)
}

// EncodeTxBytes 编码为 protobuf bytes
func (c *Client) EncodeTxBytes(txBuilder client.TxBuilder) ([]byte, error) {
	return c.ClientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
}

// BroadcastTx 通过 gRPC 广播（支持 ASYNC/SYNC/BLOCK）
func (c *Client) BroadcastTx(ctx context.Context, txBytes []byte, mode txtypes.BroadcastMode) (*txtypes.BroadcastTxResponse, error) {
	svc := txtypes.NewServiceClient(c.grpcConn)
	return svc.BroadcastTx(ctx, &txtypes.BroadcastTxRequest{
		Mode:    mode,
		TxBytes: txBytes,
	})
}

// SignAndBroadcast 一步到位：构建->设置->签名->广播
func (c *Client) SignAndBroadcast(
	ctx context.Context,
	signerName string,
	msgs ...sdk.Msg,
) (*txtypes.BroadcastTxResponse, error) {

	// 1) 构建
	txBuilder, err := c.TxBuilder(msgs...)
	if err != nil {
		return nil, err
	}
	txBuilder.SetGasLimit(c.Config.GasLimit)

	// 2) 签名（自动查询 accNum/seq）
	address, err := c.Address(signerName)
	if err != nil {
		return nil, err
	}
	num, seq, err := c.AccountNumberSequence(address)
	if err != nil {
		return nil, err
	}
	err = c.SignTxWith(signerName, txBuilder, num, seq)
	if err != nil {
		return nil, err
	}

	// 3) 编码 & 广播
	txBytes, err := c.EncodeTxBytes(txBuilder)
	if err != nil {
		return nil, err
	}
	return c.BroadcastTx(ctx, txBytes, c.Config.BroadcastMode)
}

// WaitForTx 轮询等待上链（优先 gRPC GetTx）
func (c *Client) WaitForTx(ctx context.Context, txHash string, pollInterval time.Duration) (*txtypes.GetTxResponse, error) {
	svc := txtypes.NewServiceClient(c.grpcConn)
	t := time.NewTicker(pollInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		case <-t.C:
			rsp, err := svc.GetTx(ctx, &txtypes.GetTxRequest{Hash: txHash})
			if err == nil && rsp != nil && rsp.TxResponse != nil {
				return rsp, nil
			}
			// 未找到通常返回 codes.NotFound，忽略并继续轮询
			if st, ok := status.FromError(err); ok && st.Message() != "" {
				// 其他严重错误提前返回
				// 你也可以根据 st.Code() == codes.NotFound 细分
			}
		}
	}
}

// Balance 查询某地址某 denom 的余额
func (c *Client) Balance(ctx context.Context, addr sdk.AccAddress, denom string) (sdk.Coin, error) {
	q := banktypes.NewQueryClient(c.grpcConn)
	rsp, err := q.Balance(ctx, &banktypes.QueryBalanceRequest{
		Address: addr.String(),
		Denom:   denom,
	})
	if err != nil {
		return sdk.Coin{}, err
	}
	if rsp.Balance == nil {
		return sdk.Coin{}, errs.ErrNoBalance(addr.String(), denom)
	}
	return *rsp.Balance, nil
}

// SendCoins 使用指定 signer 从 from -> to 转账
func (c *Client) SendCoins(
	ctx context.Context,
	signerName string,
	fromAddr, toAddr sdk.AccAddress,
	amount sdk.Coins,
) (*txtypes.BroadcastTxResponse, error) {
	msg := &banktypes.MsgSend{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      amount,
	}
	return c.SignAndBroadcast(ctx, signerName, msg)
}

// LatestHeight 获取最新区块高度
func (c *Client) LatestHeight(ctx context.Context) (int64, error) {
	var s *coretypes.ResultStatus
	var err error
	s, err = c.rpcClient.Status(ctx)
	if err != nil {
		return 0, err
	}
	return s.SyncInfo.LatestBlockHeight, nil
}

// Status 区块链信息
func (c *Client) Status() (*coretypes.ResultStatus, error) {
	return c.ClientCtx.Client.Status(context.Background())
}

// Close 优雅关闭
func (c *Client) Close() error {
	if c == nil {
		return nil
	}

	var firstErr error
	if c.grpcConn != nil {
		if err := c.grpcConn.Close(); err != nil {
			firstErr = err
		}
	}
	// rpchttp.HTTP 没有必须关闭的释放动作，这里忽略
	return firstErr
}

// Ping 查询节点状态（确认 RPC/gRPC 通）
func (c *Client) Ping(timeout time.Duration) error {
	tctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// 用 RPC 拉个 status（高度/同步信息）
	_, err := c.rpcClient.Status(tctx)
	if err != nil {
		return err
	}
	// gRPC 测试
	state := c.grpcConn.GetState()
	if state == connectivity.TransientFailure || state == connectivity.Shutdown {
		return errs.ErrGrpcConnFail
	}
	return nil
}
