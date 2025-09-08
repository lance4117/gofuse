package chain

import (
	"testing"
)

func TestInitClient(t *testing.T) {
	adapter := InitClient(DefaultOptions)

	// Get account from the keyring
	// 需要提前存储到keyring： gaiad keys add user
	acc, addr, err := adapter.AccountAndAddress("cosmos1rdx27mfxehx4z45wvw0c7d6hyn78gshu065420")

	t.Log(acc, addr, err)
}
