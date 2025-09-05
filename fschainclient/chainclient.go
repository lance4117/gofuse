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
	client *cosmosclient.Client
}

var InitCosmosClient = fsonce.DoWithParam(func(address string) *ClientAdapter {
	ctx := context.Background()
	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix),
		cosmosclient.WithNodeAddress(address))
	if err != nil {
		fslogger.Fatal(err, "InitCosmosClient ClientAdapter ERROR")
		return nil
	}
	return &ClientAdapter{&client}
})
