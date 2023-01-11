package usecase

import (
	"context"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/internal/model/mock"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestProductUsecase_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	mockRepo := mock.NewMockProductRepository(ctrl)
	usecase := &productUsecase{
		productRepo: mockRepo,
	}

	id := int64(123)
	resp := &model.Product{ID: id}

	t.Run("ok", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, id).Times(1).Return(resp, nil)
		res, err := usecase.FindByID(ctx, id)
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("failed - product not found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, id).Times(1).Return(nil, nil)
		res, err := usecase.FindByID(ctx, id)
		require.Error(t, err)
		require.EqualError(t, err, ErrNotFound.Error())
		require.Nil(t, res)
	})
}

func TestProductUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	mockRepo := mock.NewMockProductRepository(ctrl)
	usecase := &productUsecase{
		productRepo: mockRepo,
	}

	product := &model.Product{
		ID:          utils.GenerateID(),
		Name:        "Apple iPhone 14 Pro Max",
		Slug:        "apple-iphone-14-pro-max",
		Price:       decimal.New(19000000, 0),
		Description: "Apple iPhone 14 Pro Max",
		Quantity:    20,
	}

	t.Run("ok", func(t *testing.T) {
		mockRepo.EXPECT().Create(ctx, gomock.Any()).Times(1).Return(nil)
		mockRepo.EXPECT().FindByID(ctx, gomock.Any()).Times(1).Return(product, nil)

		product, err := usecase.Create(ctx, model.CreateProductInput{
			Name:        product.Name,
			Description: product.Description,
			Quantity:    product.Quantity,
		})
		require.NoError(t, err)
		require.NotNil(t, product)
	})
}

func TestProductUsecase_FindAllByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	mockRepo := mock.NewMockProductRepository(ctrl)
	usecase := &productUsecase{
		productRepo: mockRepo,
	}

	testProducts := []*model.Product{
		{
			ID: int64(111),
		},
		{
			ID: int64(222),
		},
	}

	t.Run("ok - found all", func(t *testing.T) {
		patch := gomonkey.ApplyMethod(reflect.TypeOf(usecase), "FindByID", func(s *productUsecase, ctx context.Context, id int64) (*model.Product, error) {
			return testProducts[0], nil
		})
		defer patch.Reset()

		res := usecase.FindAllByIDs(ctx, []int64{testProducts[0].ID, testProducts[1].ID})
		assert.NotNil(t, res)
	})

	t.Run("ok - found none", func(t *testing.T) {
		patch := gomonkey.ApplyMethod(reflect.TypeOf(usecase), "FindByID", func(s *productUsecase, ctx context.Context, id int64) (*model.Product, error) {
			return nil, nil
		})
		defer patch.Reset()

		res := usecase.FindAllByIDs(ctx, []int64{testProducts[0].ID, testProducts[1].ID})
		assert.Nil(t, res)
	})
}

func TestProductUsecase_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	mockRepo := mock.NewMockProductRepository(ctrl)
	usecase := &productUsecase{
		productRepo: mockRepo,
	}

	testSearchCriteria := model.ProductSearchCriteria{
		Page:     1,
		Size:     2,
		SortType: model.ProductSortTypeCreatedAtAsc,
	}
	ids := []int64{111, 222}

	resp := []*model.Product{
		{ID: ids[0]},
		{ID: ids[1]},
	}
	expectedCount := int64(len(ids))

	t.Run("ok", func(t *testing.T) {
		mockRepo.EXPECT().SearchByPage(ctx, testSearchCriteria).Times(1).Return(ids, expectedCount, nil)
		patch := gomonkey.ApplyMethod(reflect.TypeOf(usecase), "FindAllByIDs", func(s *productUsecase, ctx context.Context, ids []int64) []*model.Product {
			return resp
		})
		defer patch.Reset()

		spaces, count, err := usecase.Search(ctx, testSearchCriteria)

		require.NoError(t, err)
		assert.Equal(t, resp, spaces)
		assert.Equal(t, expectedCount, count)
	})

	t.Run("ok - found none", func(t *testing.T) {
		mockRepo.EXPECT().SearchByPage(ctx, testSearchCriteria).Times(1).Return([]int64{}, int64(0), nil)
		spaces, count, err := usecase.Search(ctx, testSearchCriteria)
		require.NoError(t, err)
		assert.Empty(t, spaces)
		assert.Empty(t, count)
	})
}
