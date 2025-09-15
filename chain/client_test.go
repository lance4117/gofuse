package chain

import (
	"testing"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

func TestInitClient(t *testing.T) {
	options := append(DefaultOptions, cosmosclient.WithGasAdjustment(1.5),
		cosmosclient.WithGasPrices("0stake"), cosmosclient.WithHome("D:\\code\\blogd\\blogtest"))
	client := InitClient("cosmos1rdx27mfxehx4z45wvw0c7d6hyn78gshu065420", options)
	if client == nil {
		return
	}
	t.Log(client.Account)
}
