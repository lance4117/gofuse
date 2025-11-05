package crypt

import (
	"testing"
)

func TestArgon2id(t *testing.T) {
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

func TestAESGCM(t *testing.T) {
	encryption := New()
	str := "data data"
	data, err := encryption.EncryptAESGCM([]byte("pass"), []byte(str))
	if err != nil {
		t.Fatal(err)
	}
	result, err := encryption.DecryptAESGCM([]byte("pass"), data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(result))
}
