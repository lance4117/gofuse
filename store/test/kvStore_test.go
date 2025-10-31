package test

import (
	"testing"

	"github.com/lance4117/gofuse/store/kvs"
)

func TestPebble(t *testing.T) {
	config := kvs.DefaultPebbleConfig("./data")
	peb := kvs.NewPebbleNoSync(config)

	get, err := peb.Get("a")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(get)

}
