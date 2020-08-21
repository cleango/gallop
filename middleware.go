package gallop

import "github.com/gin-gonic/gin"

type IMidHandler func(*Context)

//MidFactory 构造函数
func MidFactory(h IMidHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		webContext := &Context{Context: c}
		h(webContext)
	}
}
