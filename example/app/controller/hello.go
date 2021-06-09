package controller

import (
	"fmt"
	"github.com/cleango/gallop"
	"github.com/cleango/gallop/example/app/config"
	"github.com/cleango/gallop/example/app/job"
	"github.com/cleango/gallop/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HelloController struct {
	Demo  *config.Demo          `inject:""`
	Demo1 *config.Demo          `inject:"demo"`
	Cfg   *config.Configuration `inject:""`
}
type User struct {
	Name string `json:"name" form:"name" v:"required#请输入用户姓名"`
	Type int    `json:"type" form:"type" v:"required#请选择用户类型"`
}

func (ctr *HelloController) Valid(c *gallop.Context) string {
	var req User

	//req:=cmap.Map{}
	if err := c.ShouldBind(&req); err != nil {
		gallop.Throw(err)
	}
	return "success"
}
func (ctr *HelloController) Hello(c *gallop.Context) string {
	filed := logger.LogField{}
	filed["req_id"] = "xxxxxxxxxxxxx"
	logger.Info("1233", filed)
	gallop.AddJob("@every 3s", &job.Hello1{Demo: ctr.Demo})
	//输出： {"level":"INFO","ts":"2020-08-28 14:26:32","func":"controller/hello.go:20","msg":"1233","req_id":"xxxxxxxxxxxxx"}
	return fmt.Sprint(ctr.Demo, ctr.Demo1, ctr.Cfg.B.C)
}

func (ctr *HelloController) Json(c *gallop.Context) gallop.Json {
	return gin.H{
		"hello": "world",
	}
}

func (ctr *HelloController) File(c *gallop.Context) gallop.File {
	head := http.Header{}
	head.Set("Content-Type", "application/octet-stream")
	head.Set("Content-Disposition", "attachment; filename="+"Workbook.xls")
	head.Set("Content-Transfer-Encoding", "binary")
	return gallop.File{
		Data:        []byte("123, 123, 12312, 3123"),
		ContentType: "application/octet-stream",
		FileName:    "1.csv",
	}
}
