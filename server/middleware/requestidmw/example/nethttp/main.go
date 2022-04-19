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
	// "github.com/gorilla/mux"
)

// raw net/http
func main() {
	// inject middleware specific endpoint
	http.Handle("/hello", requestidmw.HttpMiddleware(handler()))
	runner(http.DefaultServeMux)
}

// using mux router
// func main() {
// 	r := mux.NewRouter()
// 	r.Handle("/", handler())
// 	r.Handle("/hello", handler())

// 	// inject global middleware
// 	http.Handle("/", requestidmw.HttpMiddleware(r))

// 	runner(http.DefaultServeMux)
// }

// your handler
func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request id", requestidmw.GetRequestId(r))
		fmt.Println("check your response header [request-id] is same ?")
	})
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
