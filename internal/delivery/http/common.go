package http

import (
	"github.com/go-playground/validator"
	"sync"
)

// validate singleton, it's thread safe and cached the struct validation rules
var validate *validator.Validate

var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		validate = validator.New()
	})
}

type successResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func setSuccessResponse(data any) successResponse {
	return successResponse{
		Success: true,
		Data:    data,
	}
}

type errorResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
}

func setErrorMessage(msg any) errorResponse {
	return errorResponse{
		Success: false,
		Message: msg,
	}
}
