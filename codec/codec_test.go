package codec

import (
	"encoding/base64"
	"testing"
)

type User struct {
	Uid  int    `json:"uid"`
	Name string `json:"name"`
}

func TestMsgPack(t *testing.T) {
	user1 := User{111, "name111"}

	bytes, err := MPMarshal(&user1)
	if err != nil {
		t.Fatal(err)
	}

	user, err := MPUnmarshalTo[User](bytes)
	if err != nil {
		t.Fatal(err)
	}
	if user.Uid != user1.Uid || user.Name != user1.Name {
		t.Fatalf("MPUnmarshalTo mismatch, got %+v", user)
	}

	var usertest User
	err = MPUnmarshal(bytes, &usertest)
	if err != nil {
		t.Fatal(err)
	}

	if usertest.Uid != user1.Uid || usertest.Name != user1.Name {
		t.Fatalf("MPUnmarshal mismatch, got %+v", usertest)
	}
}

func TestJSON(t *testing.T) {
	user1 := User{111, "name111"}

	bytes, err := JSONMarshal(&user1)
	if err != nil {
		t.Fatal(err)
	}

	user, err := JSONUnmarshalTo[User](bytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)

	var usertest User
	err = JSONUnmarshal(bytes, &usertest)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(usertest)
}

func TestBase64(t *testing.T) {
	// 结构体 <-> base64（标准）
	s, err := B64EncodeJSON(User{Uid: 112312, Name: "Tom"})
	if err != nil {
		t.Fatal(err)
	}

	// 验证一下底层 JSON
	dec, _ := base64.StdEncoding.DecodeString(s)
	t.Log("json:", string(dec)) // 应该是 {"uid":1,"name":"Tom"}

	u, err := B64DecodeJSON[User](s)
	if err != nil {
		t.Fatal(err)
	}

	// URL 场景
	s2, err := B64URLEncodeJSON(u)
	if err != nil {
		t.Fatal(err)
	}
	u2, err := B64URLDecodeJSON[User](s2)
	if err != nil {
		t.Fatal(err)
	}

	// 原始 []byte <-> base64
	b64 := B64Encode([]byte("hello"))
	raw, err := B64Decode(b64)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("out:", s, u, u2, raw)
}
