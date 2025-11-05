package gen

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Person struct {
	Name     string `json:"name,omitempty" yaml:"name"`
	Email    string `json:"email,omitempty" yaml:"email"`
	Phone    string `json:"phone,omitempty" yaml:"phone"`
	Gender   string `json:"gender,omitempty" yaml:"gender"`
	Birthday string `json:"birthday,omitempty" yaml:"birthday"`
	Address  string `json:"address,omitempty" yaml:"address"`
}

type Product struct {
	ID          string    `json:"id,omitempty" yaml:"id"`
	Name        string    `json:"name,omitempty" yaml:"name"`
	Category    string    `json:"category,omitempty" yaml:"category"`
	Description string    `json:"description,omitempty" yaml:"description"`
	Price       float64   `json:"price,omitempty" yaml:"price"`
	Currency    string    `json:"currency,omitempty" yaml:"currency"`
	Stock       int       `json:"stock,omitempty" yaml:"stock"`
	CreatedAt   time.Time `json:"created_at" yaml:"created_at"`
}

type OrderItem struct {
	ProductID string  `json:"product_id,omitempty"  yaml:"product_id"`
	Name      string  `json:"name,omitempty" yaml:"name"`
	UnitPrice float64 `json:"unit_price,omitempty" yaml:"unit_price"`
	Qty       int     `json:"qty,omitempty" yaml:"qty"`
}

type Order struct {
	ID        string      `json:"id,omitempty" yaml:"id"`
	User      string      `json:"user,omitempty" yaml:"user"`
	Items     []OrderItem `json:"items,omitempty" yaml:"items"`
	Total     float64     `json:"total,omitempty" yaml:"total"`
	Currency  string      `json:"currency,omitempty" yaml:"currency"`
	Status    string      `json:"status,omitempty" yaml:"status"` // created/paid/shipped/finished/refunded
	CreatedAt time.Time   `json:"created_at" yaml:"created_at"`
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
