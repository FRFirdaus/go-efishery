package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bitbucket.org/efishery/go-efishery/server/middleware/requestidmw"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// inject global middleware
	router.Use(requestidmw.GinMiddleware())

	router.GET("/hello", handler)
	router.GET("/hello/me", handler)

	// if only need specific endpoint
	// router.GET("/hello/word",requestidmw.GinMiddleware(), handler)
	runner(router)
}

// your handler
func handler(ctx *gin.Context) {
	fmt.Println("request id", requestidmw.GetGinRequestId(ctx))
	fmt.Println("check your response header [request-id] is same ?")
}

func runner(router http.Handler) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Println("server listten on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
