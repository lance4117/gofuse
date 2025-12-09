package test

import (
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/store/kvs"
)

func TestPebble(t *testing.T) {
	config := kvs.NewPebbleConfig(t.TempDir())
	peb, err := kvs.NewPebbleKV(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = peb.Close()
	})

	has, err := peb.Has("a")
	if err != nil {
		t.Error(err)
	}

	t.Log(has)

	get, err := peb.Get("a")
	if err != nil && !errors.Is(err, errs.ErrKeyNotFound) {
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
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("set TEST_REDIS_ADDR to run redis integration tests")
	}
	db := 0
	if envDB := os.Getenv("TEST_REDIS_DB"); envDB != "" {
		if parsed, err := strconv.Atoi(envDB); err == nil {
			db = parsed
		}
	}
	poolSize := 100
	if envPool := os.Getenv("TEST_REDIS_POOLSIZE"); envPool != "" {
		if parsed, err := strconv.Atoi(envPool); err == nil {
			poolSize = parsed
		}
	}
	config := kvs.NewRedisConfig(addr, os.Getenv("TEST_REDIS_PASSWORD"), db, poolSize)
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
