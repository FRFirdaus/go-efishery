package requestmw

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	e4 "github.com/labstack/echo/v4"
)

type RequestMW interface {
	InitEchoReqMWV4() e4.MiddlewareFunc
	InitEchoReqMW() echo.MiddlewareFunc

	InitGinReqMW() gin.HandlerFunc
}

type requestmw struct {
}

func Init() RequestMW {
	return &requestmw{}
}

func buildRequestMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		scheme := "HTTP"
		if r.TLS != nil {
			scheme = "HTTPS"
		}

		reqID := r.Header.Get("x-request-id")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		addrIp := r.Header.Get("x-forwarded-for")
		if addrIp == "" {
			addrIp = r.RemoteAddr
		}

		ctx := context.WithValue(r.Context(), "x-request-method", r.Method)
		ctx = context.WithValue(ctx, "x-request-scheme", scheme)
		ctx = context.WithValue(ctx, "x-forwarded-for", addrIp)
		ctx = context.WithValue(ctx, "x-request-id", reqID)

		// add response x-request-id
		w.Header().Add("x-request-id", reqID)

		// later we can add tracer id for apm here too

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
