package httpsvc

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type errorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func setErrorMessage(msg string) errorResponse {
	return errorResponse{
		Success: false,
		Message: msg,
	}
}

var (
	ErrInvalidArgument = echo.NewHTTPError(http.StatusBadRequest, setErrorMessage("invalid argument"))
	ErrInternal        = echo.NewHTTPError(http.StatusInternalServerError, setErrorMessage("internal system error"))
)
