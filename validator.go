package web

import "github.com/go-playground/validator"

type Validator interface {
	Validate(val interface{}) error
}

type defaultValidator struct {
	validate *validator.Validate
}

func newDefaultValidator() defaultValidator {
	return defaultValidator{
		validate: validator.New(),
	}
}

func (v defaultValidator) Validate(val interface{}) error {
	return v.validate.Struct(val)
}
