package chain

import (
	goctx "context"

	"gitee.com/lance4117/GoFuse/logger"
	"gitee.com/lance4117/GoFuse/once"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

const (
	DefaultAddressPrefix = "cosmos"
	DefaultAddress       = "http://localhost:26657"
	DefaultKeyring       = "test"
)

var (
	DefaultOptions = []cosmosclient.Option{
		cosmosclient.WithAddressPrefix(DefaultAddressPrefix),
		cosmosclient.WithNodeAddress(DefaultAddress),
		cosmosclient.WithKeyringBackend(DefaultKeyring),
	}
)

type Client struct {
	CosmosClient *cosmosclient.Client
}

// InitClient 获取cosmos区块链客户端
var InitClient = once.DoWithParam(func(option []cosmosclient.Option) *Client {
	ctx := goctx.Background()
	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, option...)
	if err != nil {
		logger.Fatal(err, "Init Cosmos Client Fail")
		return nil
	}
	return &Client{&client}
})

func (c *Client) Account(nameOrAddress string) (*cosmosaccount.Account, error) {
	account, err := c.CosmosClient.Account(nameOrAddress)
	return &account, err
}

func (c *Client) Address(name string) (string, error) {
	return c.CosmosClient.Address(name)
}

func (c *Client) BroadcastTx(ctx goctx.Context, account *cosmosaccount.Account, msgs ...sdktypes.Msg) (cosmosclient.Response, error) {
	return c.CosmosClient.BroadcastTx(ctx, *account, msgs...)
}
