package httpsvc

import (
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
)

// Service http service
type Service struct {
	echo           *echo.Group
	productUsecase model.ProductUsecase
}

// RouteService add dependencies and use group for routing
func RouteService(
	echo *echo.Group,
	productUsecase model.ProductUsecase,
) {
	srv := &Service{
		echo:           echo,
		productUsecase: productUsecase,
	}
	srv.initRoutes()
}

func (s *Service) initRoutes() {
	s.echo.GET("/products/", s.handleGetAllProducts())
	s.echo.POST("/products/", s.handleCreateProduct())
}

type successResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
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
		default:
			logger.Error(err)
			return ErrInternal
		}

		return c.JSON(http.StatusCreated, successResponse{
			Success: true,
			Data:    product,
		})
	}
}

func (s *Service) handleGetAllProducts() echo.HandlerFunc {
	type MetaInfo struct {
		Size       int  `json:"size"`
		Count      int  `json:"count"`
		CountPage  int  `json:"countPage"`
		HasMore    bool `json:"hasMore"`
		Cursor     int  `json:"cursor"`
		NextCursor int  `json:"nextCursor"`
	}

	type userCursor struct {
		Items    []*model.Product `json:"items"`
		MetaInfo *MetaInfo        `json:"meta_info"`
	}
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		logger := logrus.WithFields(logrus.Fields{
			"ctx": utils.DumpIncomingContext(ctx),
		})

		page := utils.StringToInt(c.QueryParam("page"))
		size := utils.StringToInt(c.QueryParam("size"))
		sortType := c.QueryParam("sortBy")
		criterias := model.ProductSearchCriteria{
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
			return ErrInternal
		}

		//var userResponses []userResponse
		//for _, user := range users {
		//	userResponses = append(userResponses, userResponse{
		//		ID:          utils.Int64ToString(user.ID),
		//		Name:        user.Name,
		//		Email:       user.Email,
		//		RoleID:      user.RoleID,
		//		PhoneNumber: user.PhoneNumber,
		//		CreatedAt:   utils.FormatTimeRFC3339(&user.CreatedAt),
		//		UpdatedAt:   utils.FormatTimeRFC3339(&user.UpdatedAt),
		//	})
		//}
		hasMore := int(count)-(criterias.Page*criterias.Size) > 0
		countPage := math.Ceil(float64(count) / float64(criterias.Size))
		res := userCursor{
			Items: items,
			MetaInfo: &MetaInfo{
				Size:      size,
				Count:     int(count),
				CountPage: int(countPage),
				HasMore:   hasMore,
				Cursor:    page,
			},
		}
		if !hasMore {
			res.MetaInfo.NextCursor = 0
		} else {
			res.MetaInfo.NextCursor = page + 1
		}

		return c.JSON(http.StatusOK, res)
	}
}
