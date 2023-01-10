package httpsvc

import (
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
	s.echo.POST("/products", s.handleCreateProduct())
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

		return c.JSON(http.StatusCreated, product)
	}
}
