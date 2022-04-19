package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type rabbitmqConn struct {
	conn *amqp.Connection
	ch   map[string]*amqp.Channel
}
type rabitmqSvc struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	pubConn   *rabbitmqConn
	subConn   *rabbitmqConn
	config    RabitmqConfiguration
	sync.Mutex
}

// RabitmqConfiguration is configuration nsq broker
type RabitmqConfiguration struct {
	// AutoRecconnect is disable or enable
	// automatically reconnect
	// default true
	AutoRecconnect *bool

	// Max retry connection
	// With exponential backoff retry
	// when < 1 no max conn retry == infinite reconnect
	MaxConnRetries int

	// ExponentialMultiplier multiplier
	// recconnection delay, default 10s
	ExponentialMultiplier int

	// OnConnectionError is hook function
	// When connection error
	OnConnectionError func(OnConnectionError)

	// URL is rabitmq url
	URL string

	// If need enable disable service
	Enable bool
}

const (
	exchangeType = "topic"
	keyRequestId = "request-id"
)

// Subscribe is subscribe to rabitmq
func (ss *rabitmqSvc) Subscribe(config SubscriberConfig, h HandlerSubscriber) error {
	// default
	exchangeName := "amq.topic"
	// parse domain topic
	splitTopic := strings.Split(config.Topic, ".")
	if len(splitTopic) > 1 {
		exchangeName = splitTopic[0]
	}

	if config.Timeout < 1 {
		config.Timeout = 15
	}

	if config.Channel == "" {
		config.Channel = filepath.Base(os.Args[0])
	}

	attempt := 0
	var err error
	// prevent open multiple connections on same time
	ss.Lock()
stateConnect:
	// receive value error
	// before call stateConnect
	if err != nil {
		// no auto reconnect
		if !ss.isAutoReconnect() {
			return err
		}

		if ss.config.OnConnectionError != nil {
			ss.config.OnConnectionError(OnConnectionError{
				Error:  err,
				Worker: WorkerSubscriber,
				Meta: map[string]interface{}{
					"topic":    config.Topic,
					"exchange": exchangeName,
				},
			})
		}

		attempt++
		if ss.config.MaxConnRetries > 0 && attempt > ss.config.MaxConnRetries {
			log.Println("max retry exceeded", attempt)
			return err
		}

		// previous delay + current delay
		delay := time.Second * time.Duration(((attempt-1)+attempt)*ss.config.ExponentialMultiplier)
		log.Println(err, "reconnect with delay", delay, "second", ", attempt", attempt)
		<-time.After(delay)

		// clean last connection
		ss.subConn = nil
	}

	if ss.subConn == nil {
		var conn *amqp.Connection
		conn, err = amqp.Dial(ss.config.URL)
		if err != nil {
			goto stateConnect
		}

		ss.subConn = &rabbitmqConn{
			conn: conn,
			ch:   make(map[string]*amqp.Channel),
		}
	}

	var ch *amqp.Channel
	ch, err = ss.subConn.conn.Channel()
	if err != nil {
		goto stateConnect
	}

	// 1 channel per subscriber thread safe
	ss.subConn.ch[uuid.NewString()] = ch

	err = ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		goto stateConnect
	}
	queueName := fmt.Sprintf("%s.%s", config.Topic, config.Channel)
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		goto stateConnect
	}

	err = ch.QueueBind(
		q.Name,       // queue name
		config.Topic, // routing key
		exchangeName, // exchange
		false,
		nil,
	)

	if err != nil {
		goto stateConnect
	}

	msgs, err := ch.Consume(
		q.Name,         // queue
		config.Channel, // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	ss.Unlock()

	go func() {
		defer ch.Close()
		for {
			select {
			case <-ss.ctx.Done():
				return
			case d := <-msgs:

				//  if connection disconnected
				//  and auto reconnect enable, will resubscribe
				if ss.subConn.conn.IsClosed() {

					if ss.isAutoReconnect() {
						ss.Subscribe(config, h)
					} else {
						if ss.config.OnConnectionError != nil {
							ss.config.OnConnectionError(OnConnectionError{
								Error:  fmt.Errorf("connection rabbitmq disconnect"),
								Worker: WorkerSubscriber,
								Meta: map[string]interface{}{
									"topic":    config.Topic,
									"exchange": exchangeName,
									"queue":    queueName,
								},
							})
						}

					}
					return
				}

				msg := Message{}
				var requestId string
				if d.Headers != nil && d.Headers[keyRequestId] != nil {
					requestId = fmt.Sprintf("%s", d.Headers[keyRequestId])
				}

				if err := json.Unmarshal(d.Body, &msg); err != nil {
					if config.OnUnmarshalFailed != nil {
						config.OnUnmarshalFailed(OnUnmarshalFailed{
							Body: d.Body,
							UUID: requestId,
						}, err)
					}

					// no requeue
					d.Reject(false)
					log.Println(err)
					continue
				}

				// if requestId header not
				// empty add it as msg.UUID
				if requestId != "" {
					msg.UUID = requestId
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
					d.Reject(true)
					continue
				}

				d.Ack(false)

			}
		}
	}()

	return nil
}

// Publish is publish to rabitmq broker
// payload & type is required
func (ss *rabitmqSvc) Publish(topic string, msg Message) error {

	// default
	exchangeName := "amq.topic"
	// parse domain topic
	splitTopic := strings.Split(topic, ".")
	if len(splitTopic) > 1 {
		exchangeName = splitTopic[0]
	}

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
	attempt := 0
stateConnect:
	// receive value error
	// before call stateConnect
	if err != nil {
		// no auto reconnect
		if !ss.isAutoReconnect() {
			return err
		}

		// hook  when disconnected
		if ss.config.OnConnectionError != nil {
			ss.config.OnConnectionError(OnConnectionError{
				Error:  err,
				Worker: WorkerPublisher,
				Meta: map[string]interface{}{
					"topic":   topic,
					"uuid":    msg.UUID,
					"payload": msg.Payload,
				},
			})

		}

		attempt++
		if ss.config.MaxConnRetries > 0 && attempt > ss.config.MaxConnRetries {
			log.Println("max retry exceeded", attempt)
			return err
		}
		// previous delay + current delay
		delay := time.Second * time.Duration(((attempt-1)+attempt)*ss.config.ExponentialMultiplier)
		log.Println("reconnect with delay", delay, "second", ", attempt", attempt)
		<-time.After(delay)

		// clean last connection
		ss.pubConn = nil
	}

	if ss.pubConn == nil {
		var conn *amqp.Connection
		conn, err = amqp.Dial(ss.config.URL)
		if err != nil {
			goto stateConnect
		}

		ss.pubConn = &rabbitmqConn{
			conn: conn,
			ch:   make(map[string]*amqp.Channel),
		}
	}

	ch := ss.pubConn.ch["default"]
	if ch == nil {
		ss.pubConn.ch["default"], err = ss.pubConn.conn.Channel()
		if err != nil {
			return err
		}
		ch = ss.pubConn.ch["default"]
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		// do reconnect when connection is closed
		if ss.pubConn.conn.IsClosed() {
			goto stateConnect
		}
		return err
	}

	err = ch.Publish(
		exchangeName, // exchange
		topic,        // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
			Headers: amqp.Table{
				keyRequestId: msg.UUID,
			},
		})

	if err != nil {
		// do reconnect when connection is closed
		if ss.pubConn.conn.IsClosed() {
			goto stateConnect
		}
	}

	return err
}

// Close is close channels and connections
func (ss *rabitmqSvc) Close() error {
	ss.ctxCancel()

	if ss.pubConn != nil {
		for _, ch := range ss.pubConn.ch {
			ch.Close()
		}
		ss.pubConn.conn.Close()
	}

	if ss.subConn != nil {
		for _, ch := range ss.subConn.ch {
			ch.Close()
		}
		ss.subConn.conn.Close()
	}

	return nil
}

func (ss *rabitmqSvc) isAutoReconnect() bool {
	if ss.config.AutoRecconnect == nil {
		return true
	}

	return *ss.config.AutoRecconnect
}

// NewRabitmq new abstraction nsq pubsub
func NewRabitmq(config RabitmqConfiguration) PubSubService {
	if config.ExponentialMultiplier < 1 {
		config.ExponentialMultiplier = 10
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	ss := &rabitmqSvc{
		config:    config,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}
	return ss
}
