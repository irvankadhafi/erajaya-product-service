package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
)

func TestProductRepository_Create(t *testing.T) {
	initializeTest()
	kit, closer := initializeRepoTestKit(t)
	defer closer()
	mock := kit.dbmock

	ctx := context.TODO()
	repo := &productRepository{
		db: kit.db,
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
		rows := sqlmock.NewRows([]string{
			"id",
			"name",
			"slug",
			"price",
			"description",
			"quantity",
		})
		rows.AddRow(product.ID,
			product.Name,
			product.Slug,
			product.Price,
			product.Description,
			product.Quantity)

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "products"`).WillReturnRows(rows)
		mock.ExpectCommit()

		err := repo.Create(ctx, product)
		require.NoError(t, err)
	})

	t.Run("failed - create product return err", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "products"`).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Create(ctx, product)
		require.Error(t, err)
	})
}

func TestProductRepository_FindByID(t *testing.T) {
	initializeTest()
	kit, closer := initializeRepoTestKit(t)
	defer closer()
	mock := kit.dbmock

	ctx := context.TODO()
	repo := &productRepository{
		db: kit.db,
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
		rows := sqlmock.NewRows([]string{
			"id",
			"name",
			"slug",
			"price",
			"description",
			"quantity",
		})
		rows.AddRow(product.ID,
			product.Name,
			product.Slug,
			product.Price,
			product.Description,
			product.Quantity)

		mock.ExpectQuery("^SELECT .+ FROM \"products\"").WillReturnRows(rows)

		res, err := repo.FindByID(ctx, product.ID)
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("failed - return err", func(t *testing.T) {
		mock.ExpectQuery("^SELECT .+ FROM \"products\"").WillReturnError(errors.New("db error"))

		res, err := repo.FindByID(ctx, product.ID)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("failed - not found", func(t *testing.T) {
		mock.ExpectQuery("^SELECT .+ FROM \"products\"").
			WillReturnError(gorm.ErrRecordNotFound)

		res, err := repo.FindByID(ctx, product.ID)
		require.NoError(t, err)
		require.Nil(t, res)
	})
}

func TestProductRepository_SearchByPage(t *testing.T) {
	initializeTest()
	kit, closer := initializeRepoTestKit(t)
	defer closer()
	mock := kit.dbmock

	ctx := context.TODO()
	repo := &productRepository{
		db: kit.db,
	}

	productIDs := []int64{int64(111), int64(222), int64(333), int64(444)}
	expectedCount := int64(len(productIDs))
	criteria := model.ProductSearchCriteria{
		Page:     1,
		Size:     10,
		SortType: model.ProductSortTypeNameDesc,
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT count(.*) FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(productIDs)))

		rows := sqlmock.NewRows([]string{"id"})
		for _, id := range productIDs {
			rows.AddRow(id)
		}
		mock.ExpectQuery(`^SELECT .+ FROM "products"`).
			WillReturnRows(rows)

		actualProductIDs, count, err := repo.SearchByPage(ctx, criteria)
		require.NoError(t, err)
		require.Equal(t, len(productIDs), len(actualProductIDs))
		require.Equal(t, expectedCount, count)
	})
}
