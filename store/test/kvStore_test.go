package test

import (
	"testing"

	"github.com/lance4117/gofuse/store/kvs"
)

func TestPebble(t *testing.T) {
	config := kvs.NewPebbleConfig("./data")
	peb, err := kvs.NewPebbleKV(config)
	if err != nil {
		t.Fatal(err)
	}

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

func TestRedis(t *testing.T) {
	config := kvs.NewRedisConfig("localhost:6379", "", 0, 100)
	store, err := kvs.NewRedisKV(config)
	if err != nil {
		t.Skip("redis not available:", err)
	}
	get, err := store.Get("ok")
	if err != nil {
		t.Skip("redis not available for get:", err)
	}
	t.Log(get)

	err = store.Put("foo", []byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	has, err := store.Has("foo")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(has)

	data, err := store.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
