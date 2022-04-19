package requestidmw

import (
	"github.com/gin-gonic/gin"
)

// GinMiddleware is inject request id and response id
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		HttpMiddleware(nil).ServeHTTP(c.Writer, c.Request)
	}
}

// GetGinRequestId is get request id from request header
func GetGinRequestId(c *gin.Context) string {
	return GetRequestId(c.Request)
}
