package chain

import (
	goctx "context"

	"gitee.com/lance4117/GoFuse/logger"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"golang.org/x/net/context"
)

const (
	DefaultAddressPrefix = "cosmos"
	DefaultAddress       = "http://localhost:26657"
	DefaultKeyring       = "test"
)

var (
	// DefaultOptions are the default options for creating a Cosmos client
	DefaultOptions = []cosmosclient.Option{
		cosmosclient.WithAddressPrefix(DefaultAddressPrefix),
		cosmosclient.WithNodeAddress(DefaultAddress),
		cosmosclient.WithKeyringBackend(DefaultKeyring),
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
	ctx := goctx.Background()
	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, option...)

	if err != nil {
		logger.Fatal(err, "Init Cosmos Client Fail")
		return nil
	}
	acc, err := client.Account(address)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Init Blog Client %s success", address)

	return &Client{&client, address, &acc}
}

// DoBroadcastTx 广播包含给定消息的交易
// ctx: 操作的上下文
// msgs: 包含在交易中的消息
// returns: 广播交易的响应和发生的任何错误
func (c *Client) DoBroadcastTx(ctx context.Context, msgs ...sdktypes.Msg) (cosmosclient.Response, error) {
	return c.CosmosClient.BroadcastTx(ctx, *c.Account, msgs...)
}

// DoBroadcastTxWithOptions 使用自定义选项广播交易
// ctx: 操作的上下文
// options: 自定义交易选项
// msgs: 包含在交易中的消息
// returns: 广播交易的响应和发生的任何错误
func (c *Client) DoBroadcastTxWithOptions(ctx context.Context, options cosmosclient.TxOptions, msgs ...sdktypes.Msg) (cosmosclient.Response, error) {
	// Create a transaction with the given options
	tx, err := c.CosmosClient.CreateTxWithOptions(ctx, *c.Account,
		options, msgs...)
	if err != nil {
		return cosmosclient.Response{}, err
	}

	// Broadcast the transaction
	return tx.Broadcast(ctx)
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
