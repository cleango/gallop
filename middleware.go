package gallop

import "github.com/gin-gonic/gin"

type IMiddleware interface {
	OnRequest(*Context) *WebError
	OnResponse(*Context) *WebError
}

//MidFactory 构造函数
func MidFactory(h IMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		webContext := &Context{Context: c}
		err := h.OnRequest(webContext)
		if err != nil {
			c.AbortWithStatusJSON(err.StatusCode, err.Data)
			return
		}
		c.Next()
		err = h.OnResponse(webContext)
		if err != nil {
			c.AbortWithStatusJSON(err.StatusCode, err.Data)
			return
		}
	}
}
