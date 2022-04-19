package requestmw

import "github.com/gin-gonic/gin"

func (rw *requestmw) InitGinReqMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		buildRequestMW(nil).ServeHTTP(c.Writer, c.Request)
	}
}
