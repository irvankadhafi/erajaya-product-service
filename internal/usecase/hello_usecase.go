package usecase

import (
	"context"
	"github.com/irvankadhafi/go-boilerplate/internal/model"
)

type helloUsecase struct {
	helloRepo model.HelloRepository
}

// NewHelloUsecase instantiate a new hello usecase
func NewHelloUsecase(helloRepo model.HelloRepository) model.HelloUsecase {
	return &helloUsecase{
		helloRepo: helloRepo,
	}
}

func (u *helloUsecase) FindByID(ctx context.Context, id int64) (*model.Greeting, error) {
	// Add Business Logic ...
	return u.helloRepo.FindByID(ctx, id)
}

func (u *helloUsecase) Create(ctx context.Context, g *model.Greeting) error {
	// Add Business Logic ...
	return u.helloRepo.Create(ctx, g)
}

func (u *helloUsecase) Update(ctx context.Context, userID int64, g *model.Greeting) error {
	// Add Business Logic ...
	return u.helloRepo.Update(ctx, g)
}
