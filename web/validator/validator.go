package validator

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var validate *validator.Validate
var once sync.Once

func New() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}
