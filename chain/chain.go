package chain

import (
	"context"

	"gitee.com/lance4117/GoFuse/logger"
	"gitee.com/lance4117/GoFuse/once"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	signtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

const (
	DefaultAddressPrefix = "cosmos"
	DefaultAddress       = "http://localhost:26657"
)

var (
	// DefaultOptions are the default options for creating a Cosmos client
	DefaultOptions = []cosmosclient.Option{
		cosmosclient.WithAddressPrefix(DefaultAddressPrefix),
		cosmosclient.WithNodeAddress(DefaultAddress),
	}
)

// Client 基于账号的Cosmos客户端
type Client struct {
	CosmosClient *cosmosclient.Client
	Address      string
	Account      *cosmosaccount.Account
}

// InitClient 初始化Cosmos区块链客户端
// address: 要使用的账户地址
// option: 连接到Cosmos节点的客户端选项
// returns: 新的Client实例
func InitClient(address string, option []cosmosclient.Option) *Client {
	ctx := context.Background()
	// Create a Cosmos client instance
	getClient := once.DoWithErr(func() (cosmosclient.Client, error) {

		return cosmosclient.New(ctx, option...)
	})
	c, err := getClient()
	if err != nil {
		logger.Fatal(err, "Init Cosmos Client Fail")
		return nil
	}
	acc, err := c.Account(address)
	if err != nil {
		logger.Fatal(err, address)
	}
	logger.Infof("Init Blog Client %s success", address)

	return &Client{&c, address, &acc}
}

// DoBroadcastTx 广播包含给定消息的交易
// ctx: 操作的上下文
// msgs: 包含在交易中的消息
// returns: 广播交易的响应和发生的任何错误
func (c *Client) DoBroadcastTx(ctx context.Context, msgs ...sdk.Msg) (cosmosclient.Response, error) {
	return c.CosmosClient.BroadcastTx(ctx, *c.Account, msgs...)
}

// DoBroadcastTxWithOptions 使用自定义选项广播交易
// ctx: 操作的上下文
// options: 自定义交易选项
// msgs: 包含在交易中的消息
// returns: 广播交易的响应和发生的任何错误
func (c *Client) DoBroadcastTxWithOptions(ctx context.Context, options cosmosclient.TxOptions, msgs []sdk.Msg) (cosmosclient.Response, error) {
	// Create a transaction with the given options
	txService, err := c.CosmosClient.CreateTxWithOptions(ctx, *c.Account,
		options, msgs...)
	if err != nil {
		return cosmosclient.Response{}, err
	}

	// Broadcast the transaction
	return txService.Broadcast(ctx)
}

// DoBroadcastTxAsyncWithOptions 通过sync模式广播包含给定消息的交易
func (c *Client) DoBroadcastTxAsyncWithOptions(options cosmosclient.TxOptions, msgs []sdk.Msg) (*sdk.TxResponse, error) {
	cli := c.CosmosClient
	clientCtx := cli.Context()
	factory, err := createTxFactory(clientCtx, c.Account.Name)
	if err != nil {
		return nil, err
	}
	// 设置费用
	factory = factory.WithFees("0stake")

	// 构造未签名交易信息
	builder, err := factory.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	// 设置 gas 限制
	builder.SetGasLimit(options.GasLimit)

	// 手动签名
	signature, err := c.DoSign(builder)
	if err != nil {
		return nil, err
	}
	err = builder.SetSignatures(signature)
	if err != nil {
		return nil, err
	}

	// 构造tx的比特流
	txBytes, err := clientCtx.TxConfig.TxEncoder()(builder.GetTx())
	if err != nil {
		return nil, err
	}

	// 将tx以async的方式广播
	return clientCtx.BroadcastTxAsync(txBytes)
}

func (c *Client) DoSign(tx client.TxBuilder) (signtypes.SignatureV2, error) {
	account := c.Account
	cli := c.CosmosClient
	factory := cli.TxFactory
	ctx := cli.Context()
	signMode := factory.SignMode()
	sequence := factory.Sequence()
	// 首先需要获取签名者的公钥和序列号
	pk, err := account.Record.GetPubKey()
	if err != nil {
		return signtypes.SignatureV2{}, err
	}
	address, err := account.Record.GetAddress()
	if err != nil {
		return signtypes.SignatureV2{}, err
	}

	// 2. 生成需要签名的字节
	signerData := signing.SignerData{
		ChainID:       cli.Context().ChainID,
		AccountNumber: factory.AccountNumber(),
		Sequence:      sequence,
		PubKey:        pk,
		Address:       address.String(),
	}

	// 3. 获取需要签名的字节
	signBytes, err := signing.GetSignBytesAdapter(
		context.Background(),
		ctx.TxConfig.SignModeHandler(),
		signMode,
		signerData,
		tx.GetTx(),
	)
	if err != nil {
		return signtypes.SignatureV2{}, err
	}

	signByAddress, _, err := ctx.Keyring.SignByAddress(address, signBytes, signMode)
	if err != nil {
		return signtypes.SignatureV2{}, err
	}

	// 5. 构建SignatureData
	sigData := signtypes.SingleSignatureData{
		SignMode: signMode, Signature: signByAddress,
	}
	// 返回SignatureV2对象
	return signtypes.SignatureV2{
		PubKey:   pk,
		Data:     &sigData,
		Sequence: sequence,
	}, nil
}

// BankBalance 查询客户端账户的余额
// ctx: 操作的上下文
// pagination: 查询的分页参数
// returns: 账户中的代币和发生的任何错误
func (c *Client) BankBalance(ctx context.Context, pagination *query.PageRequest) (sdk.Coins, error) {
	return c.CosmosClient.BankBalances(ctx, c.Address, pagination)
}

// Status 查询Cosmos节点的状态
// ctx: 操作的上下文
// returns: 状态结果和发生的任何错误
func (c *Client) Status(ctx context.Context) (*ctypes.ResultStatus, error) {
	return c.CosmosClient.Status(ctx)
}

// 完整的工作示例
func createTxFactory(clientCtx client.Context, from string) (tx.Factory, error) {
	// 确保设置了 AccountRetriever
	if clientCtx.AccountRetriever == nil {
		clientCtx = clientCtx.WithAccountRetriever(authtypes.AccountRetriever{})
	}

	// 确保设置了 FromAddress
	if clientCtx.FromAddress.Empty() {
		fromAddr, fromName, _, err := client.GetFromFields(clientCtx, clientCtx.Keyring, from)
		if err != nil {
			return tx.Factory{}, err
		}
		clientCtx = clientCtx.WithFrom(from).
			WithFromAddress(fromAddr).
			WithFromName(fromName)
	}

	// 创建 Factory
	factory := tx.Factory{}.
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		WithTxConfig(clientCtx.TxConfig)

	// Prepare 方法会自动填充 accountNumber 和 sequence
	preparedFactory, err := factory.Prepare(clientCtx)
	if err != nil {
		return tx.Factory{}, err
	}

	return preparedFactory, nil
}
