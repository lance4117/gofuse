package chain

import (
	"testing"
)

func TestInitClient(t *testing.T) {
	client := InitClient("cosmos1rdx27mfxehx4z45wvw0c7d6hyn78gshu065420", DefaultOptions)
	if client == nil {
		return
	}
	t.Log(client.Account)
}
