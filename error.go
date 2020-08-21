package gallop

import (
	"fmt"
	"net/http"
)

func Throw(err error) {
	if err != nil {
		panic(err)
	}
}
func ThrowWithErr(err, base error) {
	if err != nil {
		panic(fmt.Errorf("%w %s", base, err.Error()))
	}
}

type WebError struct {
	StatusCode int
	Data       interface{}
}

func Error200(data interface{}) *WebError {
	return &WebError{
		StatusCode: http.StatusOK,
		Data:       data,
	}
}
