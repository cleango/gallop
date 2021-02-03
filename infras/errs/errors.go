package errs

import (
	"errors"
	"fmt"
)

var (
	ViewError = errors.New("")
	TokenErr  = errors.New("登录令牌已失效，请重新登录")
	BizErr    = errors.New("")
	DBErr     = errors.New("")
)

func WithMessage(err error, message string) error {
	if err != nil {
		return fmt.Errorf("%w%s", err, message)
	}
	return nil
}
func WithError(err, err1 error) error {
	if err1 != nil {
		return fmt.Errorf("%w %s", err, err1.Error())
	}
	return nil
}
func WithViewErr(msg error) error {
	return WithError(ViewError, msg)
}
func WithViewErrf(layout string, msg ...interface{}) error {
	return WithMessage(ViewError, fmt.Sprintf(layout, msg...))
}
func WithTokenErr(msg error) error {
	return WithError(TokenErr, msg)
}
func WithTokenErrf(layout string, msg ...interface{}) error {
	return WithMessage(TokenErr, fmt.Sprintf(layout, msg...))
}
func WithBizErr(msg error) error {
	return WithError(BizErr, msg)
}
func WithBizErrf(layout string, msg ...interface{}) error {
	return WithMessage(BizErr, fmt.Sprintf(layout, msg...))
}
func WithDBErr(msg error) error {
	return WithError(DBErr, msg)
}
func WithDBErrf(layout string, msg ...interface{}) error {
	return WithMessage(DBErr, fmt.Sprintf(layout, msg...))
}
