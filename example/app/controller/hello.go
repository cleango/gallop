package controller

import (
	"fmt"
	"github.com/cleango/gallop"
	"github.com/cleango/gallop/example/app/config"
	"github.com/gin-gonic/gin"
)

type HelloController struct {
	Demo *config.Demo  `inject:""`
	Demo1 *config.Demo `inject:"demo"`
}

func (ctr *HelloController) Hello(c *gallop.Context) string {
	return fmt.Sprint(ctr.Demo,ctr.Demo1)
}

func (ctr *HelloController) Json(c *gallop.Context) gallop.Json{
	return gin.H{
		"hello":"world",
	}
}