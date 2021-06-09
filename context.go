package gallop

import (
	"github.com/cleango/gallop/infras/errs"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/util/gvalid"
)

type Context struct {
	*gin.Context
}

func (c *Context) ShouldBind(obj interface{}) error {
	if err := c.Context.ShouldBind(obj); err != nil {
		return errs.WithViewErr(err)
	}
	//err := gvalid.New().CheckValue(obj)
	//if err != nil {
	//	return errs.WithViewErr(err)
	//}
	return nil
}
func (c *Context) ShouldBindUri(obj interface{}) error {
	if err := c.Context.ShouldBindUri(obj); err != nil {
		return errs.WithViewErr(err)
	}

	//err := gvalid.New().CheckStruct(obj)
	//if err != nil {
	//	return errs.WithViewErr(err)
	//}
	return nil
}
func (c *Context) ShouldBindHeader(obj interface{}) error {
	if err := c.Context.ShouldBindHeader(obj); err != nil {
		return errs.WithViewErr(err)
	}

	err := gvalid.New().CheckStruct(obj)
	if err != nil {
		return errs.WithViewErr(err)
	}
	return nil
}
