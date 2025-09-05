package fsutils

import "testing"

func TestBytes(t *testing.T) {
	type User struct {
		Uid  int
		Name string
	}
	user1 := User{111, "name111"}

	bytes, err := AnyToBytes(&user1)
	if err != nil {
		t.Fatal(err)
	}

	user, err := BytesToAny[User](bytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(user)
}
