package chain

import (
	"testing"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

func TestInitClient(t *testing.T) {
	options := append(DefaultOptions, cosmosclient.WithGasAdjustment(1.5),
		cosmosclient.WithGasPrices("0stake"), cosmosclient.WithHome("D:\\code\\blogd\\blogdata"))
	client := NewClient("cosmos1rdx27mfxehx4z45wvw0c7d6hyn78gshu065420", options)
	if client == nil {
		return
	}
	t.Log(client.Account)
}

func TestPubkeyClient(t *testing.T) {

	address := PubToAddress("AqbnviTZBMHO0lp7X0S/9uSHgXT3Hjlx5pw5Fhu80K5y")

	options := append(DefaultOptions, cosmosclient.WithGasAdjustment(1.5),
		cosmosclient.WithGasPrices("0stake"), cosmosclient.WithHome("D:\\code\\blogd\\blogdata"))
	client := NewClient(address, options)
	if client == nil {
		return
	}
	t.Log(client.Account)
}

func TestNewAccountClient(t *testing.T) {
	key, err := NewDefaultKey()
	if err != nil {
		t.Fatal(err)
	}
	options := append(DefaultOptions, cosmosclient.WithGasAdjustment(1.5),
		cosmosclient.WithGasPrices("0stake"), cosmosclient.WithHome("D:\\code\\blogd\\blogdata"))
	client := NewClient(key.Address, options)
	if client == nil {
		return
	}
	t.Log(client.Account)
}
