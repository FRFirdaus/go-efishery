# pubsub 

pubsub is high level abstraction of publisher and subscriber protocol.
goal is can use any broker vendor, currently only support [nsq](https://nsq.io/)

## Example

- [Publisher](https://bitbucket.org/efishery/go-efishery/src/master/libs/broker/pubsub/example/publisher/publisher.go)
- [Subscriber](https://bitbucket.org/efishery/go-efishery/src/master/libs/broker/pubsub/example/subscriber/subscriber.go)

## Base knowledge

### Init
Init is initialize broker service
e.g
Init broker pubsub service using nsq broker 
```go
import	"bitbucket.org/efishery/go-efishery/libs/broker/pubsub"

// init pubsub broker using nsq broker
// on service start
pubsub.Init(pubsub.NewNsq(pubsub.NsqConfiguration{
    URL:      "127.0.0.1:4150",
    LogLevel: 1, //  Info
}))

```

Structure:
- [Nsq Configuration](https://bitbucket.org/efishery/go-efishery/src/master/libs/broker/pubsub/nsq.go#libs/broker/pubsub/nsq.go-22)

### Close

Close is high level gracefully close connection to broker after init
```go
import	"bitbucket.org/efishery/go-efishery/libs/broker/pubsub"

// close on server shutdown
// gracefully shutdown
pubsub.Close()
```


### Publish

Publish is high level function of publish message to broker.
Need call `Init` first.

```go

import	"bitbucket.org/efishery/go-efishery/libs/broker/pubsub"

// publish messsage
// no need serialize close after publish to prevent open close connection
	pubsub.Publish("module-testing",
		pubsub.Message{Type: "order",
			Payload: map[string]interface{}{
				"setdata1": 1,
				"setdata2": "data",
			}},
	)
```

Structure:
-  [Message](https://bitbucket.org/efishery/go-efishery/src/master/libs/broker/pubsub/pubsub.go#libs/broker/pubsub/pubsub.go-34)


### Subscribe

Subscribe is high level function of subscribe to broker.
Need call `Init` first.

```go
import	"bitbucket.org/efishery/go-efishery/libs/broker/pubsub"
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
```

Structure:
-  [SubscriberConfig](https://bitbucket.org/efishery/go-efishery/src/master/libs/broker/pubsub/pubsub.go#libs/broker/pubsub/pubsub.go-13)
-  [Message](https://bitbucket.org/efishery/go-efishery/src/master/libs/broker/pubsub/pubsub.go#libs/broker/pubsub/pubsub.go-34)
