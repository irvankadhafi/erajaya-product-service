package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/irvankadhafi/erajaya-product-service/cache"
	"github.com/irvankadhafi/erajaya-product-service/internal/config"
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type productRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewProductRepository create new repository
func NewProductRepository(db *gorm.DB, cache cache.Cache) model.ProductRepository {
	return &productRepository{
		db:    db,
		cache: cache,
	}
}

func (p *productRepository) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"productID": id,
	})

	cacheKey := p.newCacheKeyByID(id)
	if !config.DisableCaching() {
		reply, err := p.findFromCacheByKey(cacheKey)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		if reply != nil {
			return reply, nil
		}
	}

	product := &model.Product{}
	err := p.db.WithContext(ctx).Take(product, "id = ?", id).Error
	switch err {
	case nil:
		err := p.cache.Store(cache.NewItem(cacheKey, utils.Dump(product)))
		if err != nil {
			logger.Error(err)
		}
		return product, nil
	case gorm.ErrRecordNotFound:
		p.storeNilCacheByKey(cacheKey)
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

	err = p.cache.DeleteByKeys([]string{p.newCacheKeyByID(product.ID)})
	if err != nil {
		logger.Error(err)
	}
	return nil
}

func (p *productRepository) SearchByPage(ctx context.Context, criteria model.ProductCriteria) (ids []int64, count int64, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":      utils.DumpIncomingContext(ctx),
		"criteria": utils.Dump(criteria),
	})

	err = p.db.Debug().WithContext(ctx).Model(model.Product{}).
		Count(&count).
		Error
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}

	if count <= 0 {
		return nil, 0, nil
	}

	err = p.db.Debug().WithContext(ctx).
		Model(model.Product{}).
		Scopes(scopeByPageAndLimit(criteria.Page, criteria.Size)).
		Order(orderByProductSortType(criteria.SortType)).
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

func (p *productRepository) newCacheKeyByID(id int64) string {
	return fmt.Sprintf("cache:object:product:id:%d", id)
}

func (p *productRepository) findFromCacheByKey(key string) (reply *model.Product, err error) {
	var rep interface{}
	rep, err = p.cache.Get(key)
	if err != nil || rep == nil {
		return
	}

	bt, _ := rep.([]byte)
	if bt == nil {
		return
	}

	if err = json.Unmarshal(bt, &reply); err != nil {
		return
	}

	return
}

func (p *productRepository) storeNilCacheByKey(key string) {
	err := p.cache.StoreNil(key)
	if err != nil {
		logrus.Error(err)
	}
}

func orderByProductSortType(sortType model.ProductSortType) string {
	if orderBy, ok := model.QueryProductSortByMap[sortType]; ok {
		return orderBy
	}

	return model.QueryProductSortByMap[model.ProductSortTypeCreatedAtDesc]
}
