package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
)

type nsqSvc struct {
	producer *nsq.Producer
	consumer []*nsq.Consumer
	config   NsqConfiguration
}

// NsqConfiguration is configuration nsq broker
type NsqConfiguration struct {
	// URL is nsqd url
	URL string
	// log level default debug 0-3
	LogLevel int
	// If need enable disable service
	Enable bool
}

// Subscribe is subscribe to nsqd
func (ss *nsqSvc) Subscribe(config SubscriberConfig, h HandlerSubscriber) error {
	nsqconfig := nsq.NewConfig()
	if config.Timeout < 1 {
		config.Timeout = 15
	}

	if config.Channel == "" {
		config.Channel = filepath.Base(os.Args[0])
	}

	q, err := nsq.NewConsumer(config.Topic, config.Channel, nsqconfig)
	if err != nil {
		return err
	}

	q.SetLoggerLevel(nsq.LogLevel(ss.config.LogLevel))
	q.AddHandler(nsq.HandlerFunc(func(nsqMessage *nsq.Message) error {
		msg := Message{}
		if err := json.Unmarshal(nsqMessage.Body, &msg); err != nil {
			config.OnUnmarshalFailed(OnUnmarshalFailed{
				Body: nsqMessage.Body,
			}, err)

			return err
		}

		// create value context timeout
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Second)
		defer cancelFunc()
		errChan := make(chan error)
		go func(ctx context.Context, msg Message) {
			errChan <- h(ctx, msg)
		}(ctx, msg)
		go func(ctx context.Context) {
			// wait timeout
			for {
				select {
				case <-ctx.Done():
					errChan <- ErrorConsumeTimeout
					return
				default:
					<-time.After(time.Duration(config.Timeout) * time.Second)
				}
			}
		}(ctx)

		// if error not nil wil requeue message
		err = <-errChan
		if err != nil {
			log.Println(err, msg.UUID)
		}

		return err
	}))

	err = q.ConnectToNSQD(ss.config.URL)
	if err != nil {
		return err
	}

	ss.consumer = append(ss.consumer, q)

	return nil
}

// Publish is publish to nsqd broker
// payload & type is required
func (ss *nsqSvc) Publish(topic string, msg Message) error {
	// validation
	if msg.Payload == nil {
		return fmt.Errorf("Payload required")
	}

	// default now
	if msg.Timestamp < 1 {
		msg.Timestamp = int(toUnixMilliSecond(time.Now()))
	}

	// default uuid
	if msg.UUID == "" {
		msg.UUID = uuid.New().String()
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if ss.producer == nil {
		config := nsq.NewConfig()
		ss.producer, err = nsq.NewProducer(ss.config.URL, config)
		if err != nil {
			return err
		}
		ss.producer.SetLogger(log.New(os.Stderr, "", log.Flags()), nsq.LogLevel(ss.config.LogLevel))
	}

	return ss.producer.Publish(topic, data)
}

func (ss *nsqSvc) Close() error {
	if ss.consumer != nil {
		for _, v := range ss.consumer {
			v.Stop()
		}
	}

	if ss.producer != nil {
		ss.producer.Stop()
	}

	return nil
}

// NewNsq new abstraction nsq pubsub
func NewNsq(config NsqConfiguration) PubSubService {
	ss := &nsqSvc{
		config: config,
	}
	return ss
}
