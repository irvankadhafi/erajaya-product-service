package http

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	ErrInvalidArgument         = echo.NewHTTPError(http.StatusBadRequest, setErrorMessage("invalid argument"))
	ErrInternal                = echo.NewHTTPError(http.StatusInternalServerError, setErrorMessage("internal system error"))
	ErrProductNameAlreadyExist = echo.NewHTTPError(http.StatusBadRequest, setErrorMessage("product name already exist"))
)

// httpValidationOrInternalErr return valdiation or internal error
func httpValidationOrInternalErr(err error) error {
	// Memeriksa apakah err merupakan validator.ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Jika tidak ada kesalahan validasi, mengembalikan kesalahan internal
		return ErrInternal
	}

	// Mengubah validator.ValidationErrors menjadi map dengan kunci field dan nilai pesan kesalahan
	fields := make(map[string]string)
	for _, validationError := range validationErrors {
		// Mengambil tag yang digunakan untuk menyebabkan kesalahan validasi
		tag := validationError.Tag()
		// Menambahkan kesalahan validasi ke map
		fields[validationError.Field()] = fmt.Sprintf("Failed on the '%s' tag", tag)
	}

	// Mengembalikan kesalahan HTTP dengan status bad request dan daftar field yang gagal validasi
	return echo.NewHTTPError(http.StatusBadRequest, setErrorMessage(fields))
}
