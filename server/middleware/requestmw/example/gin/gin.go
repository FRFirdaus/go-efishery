package main

import (
	"bitbucket.org/efishery/go-efishery/server/middleware/requestmw"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// inject global middleware
	rw := requestmw.Init()
	router.Use(rw.InitGinReqMW())
}
