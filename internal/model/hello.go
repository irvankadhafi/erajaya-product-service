package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// HelloUsecase :nodoc:
type HelloUsecase interface {
	FindByID(ctx context.Context, id int64) (*Greeting, error)
	Create(ctx context.Context, g *Greeting) error
	Update(ctx context.Context, userID int64, g *Greeting) error
}

// HelloRepository :nodoc:
type HelloRepository interface {
	FindByID(c context.Context, id int64) (*Greeting, error)
	Create(c context.Context, g *Greeting) error
	Update(c context.Context, g *Greeting) error
}

// Greeting :nodoc:
type Greeting struct {
	ID        int64          `gorm:"primary_key" json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at" sql:"DEFAULT:'now()':::STRING::TIMESTAMP" gorm:"->;<-:create"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
