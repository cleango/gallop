package controller

import (
	"fmt"
	"github.com/cleango/gallop"
	"github.com/cleango/gallop/example/app/config"
	"github.com/cleango/gallop/logger"
	"github.com/gin-gonic/gin"
)

type HelloController struct {
	Demo *config.Demo  `inject:""`
	Demo1 *config.Demo `inject:"demo"`
	Cfg *config.Configuration `inject:""`
}

func (ctr *HelloController) Hello(c *gallop.Context) string {
	filed:=logger.LogField{}
	filed["req_id"]="xxxxxxxxxxxxx"
	logger.Info("1233",filed)
	//输出： {"level":"INFO","ts":"2020-08-28 14:26:32","func":"controller/hello.go:20","msg":"1233","req_id":"xxxxxxxxxxxxx"}
	return fmt.Sprint(ctr.Demo,ctr.Demo1,ctr.Cfg.B.C)
}

func (ctr *HelloController) Json(c *gallop.Context) gallop.Json{
	return gin.H{
		"hello":"world",
	}
}