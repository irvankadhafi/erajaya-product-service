package repository

import (
	"context"
	"github.com/irvankadhafi/go-boilerplate/internal/model"
	"github.com/irvankadhafi/go-boilerplate/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type helloRepo struct {
	db *gorm.DB
}

// NewHelloRepository create new repository
func NewHelloRepository(db *gorm.DB) model.HelloRepository {
	return &helloRepo{
		db: db,
	}
}

// FindByID find object with specific id
func (r *helloRepo) FindByID(ctx context.Context, id int64) (*model.Greeting, error) {
	var greeting *model.Greeting

	logger := log.WithFields(log.Fields{
		"context": utils.Dump(ctx),
		"id":      id})

	var g model.Greeting
	err := r.db.WithContext(ctx).First(&g, id).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	greeting = &g

	return greeting, nil
}

// Create Greeting
func (r *helloRepo) Create(ctx context.Context, g *model.Greeting) error {
	logger := log.WithFields(log.Fields{
		"context":  utils.DumpIncomingContext(ctx),
		"greeting": utils.Dump(g),
	})
	err := r.db.WithContext(ctx).Create(g).Error
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

// Update Greeting by ID
func (r *helloRepo) Update(ctx context.Context, g *model.Greeting) error {
	tx := r.db.Begin()

	err := tx.Model(g).Updates(g).Error
	if err != nil {
		log.WithFields(log.Fields{
			"context":  utils.DumpIncomingContext(ctx),
			"greeting": utils.Dump(g)}).
			Error(err)
		tx.Rollback()
		return err
	}

	err = tx.First(&g, g.ID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"context": utils.DumpIncomingContext(ctx),
			"id":      g.ID}).
			Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		log.WithField("tx", utils.Dump(tx)).Error(err)
		return err
	}

	return nil
}
