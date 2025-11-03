package test

import (
	"testing"

	"github.com/lance4117/gofuse/store/kvs"
)

func TestPebble(t *testing.T) {
	config := kvs.DefaultPebbleConfig("./data")
	peb := kvs.NewPebbleKV(config)

	has, err := peb.Has("a")
	if err != nil {
		t.Error(err)
	}

	t.Log(has)

	get, err := peb.Get("a")
	if err != nil {
		t.Error(err)
	}

	t.Log(get)

	err = peb.Put("test", []byte("foo"))
	if err != nil {
		t.Fatal(err)
	}

	has, err = peb.Has("test")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(has)

	data, err := peb.Get("test")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))

}
