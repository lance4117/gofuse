package gen

import (
	"testing"
	"time"
)

func TestId(t *testing.T) {
	id, err := NewId()
	t.Log(id, err)
}

func TestArticle(t *testing.T) {
	article := NewArticle(3, 20)
	t.Log(article.Title)
	t.Log(article.Author)
	t.Log(article.Summary)
	t.Log(article.Contents)
}

func TestObjs(t *testing.T) {
	order := NewOrder(5, 10)
	t.Log("orders:", order)

	person := NewPerson()
	t.Log("person:", person)

	product := NewProduct()
	t.Log("product:", product)
}

func TestLettersAndNumbers(t *testing.T) {
	t.Log(Letters(10))

	numbers := LettersAndNumbers(10)
	t.Log(numbers)
}

func TestTime(t *testing.T) {
	for range 10 {
		recent := NowRecent(time.Hour)
		t.Log(recent.Format(time.DateTime))
	}
}

func TestChoice(t *testing.T) {
	var list []Person
	for range 5 {
		person := NewPerson()
		t.Log(person.Address)
		list = append(list, person)
	}

	for i := 0; i < 5; i++ {
		choice, b := Choice(list)
		t.Log(choice, b)
	}
}
