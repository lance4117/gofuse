package fschain

import (
	"testing"
)

func TestInitClient(t *testing.T) {
	adapter := InitClient(DefaultOptions)

	// Get account from the keyring
	// 需要提前存储到keyring： gaiad keys add user
	account, err := adapter.Account("cosmos15sya5v5hzmx37xuaw6jx2rl2gnla6z3j8v93mr")
	if err != nil {
		t.Fatal(err)
	}

	addr, err := adapter.Address("cosmos15sya5v5hzmx37xuaw6jx2rl2gnla6z3j8v93mr")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(account, addr)
}
