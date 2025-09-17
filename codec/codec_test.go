package codec

import "testing"

type User struct {
	Uid  int
	Name string
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
	t.Log(user)

	var usertest User
	err = MPUnmarshal(bytes, &usertest)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(usertest)
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
