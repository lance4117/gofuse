package conv

import "testing"

func TestBytesInt64(t *testing.T) {
	bytes := Int64ToBytes(1153202979583686590)
	t.Log(bytes)
	n, err := BytesToInt64(bytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n)
}
