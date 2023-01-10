package httpsvc

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	ErrInvalidArgument = echo.NewHTTPError(http.StatusBadRequest, "invalid argument")
	ErrInternal        = echo.NewHTTPError(http.StatusInternalServerError, "internal system error")
)
