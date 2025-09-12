// Package chain provides a client for interacting with Cosmos blockchain
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
	// DefaultAddressPrefix is the default address prefix for Cosmos addresses
	DefaultAddressPrefix = "cosmos"
	// DefaultAddress is the default node address for the Cosmos client
	DefaultAddress = "http://localhost:26657"
	// DefaultKeyring is the default keyring backend for the Cosmos client
	DefaultKeyring = "test"
)

var (
	// DefaultOptions are the default options for creating a Cosmos client
	DefaultOptions = []cosmosclient.Option{
		cosmosclient.WithAddressPrefix(DefaultAddressPrefix),
		cosmosclient.WithNodeAddress(DefaultAddress),
		cosmosclient.WithKeyringBackend(DefaultKeyring),
	}
)

// Client is a wrapper around the Cosmos client with account information
type Client struct {
	CosmosClient *cosmosclient.Client
	Address      string
	Account      *cosmosaccount.Account
}

// InitClient 获取cosmos区块链客户端
// address: the address of the account to use
// option: client options for connecting to the Cosmos node
// returns: a new Client instance
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

// DoBroadcastTx broadcasts a transaction with the given messages
// ctx: the context for the operation
// msgs: the messages to include in the transaction
// returns: the response from broadcasting the transaction and any error that occurred
func (c *Client) DoBroadcastTx(ctx context.Context, msgs ...sdktypes.Msg) (cosmosclient.Response, error) {
	return c.CosmosClient.BroadcastTx(ctx, *c.Account, msgs...)
}

// DoBroadcastTxWithOptions broadcasts a transaction with custom options
// ctx: the context for the operation
// options: custom transaction options
// msgs: the messages to include in the transaction
// returns: the response from broadcasting the transaction and any error that occurred
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

func (c *Client) BankBalance(ctx context.Context, pagination *query.PageRequest) (sdk.Coins, error) {
	return c.CosmosClient.BankBalances(ctx, c.Address, pagination)
}

func (c *Client) Status(ctx context.Context) (*ctypes.ResultStatus, error) {
	return c.CosmosClient.Status(ctx)
}
