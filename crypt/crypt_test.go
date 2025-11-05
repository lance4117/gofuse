package crypt

import (
	"testing"
)

func TestEnctypt(t *testing.T) {
	encryption := New()
	str := "data data"
	data, err := encryption.EncryptArgon2id("pass", []byte(str))
	if err != nil {
		t.Fatal(err)
	}
	result, err := encryption.DecryptArgon2id("pass", data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(result))
}
