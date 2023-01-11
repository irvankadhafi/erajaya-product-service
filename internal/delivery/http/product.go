package http

import (
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

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
		default:
			logger.Error(err)
			return httpValidationOrInternalErr(err)
		}

		return c.JSON(http.StatusCreated, setSuccessResponse(product))
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
		Items    []*model.Product `json:"items"`
		MetaInfo *metaInfo        `json:"meta_info"`
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

		items, count, err := s.productUsecase.Search(ctx, criterias)
		switch err {
		case nil:
			break
		default:
			logger.Error(err)
			return httpValidationOrInternalErr(err)
		}

		hasMore := int(count)-(criterias.Page*criterias.Size) > 0
		res := userCursor{
			Items: items,
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
