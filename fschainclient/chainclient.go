package fschainclient

import (
	"context"

	"gitee.com/lance4117/GoFuse/fslogger"
	"gitee.com/lance4117/GoFuse/fsonce"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

const (
	addressPrefix      = "cosmos"
	DefaultAddress     = "http://localhost:26657"
	DefaultAccountName = "alice"
)

type ClientAdapter struct {
	Client *cosmosclient.Client
}

// InitCosmosClient 获取Cosmos前缀的cosmos区块链客户端
var InitCosmosClient = fsonce.DoWithParam(func(address string) *ClientAdapter {
	ctx := context.Background()
	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix),
		cosmosclient.WithNodeAddress(address))
	if err != nil {
		fslogger.Fatal(err, "Init Client Adapter ERROR")
		return nil
	}
	return &ClientAdapter{&client}
})
