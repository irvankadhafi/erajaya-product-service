package httpsvc

import (
	"github.com/irvankadhafi/go-boilerplate/internal/model"
	"github.com/labstack/echo/v4"
)

// Service http service
type Service struct {
	echo         *echo.Group
	helloUsecase model.HelloUsecase
	// ... usecases
	//httpMiddleware *auth.AuthenticationMiddleware
}

// RouteService add dependencies and use group for routing
func RouteService(
	echo *echo.Group,
	helloUsecase model.HelloUsecase,
	//authMiddleware *auth.AuthenticationMiddleware,
) {
	srv := &Service{
		echo:         echo,
		helloUsecase: helloUsecase,
		// ... usecases
		//httpMiddleware: authMiddleware,
	}
	srv.initRoutes()
}

func (s *Service) initRoutes() {
	// ... routes
}
