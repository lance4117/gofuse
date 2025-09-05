package fschainclient

import (
	"testing"
)

func TestInitClient(t *testing.T) {
	adapter := InitClient(DefaultOptions)

	// Get account from the keyring
	// 需要提前存储到keyring： gaiad keys add user
	account, err := adapter.Client.Account("cosmos15sya5v5hzmx37xuaw6jx2rl2gnla6z3j8v93mr")
	if err != nil {
		t.Fatal(err)
	}

	addr, err := adapter.Client.Address("cosmos15sya5v5hzmx37xuaw6jx2rl2gnla6z3j8v93mr")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(account, addr)
}
