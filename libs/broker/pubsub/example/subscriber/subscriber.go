package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/efishery/go-efishery/libs/broker/pubsub"
)

func main() {
	// init pubsub broker using nsq broker
	// on service start
	// pubsub.Init(pubsub.NewNsq(pubsub.NsqConfiguration{
	// 	URL:      "127.0.0.1:4150",
	// 	LogLevel: 1, //  Info
	// }))

	// rabitmq
	// by default this rabbitmq lib is auto reconnect
	config := pubsub.RabitmqConfiguration{
		URL: "amqp://guest:guest@localhost:5672/",
		OnConnectionError: func(oce pubsub.OnConnectionError) {
			// add logic on connection error here
			// e.g hook slack notif, hook whatsapp notif etc
			fmt.Println("handler on connection error here", oce)
		},
	}

	pubsub.Init(pubsub.NewRabitmq(config))

	// subscribe topic
	pubsub.Subscribe(pubsub.SubscriberConfig{
		Topic:   "module-testing",
		Channel: "my-servicename-prod",
	},
		func(c context.Context, m pubsub.Message) error {
			// recommended use parent context if need log proses
			// to make aware timeout&gracefully shutdown all child's and long proses
			// http.NewRequestWithContext(ctx,..,...)

			fmt.Println(m.Payload, m.Type)
			return nil
		})

	fmt.Println("CTRL + C to shutdown")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// close on server shutdown
	//  gracefully shutdown
	pubsub.Close()
}
