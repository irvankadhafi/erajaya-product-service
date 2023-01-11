package httpsvc

import (
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"math"
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
		Size      int  `json:"size"`
		Count     int  `json:"count"`
		CountPage int  `json:"count_page"`
		HasMore   bool `json:"has_more"`
		Page      int  `json:"page"`
		NextPage  int  `json:"next_page"`
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
		criterias := model.ProductCriteria{
			Page:     page,
			Size:     size,
			SortType: model.ProductSortType(sortType),
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
		countPage := math.Ceil(float64(count) / float64(criterias.Size))
		res := userCursor{
			Items: items,
			MetaInfo: &metaInfo{
				Size:      size,
				Count:     int(count),
				CountPage: int(countPage),
				HasMore:   hasMore,
				Page:      page,
			},
		}
		if !hasMore {
			res.MetaInfo.NextPage = 0
		} else {
			res.MetaInfo.NextPage = page + 1
		}

		return c.JSON(http.StatusOK, setSuccessResponse(res))
	}
}
