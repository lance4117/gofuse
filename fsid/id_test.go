package fsid

import "testing"

func TestIdGen(t *testing.T) {
	id, err := NewId()
	t.Log(id, err)
}
