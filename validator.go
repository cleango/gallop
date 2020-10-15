package gallop

import (
	"github.com/go-playground/validator/v10"
)

type IValidator interface {
	Validate() (f validator.Func, name string)
}
