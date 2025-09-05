package fschainclient

import (
	"testing"
)

func TestInitClient(t *testing.T) {
	adapter := InitCosmosClient(DefaultAddress)

	// Get account from the keyring
	// 需要提前存储到keyring： gaiad keys add user
	account, err := adapter.Client.Account("cosmos1q3dllmcgcx0x3a5y68pytma2zcjqt2ud69c79c")
	if err != nil {
		t.Fatal(err)
	}

	addr, err := adapter.Client.Address("cosmos1q3dllmcgcx0x3a5y68pytma2zcjqt2ud69c79c")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(account, addr)
}
