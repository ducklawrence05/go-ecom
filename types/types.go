package types

import (
	"context"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CartCheckoutItem struct {
	ProductID int `json:"productID" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userID"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderID"`
	ProductID int       `json:"productID"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

// store
type UserStore interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	CreateUser(ctx context.Context, user User) error
}

type ProductStore interface {
	GetProducts(ctx context.Context) ([]Product, error)
	GetProductsByIDs(ctx context.Context, productIDs []int) ([]Product, error)
	CreateProduct(ctx context.Context, product Product) error
	UpdateProduct(ctx context.Context, product Product) error
}

type OrderStore interface {
	CreateOrder(ctx context.Context, order Order) (int, error)
	CreateOrderItem(ctx context.Context, orderItem OrderItem) error
}

// payload
type RegisterUserPayLoad struct {
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastName" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=30"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type LoginUserPayLoad struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateProductPayLoad struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Image       string  `json:"image" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
}

type CreateOrderPayLoad struct {
	UserID  int     `json:"userID" validate:"required"`
	Total   float64 `json:"total" validate:"required,gt=0"`
	Status  string  `json:"status" validate:"required,oneof=pending completed cancelled"`
	Address string  `json:"address" validate:"required"`
}

type CartCheckoutPayLoad struct {
	Items []CartCheckoutItem `json:"items" validate:"required"`
}
