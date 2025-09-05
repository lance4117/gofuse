package fschainclient

import (
	"context"

	"gitee.com/lance4117/GoFuse/fslogger"
	"gitee.com/lance4117/GoFuse/fsonce"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

const (
	addressPrefix  = "cosmos"
	DefaultAddress = "http://localhost:26657"
)

var (
	DefaultOptions = []cosmosclient.Option{
		cosmosclient.WithAddressPrefix(addressPrefix),
		cosmosclient.WithNodeAddress(DefaultAddress),
		cosmosclient.WithKeyringBackend("test"),
	}
)

type ClientAdapter struct {
	Client *cosmosclient.Client
}

// InitClient 获取cosmos区块链客户端
var InitClient = fsonce.DoWithParam(func(option []cosmosclient.Option) *ClientAdapter {
	ctx := context.Background()
	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, option...)
	if err != nil {
		fslogger.Fatal(err, "Init Cosmos Client Fail")
		return nil
	}
	return &ClientAdapter{&client}
})
