package http

import (
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/internal/usecase"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type productResponse struct {
	*model.Product
	Price     string `json:"price"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (s *Service) handleCreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		logger := logrus.WithFields(logrus.Fields{
			"ctx": utils.DumpIncomingContext(ctx),
		})

		req := model.CreateProductInput{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}

		product, err := s.productUsecase.Create(ctx, req)
		switch err {
		case nil:
			break
		case usecase.ErrAlreadyExist:
			return ErrProductNameAlreadyExist
		default:
			logger.Error(err)
			return httpValidationOrInternalErr(err)
		}

		return c.JSON(http.StatusCreated, setSuccessResponse(productResponse{
			Product:   product,
			Price:     utils.Int64ToRupiah(product.Price),
			CreatedAt: utils.FormatTimeRFC3339(&product.CreatedAt),
			UpdatedAt: utils.FormatTimeRFC3339(&product.UpdatedAt),
		}))
	}
}

func (s *Service) handleGetAllProducts() echo.HandlerFunc {
	type metaInfo struct {
		Size      int `json:"size"`
		Count     int `json:"count"`
		CountPage int `json:"count_page"`
		Page      int `json:"page"`
		NextPage  int `json:"next_page"`
	}

	type userCursor struct {
		Items    []productResponse `json:"items"`
		MetaInfo *metaInfo         `json:"meta_info"`
	}
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		logger := logrus.WithFields(logrus.Fields{
			"ctx": utils.DumpIncomingContext(ctx),
		})

		page := utils.StringToInt(c.QueryParam("page"))
		size := utils.StringToInt(c.QueryParam("size"))
		sortType := c.QueryParam("sortBy")
		query := c.QueryParam("query")
		criterias := model.ProductSearchCriteria{
			Page:     page,
			Size:     size,
			SortType: model.ProductSortType(sortType),
			Query:    query,
		}

		products, count, err := s.productUsecase.Search(ctx, criterias)
		switch err {
		case nil:
			break
		default:
			logger.Error(err)
			return httpValidationOrInternalErr(err)
		}

		var productResponses []productResponse
		for _, product := range products {
			productResponses = append(productResponses, productResponse{
				Product:   product,
				Price:     utils.Int64ToRupiah(product.Price),
				CreatedAt: utils.FormatTimeRFC3339(&product.CreatedAt),
				UpdatedAt: utils.FormatTimeRFC3339(&product.UpdatedAt),
			})
		}

		hasMore := int(count)-(criterias.Page*criterias.Size) > 0
		res := userCursor{
			Items: productResponses,
			MetaInfo: &metaInfo{
				Size:      size,
				Count:     int(count),
				CountPage: utils.CalculatePages(int(count), criterias.Size),
				Page:      page,
			},
		}
		if hasMore {
			res.MetaInfo.NextPage = page + 1
		}

		return c.JSON(http.StatusOK, setSuccessResponse(res))
	}
}
