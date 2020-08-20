package gallop

import "net/http"

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
