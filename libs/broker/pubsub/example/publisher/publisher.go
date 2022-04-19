package main

import (
	"fmt"

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

	// publish messsage
	// no need serialize close after publish to prevent open close connection
	pubsub.Publish("module-testing",
		pubsub.Message{Type: "order",
			Payload: map[string]interface{}{
				"setdata1": 1,
				"setdata2": "data",
			}},
	)

	// close on server shutdown
	//  gracefully shutdown
	pubsub.Close()
}
