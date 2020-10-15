package validator

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var BookingErr map[string]string

var BookingName string

type BookingValidator struct {
	CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

func init() {
	BookingName = "bookabledate"
	BookingErr = make(map[string]string)
	BookingErr[BookingName] = "test"
}

func NewBookingValidator() *BookingValidator {
	return &BookingValidator{}
}
func (b *BookingValidator) Validate() (validator.Func, string) {
	return func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(time.Time)
		if ok {
			today := time.Now()
			if today.After(date) {
				return false
			}
		}
		return true
	}, BookingName
}
