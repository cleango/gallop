package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

func TranslateCustomErr(err error, errM map[string]string) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, err := range errs {
			if v, ok := errM[err.Tag()]; ok {
				return errors.New(v)
			}
		}
	}
	return err

}
