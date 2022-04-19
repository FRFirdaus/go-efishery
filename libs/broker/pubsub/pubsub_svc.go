package pubsub

import (
	"fmt"
	"sync"
)

var _svcBroker PubSubService
var _mutex sync.Mutex

var (
	errInitSvc = fmt.Errorf("Please init broker first")
)

// Publish is high level function of publish message to broker
func Publish(topic string, msg Message) error {
	_mutex.Lock()
	defer _mutex.Unlock()

	if _svcBroker == nil {
		return errInitSvc
	}
	return _svcBroker.Publish(topic, msg)
}

// Subscribe is high level function of subscribe to broker
func Subscribe(config SubscriberConfig, h HandlerSubscriber) error {
	_mutex.Lock()
	defer _mutex.Unlock()

	if _svcBroker == nil {
		return errInitSvc
	}

	return _svcBroker.Subscribe(config, h)
}

// Close is high level gracefully close connection to broker
func Close() error {
	_mutex.Lock()
	defer _mutex.Unlock()

	if _svcBroker == nil {
		return errInitSvc
	}

	if err := _svcBroker.Close(); err != nil {
		return err
	}

	_svcBroker = nil

	return nil
}

// Init is initialize broker service
func Init(svc PubSubService) error {
	_mutex.Lock()
	defer _mutex.Unlock()

	if _svcBroker != nil {
		return fmt.Errorf("Already initialized")
	}
	_svcBroker = svc
	return nil
}
