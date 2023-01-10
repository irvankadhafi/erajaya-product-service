package usecase

import (
	"context"
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/sirupsen/logrus"
	"sync"
)

type productUsecase struct {
	productRepository model.ProductRepository
}

// NewProductUsecase instantiate a new product usecase
func NewProductUsecase(productRepository model.ProductRepository) model.ProductUsecase {
	return &productUsecase{
		productRepository: productRepository,
	}
}

func (p *productUsecase) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"productID": id,
	})

	product, err := p.productRepository.FindByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if product == nil {
		return nil, ErrNotFound
	}

	return product, nil
}

func (p *productUsecase) FindAllByIDs(ctx context.Context, ids []int64) []*model.Product {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"ids": ids,
	})

	var wg sync.WaitGroup
	c := make(chan *model.Product, len(ids))

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

	// Put 'em in a buffer
	rs := map[int64]*model.Product{}
	for product := range c {
		if product != nil {
			rs[product.ID] = product
		}
	}

	// Sort 'em out
	var products []*model.Product
	for _, id := range ids {
		if product, ok := rs[id]; ok {
			products = append(products, product)
		}
	}

	return products
}

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
		ID:   utils.GenerateID(),
		Name: input.Name,
		//Slug:        input.,
		Description: input.Description,
		Quantity:    input.Quantity,
	}
	err = p.productRepository.Create(ctx, product)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return p.FindByID(ctx, product.ID)
}

func (p *productUsecase) Search(ctx context.Context, criteria model.ProductSearchCriteria) (products []*model.Product, count int64, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":      utils.DumpIncomingContext(ctx),
		"criteria": utils.Dump(criteria),
	})

	productIDs, count, err := p.productRepository.SearchByPage(ctx, criteria)
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
