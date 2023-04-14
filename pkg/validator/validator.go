package validator

import "github.com/go-playground/validator/v10"

type validatorImpl struct {
	v *validator.Validate
}

func (cv *validatorImpl) Validate(i interface{}) error {
	return cv.v.Struct(i)
}

type Validator interface {
	Validate(i interface{}) error
}

func NewValidator() Validator {
	return &validatorImpl{
		v: validator.New(),
	}
}
