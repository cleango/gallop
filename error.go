package gallop

import "net/http"

func Throw(err error) {
	if err != nil {
		panic(err)
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
