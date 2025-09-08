package gen

import "testing"

func TestGenArticle(t *testing.T) {
	article := NewArticle(3, 3, 20)
	t.Log(article)
}

func TestIdGen(t *testing.T) {
	id, err := NewId()
	t.Log(id, err)
}
