package chain

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	config := DefaultConfig("zero", "D:\\code\\zerod")

	client, err := New(config)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	status, err := client.Status()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config)
	t.Log(status)
}
