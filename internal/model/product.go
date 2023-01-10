package model

import (
	"context"
	"time"
)

// Product :nodoc:
type Product struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at" sql:"DEFAULT:'now()':::STRING::TIMESTAMP" gorm:"->;<-:create"` // create & read only
	UpdatedAt   time.Time `json:"updated_at" sql:"DEFAULT:'now()':::STRING::TIMESTAMP"`
}

// ProductUsecase :nodoc:
type ProductUsecase interface {
	FindByID(ctx context.Context, id int64) (*Product, error)
	FindAllByIDs(ctx context.Context, ids []int64) (products []*Product)
	Search(ctx context.Context, criteria ProductSearchCriteria) (products []*Product, count int64, err error)

	Create(ctx context.Context, input CreateProductInput) (*Product, error)
}

// ProductRepository :nodoc:
type ProductRepository interface {
	FindByID(ctx context.Context, id int64) (*Product, error)
	SearchByPage(ctx context.Context, criteria ProductSearchCriteria) (ids []int64, count int64, err error)

	Create(ctx context.Context, product *Product) error
}

type CreateProductInput struct {
	Name        string `json:"name" validate:"required,min=3,max=60"`
	Description string `json:"description" validate:"max=80"`
	Quantity    int    `json:"quantity"`
}

// Validate :nodoc:
func (c *CreateProductInput) Validate() error {
	return validate.Struct(c)
}

type ProductSearchCriteria struct {
	Query string `json:"query"`
	Page  int64  `json:"page"`
	Size  int64  `json:"size"`
}
