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

	err := gvalid.CheckStruct(obj, nil)
	if err != nil {
		return errs.WithViewErr(err)
	}
	return nil
}
