package repository

import (
	"context"
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository create new repository
func NewProductRepository(db *gorm.DB) model.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (p *productRepository) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"productID": id,
	})

	var product model.Product
	err := p.db.WithContext(ctx).Take(&product, "id = ?", id).Error
	switch err {
	case nil:
		return &product, nil
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}
}

func (p *productRepository) Create(ctx context.Context, product *model.Product) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":     utils.DumpIncomingContext(ctx),
		"product": utils.Dump(product),
	})

	err := p.db.WithContext(ctx).Create(product).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (p *productRepository) SearchByPage(ctx context.Context, criteria model.ProductSearchCriteria) (ids []int64, count int64, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":      utils.DumpIncomingContext(ctx),
		"criteria": utils.Dump(criteria),
	})

	err = p.db.WithContext(ctx).Model(model.Product{}).
		Count(&count).
		Error
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}

	if count <= 0 {
		return nil, 0, nil
	}

	err = p.db.WithContext(ctx).
		Model(model.Product{}).
		Scopes(
			orderByCreatedAtAsc,
			withNameLike(criteria.Query),
			withSize(criteria.Size),
			withOffset(utils.Offset(criteria.Page, criteria.Size)),
		).
		Pluck("id", &ids).Error
	switch err {
	case nil:
		return ids, count, nil
	case gorm.ErrRecordNotFound:
		return nil, 0, nil
	default:
		logger.Error(err)
		return nil, 0, err
	}
}
