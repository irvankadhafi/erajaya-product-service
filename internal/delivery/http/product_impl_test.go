package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/irvankadhafi/erajaya-product-service/internal/helper"
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/internal/model/mock"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHTTP_handleCreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	product := &model.Product{
		ID:          utils.GenerateID(),
		Name:        "Apple iPhone 14 Pro Max",
		Description: "Apple iPhone 14 Pro Max",
		Slug:        "apple-iphone-14-pro-max",
		Price:       19000000,
		Quantity:    10,
	}

	inputProduct := model.CreateProductInput{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
	}

	mockProductUsecase := mock.NewMockProductUsecase(ctrl)
	server := &Service{
		productUsecase: mockProductUsecase,
	}

	t.Run("ok", func(t *testing.T) {
		ec := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/products/", strings.NewReader(`
		{
			"name": "Apple iPhone 14 Pro Max",
			"description": "Apple iPhone 14 Pro Max",
			"price": 19000000,
			"quantity": 10
		}`,
		))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		echoCtx := ec.NewContext(req, rec)
		ctx := context.TODO()

		mockProductUsecase.EXPECT().Create(ctx, inputProduct).Times(1).Return(product, nil)

		err := server.handleCreateProduct()(echoCtx)
		require.NoError(t, err)

		defer helper.WrapCloser(rec.Result().Body.Close)

		resBody := map[string]interface{}{}
		err = json.NewDecoder(rec.Result().Body).Decode(&resBody)
		require.NoError(t, err)
		require.EqualValues(t, http.StatusCreated, rec.Result().StatusCode)
	})

	t.Run("handle internal server error", func(t *testing.T) {
		ec := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/products/", strings.NewReader(`
		{
			"name": "Apple iPhone 14 Pro Max",
			"description": "Apple iPhone 14 Pro Max",
			"price": 19000000,
			"quantity": 10
		}`,
		))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		echoCtx := ec.NewContext(req, rec)
		ctx := context.TODO()

		mockProductUsecase.EXPECT().Create(ctx, inputProduct).Times(1).Return(nil, errors.New("usecase error"))

		err := server.handleCreateProduct()(echoCtx)
		ec.DefaultHTTPErrorHandler(err, echoCtx)
		require.Error(t, err)
		require.EqualValues(t, http.StatusInternalServerError, rec.Result().StatusCode)
	})
}

func TestHTTP_handleGetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUsecase := mock.NewMockProductUsecase(ctrl)
	svc := &Service{
		productUsecase: mockProductUsecase,
	}

	criteria := model.ProductSearchCriteria{
		Page:     1,
		Size:     4,
		SortType: model.ProductSortTypeNameDesc,
	}

	products := []*model.Product{
		{ID: utils.GenerateID()},
		{ID: utils.GenerateID()},
		{ID: utils.GenerateID()},
	}

	urlVal := url.Values{}
	urlVal.Add("page", utils.Int64ToString(int64(criteria.Page)))
	urlVal.Add("size", utils.Int64ToString(int64(criteria.Size)))
	urlVal.Add("sortBy", string(criteria.SortType))
	queryParam := urlVal.Encode()
	logrus.Warn("QUERY PARAM >>>>> ", queryParam)

	t.Run("ok", func(t *testing.T) {
		ec := echo.New()
		req := httptest.NewRequest("GET", "/products?"+queryParam, nil)
		rec := httptest.NewRecorder()
		echoCtx := ec.NewContext(req, rec)

		mockProductUsecase.EXPECT().Search(gomock.Any(), criteria).Times(1).Return(products, int64(len(products)), nil)

		err := svc.handleGetAllProducts()(echoCtx)
		require.NoError(t, err)

		_, err = io.ReadAll(rec.Result().Body)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	})

	t.Run("handle error", func(t *testing.T) {
		ec := echo.New()
		ec.HTTPErrorHandler = ec.DefaultHTTPErrorHandler
		req := httptest.NewRequest("GET", "/products?"+queryParam, nil)
		rec := httptest.NewRecorder()
		echoCtx := ec.NewContext(req, rec)

		mockProductUsecase.EXPECT().Search(gomock.Any(), criteria).Times(1).Return(nil, int64(0), errors.New("usecase error"))

		err := svc.handleGetAllProducts()(echoCtx)
		require.Error(t, err)
	})
}
