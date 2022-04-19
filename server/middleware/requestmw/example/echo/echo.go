package main

import (
	"bitbucket.org/efishery/go-efishery/server/middleware/requestmw"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	// inject global middleware
	rw := requestmw.Init()
	e.Use(rw.InitEchoReqMW())
}
