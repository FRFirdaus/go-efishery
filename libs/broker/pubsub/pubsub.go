package pubsub

import (
	"context"
	"errors"
)

var (
	//ErrorConsumeTimeout if flag handler consume timeout and requeue
	ErrorConsumeTimeout = errors.New("Consume Timeout")
)

// SubscriberConfig subscriber high-level configuration
type SubscriberConfig struct {
	// Timeout in second
	Timeout int
	// Topic will subscribe
	Topic string
	// Channel is queue label
	// default using service name
	Channel string

	// HookOnUnmarshalFailed is hook function
	// When error unmarshal Message struct
	OnUnmarshalFailed func(OnUnmarshalFailed, error)
}

// PubSubService abstraction for publish and subscribe broker
type PubSubService interface {
	Publish(topic string, msg Message) error
	Subscribe(config SubscriberConfig, h HandlerSubscriber) error
	Close() error
}

// HandlerSubscriber function for subscriber
type HandlerSubscriber func(context.Context, Message) error

// Message is message format
type Message struct {
	// UUID is message id
	// Default uuid()
	UUID string `json:"uuid,omitempty"`
	// Timestamp unix on queue publish
	// Default now()
	Timestamp int `json:"timestamp,omitempty"`
	// Type can event type or etc
	Type string `json:"type,omitempty"`
	// Payload is data
	Payload interface{} `json:"payload,omitempty"`
}

// OnUnmarshalFailed structure data on unmarshal errors
type OnUnmarshalFailed struct {
	// body is original data from broker
	Body []byte

	// UUID is request-id or message.uuid
	UUID string `json:"uuid,omitempty"`
}

// WorkerName is identifier is subscriber or publisher
type WorkerName string

const (
	WorkerPublisher  WorkerName = "publisher"
	WorkerSubscriber WorkerName = "subscriber"
)

// OnConnectionError is struct data on connection error
type OnConnectionError struct {
	Worker WorkerName
	Error  error
	Meta   map[string]interface{}
}
