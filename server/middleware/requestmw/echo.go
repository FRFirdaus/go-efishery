package requestmw

import (
	"github.com/labstack/echo"
	e4 "github.com/labstack/echo/v4"
)

func (rw *requestmw) InitEchoReqMWV4() e4.MiddlewareFunc {
	return e4.WrapMiddleware(buildRequestMW)
}

func (rw *requestmw) InitEchoReqMW() echo.MiddlewareFunc {
	return echo.WrapMiddleware(buildRequestMW)
}
