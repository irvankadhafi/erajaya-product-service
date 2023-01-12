package repository

import (
	"context"
	"fmt"
	"github.com/irvankadhafi/erajaya-product-service/cache"
	"github.com/irvankadhafi/erajaya-product-service/internal/config"
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
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

// FindByID find product by id
func (p *productRepository) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"productID": id,
	})

	cacheKey := p.newCacheKeyByID(id)
	if !config.DisableCaching() {
		reply, err := findFromCacheByKey[*model.Product](p.cache, cacheKey)
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
		storeNilCacheByKey(p.cache, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}
}

// FindBySlug find product with specific slug
func (p *productRepository) FindBySlug(ctx context.Context, slug string) (*model.Product, error) {
	logger := logrus.WithFields(logrus.Fields{
		"context": utils.DumpIncomingContext(ctx),
		"slug":    slug})

	cacheKey := p.newCacheKeyBySlug(slug)
	if !config.DisableCaching() {
		reply, err := findFromCacheByKey[int64](p.cache, cacheKey)

		if err != nil {
			logger.Error(err)
			return nil, err
		}
		if reply > 0 {
			return p.FindByID(ctx, reply)
		}
	}

	var product *model.Product
	var id int64
	var ids []int64
	err := p.db.WithContext(ctx).Model(product).Where("slug = ?", slug).Pluck("id", &ids).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(ids) == 0 {
		storeNilCacheByKey(p.cache, cacheKey)
		return nil, nil
	}

	id = ids[0]

	err = p.cache.Store(cache.NewItem(cacheKey, id))
	if err != nil {
		logger.Error(err)
	}

	return p.FindByID(ctx, id)
}

// Create product
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

	if err = p.deleteCaches(product.ID); err != nil {
		logger.Error(err)
	}
	return nil
}

// SearchByPage find all product with specific criteria
func (p *productRepository) SearchByPage(ctx context.Context, criteria model.ProductSearchCriteria) (ids []int64, count int64, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":      utils.DumpIncomingContext(ctx),
		"criteria": utils.Dump(criteria),
	})

	cacheKey := p.newProductCacheKeyByCriteria(criteria)
	if !config.DisableCaching() {
		products, err := findFromCacheByKey[*model.MultiCacheValue](p.cache, cacheKey)
		if err != nil {
			logger.Error(err)
			return nil, 0, err
		}

		if products != nil {
			return products.IDs, products.Count, nil
		}
	}

	count, err = p.countAll(ctx, criteria)
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}

	if count <= 0 {
		storeNilMultiValueByKey(p.cache, cacheKey)
		return nil, 0, nil
	}

	ids, err = p.findAllIDsByCriteria(ctx, criteria)
	switch err {
	case nil:
		cacheItem := cache.NewItemWithCustomTTL(
			cacheKey,
			utils.Dump(model.MultiCacheValue{IDs: ids, Count: count}),
			15*time.Second,
		)
		if err := p.cache.Store(cacheItem); err != nil {
			logger.Error(err)
		}
		return ids, count, nil
	case gorm.ErrRecordNotFound:
		return nil, 0, nil
	default:
		logger.Error(err)
		return nil, 0, err
	}
}

func (p *productRepository) findAllIDsByCriteria(ctx context.Context, criteria model.ProductSearchCriteria) ([]int64, error) {
	var scopes []func(*gorm.DB) *gorm.DB
	scopes = append(scopes, scopeByPageAndLimit(criteria.Page, criteria.Size))
	if criteria.Query != "" {
		scopes = append(scopes, scopeMatchTSQuery(criteria.Query))
	}

	var ids []int64
	err := p.db.WithContext(ctx).
		Model(model.Product{}).
		Scopes(scopes...).
		Order(orderByProductSortType(criteria.SortType)).
		Pluck("id", &ids).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":      utils.DumpIncomingContext(ctx),
			"criteria": utils.Dump(criteria),
		}).Error(err)
		return nil, err
	}

	return ids, nil
}

func (p *productRepository) countAll(ctx context.Context, criteria model.ProductSearchCriteria) (int64, error) {
	var scopes []func(*gorm.DB) *gorm.DB
	if criteria.Query != "" {
		scopes = append(scopes, scopeMatchTSQuery(criteria.Query))
	}

	var count int64
	err := p.db.WithContext(ctx).Model(model.Product{}).
		Scopes(scopes...).
		Count(&count).
		Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":      utils.DumpIncomingContext(ctx),
			"criteria": utils.Dump(criteria),
		}).Error(err)
		return 0, err
	}

	return count, nil
}

func (p *productRepository) newCacheKeyByID(id int64) string {
	return fmt.Sprintf("cache:object:product:id:%d", id)
}

func (p *productRepository) newCacheKeyBySlug(slug string) string {
	return fmt.Sprintf("cache:object:product:slug:%s", slug)
}

func (p *productRepository) newProductCacheKeyBucket() string {
	return fmt.Sprintf("cache:object:product")
}

func (p *productRepository) newProductCacheKeyByCriteria(criteria model.ProductSearchCriteria) string {
	key := fmt.Sprintf("cache:object:productMultiValue:page:%d:size:%d:sortType:%s", criteria.Page, criteria.Size, string(criteria.SortType))

	if criteria.Query != "" {
		return key + ":query:" + criteria.Query
	}

	return key
}

func (p *productRepository) deleteCaches(productID int64) error {
	productCacheKeyBucket := p.newProductCacheKeyBucket()
	productCacheKey := p.newCacheKeyByID(productID)

	err := p.cache.DeleteByKeys([]string{
		productCacheKeyBucket,
		productCacheKey,
	})
	if err != nil {
		return err
	}

	return nil
}

func orderByProductSortType(sortType model.ProductSortType) string {
	if orderBy, ok := model.QueryProductSortByMap[sortType]; ok {
		return orderBy
	}

	return model.QueryProductSortByMap[model.ProductSortTypeCreatedAtDesc]
}
