package usecase

import (
	"context"
	"github.com/gosimple/slug"
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/sirupsen/logrus"
	"sync"
)

type productUsecase struct {
	productRepo model.ProductRepository
}

// NewProductUsecase instantiate a new product usecase
func NewProductUsecase(productRepo model.ProductRepository) model.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

// FindByID find product by specific id
func (p *productUsecase) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"productID": id,
	})

	product, err := p.productRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if product == nil {
		return nil, ErrNotFound
	}

	return product, nil
}

// FindAllByIDs find all products with IDs
func (p *productUsecase) FindAllByIDs(ctx context.Context, ids []int64) []*model.Product {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"ids": ids,
	})

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	// Creating channel to receive found products
	c := make(chan *model.Product, len(ids))

	// Iterating through received ids
	for _, id := range ids {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			product, err := p.FindByID(ctx, id)
			if err != nil {
				logger.Error(err)
				return
			}
			c <- product
		}(id)
	}
	wg.Wait()
	close(c)

	if len(c) <= 0 {
		return nil
	}

	// put all products in a map with product id as key
	rs := map[int64]*model.Product{}
	for product := range c {
		if product != nil {
			rs[product.ID] = product
		}
	}

	// sort products based on the order of received ids
	var products []*model.Product
	for _, id := range ids {
		if product, ok := rs[id]; ok {
			products = append(products, product)
		}
	}

	return products
}

// Create product from input
func (p *productUsecase) Create(ctx context.Context, input model.CreateProductInput) (*model.Product, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	err := input.Validate()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	product := &model.Product{
		ID:          utils.GenerateID(),
		Name:        input.Name,
		Slug:        slug.Make(input.Name),
		Price:       input.Price,
		Description: input.Description,
		Quantity:    input.Quantity,
	}

	existingProduct, err := p.productRepo.FindBySlug(ctx, product.Slug)
	if err != nil && err != ErrNotFound {
		logger.Error(err)
		return nil, err
	}
	if existingProduct != nil {
		return nil, ErrAlreadyExist
	}

	err = p.productRepo.Create(ctx, product)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return p.FindByID(ctx, product.ID)
}

// Search product with given search criteria
func (p *productUsecase) Search(ctx context.Context, criteria model.ProductSearchCriteria) (products []*model.Product, count int64, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":      utils.DumpIncomingContext(ctx),
		"criteria": utils.Dump(criteria),
	})

	criteria.SetDefaultValue()
	productIDs, count, err := p.productRepo.SearchByPage(ctx, criteria)
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}
	if len(productIDs) <= 0 || count <= 0 {
		return nil, 0, err
	}

	products = p.FindAllByIDs(ctx, productIDs)
	if len(products) <= 0 {
		err = ErrNotFound
		return
	}

	return
}
