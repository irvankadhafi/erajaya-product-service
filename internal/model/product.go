package model

import (
	"context"
	"time"
)

// Product model
type Product struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Price       int64     `json:"price" sql:"type:decimal(20,0)" gorm:"type:numeric(20,0)"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at" sql:"DEFAULT:'now()':::STRING::TIMESTAMP" gorm:"->;<-:create"` // create & read only
	UpdatedAt   time.Time `json:"updated_at" sql:"DEFAULT:'now()':::STRING::TIMESTAMP"`
}

// ProductRepository repository
type ProductRepository interface {
	FindByID(ctx context.Context, id int64) (*Product, error)
	SearchByPage(ctx context.Context, criteria ProductSearchCriteria) (ids []int64, count int64, err error)
	FindBySlug(ctx context.Context, slug string) (*Product, error)
	Create(ctx context.Context, product *Product) error
}

// ProductUsecase usecase
type ProductUsecase interface {
	FindByID(ctx context.Context, id int64) (*Product, error)
	FindAllByIDs(ctx context.Context, ids []int64) (products []*Product)
	Search(ctx context.Context, criteria ProductSearchCriteria) (products []*Product, count int64, err error)

	Create(ctx context.Context, input CreateProductInput) (*Product, error)
}

// CreateProductInput create product input
type CreateProductInput struct {
	Name        string `json:"name" validate:"required,min=3,max=60"`
	Description string `json:"description" validate:"max=80"`
	Price       int64  `json:"price" validate:"gte=0"`
	Quantity    int    `json:"quantity" validate:"gt=0"`
}

// Validate validate product input
func (c *CreateProductInput) Validate() error {
	return validate.Struct(c)
}

// ProductSortType sort type for product search
type ProductSortType string

const (
	ProductSortTypeCreatedAtAsc  ProductSortType = "CREATED_AT_ASC"
	ProductSortTypeCreatedAtDesc ProductSortType = "CREATED_AT_DESC"
	ProductSortTypePriceAsc      ProductSortType = "PRICE_ASC"
	ProductSortTypePriceDesc     ProductSortType = "PRICE_DESC"
	ProductSortTypeNameAsc       ProductSortType = "NAME_ASC"
	ProductSortTypeNameDesc      ProductSortType = "NAME_DESC"
)

// QueryProductSortByMap sort type to query string map for database ordering
var QueryProductSortByMap = map[ProductSortType]string{
	ProductSortTypeCreatedAtAsc:  "created_at ASC",
	ProductSortTypeCreatedAtDesc: "created_at DESC",
	ProductSortTypePriceAsc:      "price ASC",
	ProductSortTypePriceDesc:     "price DESC",
	ProductSortTypeNameAsc:       "name ASC",
	ProductSortTypeNameDesc:      "name DESC",
}

type ProductSearchCriteria struct {
	Query    string          `json:"query"`
	Page     int             `json:"page"`
	Size     int             `json:"size"`
	SortType ProductSortType `json:"sort_type"`
}
