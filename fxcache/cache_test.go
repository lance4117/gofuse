package fxcache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache, err := NewCache(10 * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = cache.Set("key111", "value111")
	if err != nil {
		t.Fatal(err)
	}
	// 应该覆盖原有内容
	err = cache.Set("key111", "value123")
	if err != nil {
		t.Fatal(err)
	}
	err = cache.Set("key222", "value222")
	if err != nil {
		t.Fatal(err)
	}
	var value any
	err = cache.Get("key111", &value)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(value)
	t.Log(cache.Len())
}
