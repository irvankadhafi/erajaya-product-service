package http

import (
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/labstack/echo/v4"
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
