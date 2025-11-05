package gen

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Person struct {
	Name     string
	Email    string
	Phone    string
	Gender   string
	Birthday string
	Address  string
}

type Product struct {
	ID          string
	Name        string
	Category    string
	Description string
	Price       float64
	Currency    string
	Stock       int
	CreatedAt   time.Time
}

type OrderItem struct {
	ProductID string
	Name      string
	UnitPrice float64
	Qty       int
}

type Order struct {
	ID        string
	User      string
	Items     []OrderItem
	Total     float64
	Currency  string
	Status    string // created/paid/shipped/finished/refunded
	CreatedAt time.Time
}

// Choice 从切片里等概率选择一个元素（空切片返回零值与 false）
func Choice[T any](arr []T) (T, bool) {
	var zero T
	if len(arr) == 0 {
		return zero, false
	}
	return arr[IntN(len(arr))], true
}

// NewPerson 生成一个随机人
func NewPerson() Person {
	p := gofakeit.Person()
	return Person{
		Name:     p.FirstName + " " + p.LastName,
		Email:    gofakeit.Email(),
		Phone:    gofakeit.Phone(),
		Gender:   p.Gender,
		Birthday: gofakeit.Date().Format("2006-01-02"),
		Address:  p.Address.Address,
	}
}

// NewProduct 生成模拟商品
func NewProduct() Product {
	price := FloatRange(19, 9999)
	return Product{
		ID:          ShortID(),
		Name:        gofakeit.ProductName(),
		Category:    gofakeit.RandomString([]string{"digital", "clothing", "book", "food", "toy"}),
		Description: Sentence(12),
		Price:       price,
		Currency:    gofakeit.Cusip(),
		Stock:       IntRange(0, 9999),
		CreatedAt:   time.Now(),
	}
}

// NewOrder 生成模拟订单（含若干条目）
func NewOrder(itemMin, itemMax int) Order {
	if itemMin < 1 {
		itemMin = 1
	}
	if itemMax < itemMin {
		itemMax = itemMin
	}
	items := make([]OrderItem, 0, itemMax)
	count := IntRange(itemMin, itemMax)
	var total float64
	for i := 0; i < count; i++ {
		p := NewProduct()
		qty := IntRange(1, 3)
		items = append(items, OrderItem{ProductID: p.ID, Name: p.Name, UnitPrice: p.Price, Qty: qty})
		total += p.Price * float64(qty)
	}
	return Order{
		ID:        "ORD-" + ShortID(),
		User:      gofakeit.Name(),
		Items:     items,
		Total:     total,
		Currency:  gofakeit.CurrencyShort(),
		Status:    gofakeit.RandomString([]string{"created", "paid", "shipped", "finished", "refunded"}),
		CreatedAt: NowRecent(time.Hour * 24),
	}
}
