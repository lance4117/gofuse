package fschain

import (
	goctx "context"
	"strings"

	"gitee.com/lance4117/GoFuse/fslogger"
	"gitee.com/lance4117/GoFuse/fsonce"
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
var InitClient = fsonce.DoWithParam(func(option []cosmosclient.Option) *Client {
	ctx := goctx.Background()
	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, option...)
	if err != nil {
		fslogger.Fatal(err, "Init Cosmos Client Fail")
		return nil
	}
	return &Client{&client}
})

func (c *Client) AccountAndAddress(address string) (*cosmosaccount.Account, string, error) {
	account, err := c.CosmosClient.Account(address)
	if err != nil {
		return nil, "", err
	}
	var prefix string
	pos := strings.Index(address, "1")
	if pos <= 0 {
		prefix = ""
	} else {
		prefix = address[:pos]
	}
	addr, err := account.Address(prefix)
	return &account, addr, err
}

func (c *Client) BroadcastTx(ctx goctx.Context, account *cosmosaccount.Account, msgs ...sdktypes.Msg) (cosmosclient.Response, error) {
	return c.CosmosClient.BroadcastTx(ctx, *account, msgs...)
}
