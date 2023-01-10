package model

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
