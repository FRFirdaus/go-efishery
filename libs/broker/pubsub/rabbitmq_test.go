package pubsub

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRabitmqPubSub(t *testing.T) {

	svc := NewRabitmq(RabitmqConfiguration{
		URL: "amqp://guest:guest@riset.public.efishery.com:5672/",
	})
	defer svc.Close()

	closed := make(chan (bool))
	topic := "module-testing"

	err := svc.Subscribe(SubscriberConfig{
		Topic: topic,
	}, func(ctx context.Context, m Message) error {
		t.Log(m)
		closed <- true
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	if err := svc.Publish(topic, Message{UUID: "id1", Type: "test", Payload: "any data"}); err != nil {
		t.Error(err)
		return
	}
	<-closed
}

func TestErrRabitmqPubSub(t *testing.T) {

	svc := NewRabitmq(RabitmqConfiguration{
		URL: "amqp://guest:guest@riset.public.efishery.com:5672/",
	})

	closed := make(chan (bool))
	topic := "module-testing"
	flagErr := false

	err := svc.Subscribe(SubscriberConfig{
		Topic: topic,
	}, func(ctx context.Context, m Message) error {
		t.Log(m)

		if !flagErr {
			flagErr = true
			return fmt.Errorf("check requeue is success ")
		}
		closed <- true

		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err := svc.Publish(topic, Message{UUID: "id1", Type: "test", Payload: "any data"}); err != nil {
		t.Error(err)
		return
	}

	<-closed

	if err := svc.Close(); err != nil {
		t.Error(err)
		return
	}
}

func TestRabitmqPubSubTimeout(t *testing.T) {

	svc := NewRabitmq(RabitmqConfiguration{
		URL: "amqp://guest:guest@riset.public.efishery.com:5672/",
	})
	defer svc.Close()
	var wg sync.WaitGroup
	total := 2
	wg.Add(total)
	topic := "module-testing"

	err := svc.Subscribe(SubscriberConfig{
		Topic:   topic,
		Timeout: 2,
	}, func(ctx context.Context, m Message) error {
		fmt.Println(m)
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			default:
				// long process until timeout
				<-time.After(1 * time.Second)
			}
		}
		wg.Done()
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err := svc.Publish(topic, Message{UUID: "id1", Type: "test", Payload: "any data"}); err != nil {
		t.Error(err)
		return
	}
	wg.Wait()
}

func TestRabitmqWhenSubscriberDown(t *testing.T) {

	svc := NewRabitmq(RabitmqConfiguration{
		URL: "amqp://guest:guest@riset.public.efishery.com:5672/",
	})

	var wg sync.WaitGroup
	total := 2
	wg.Add(total)
	topic := "module-testing"

	err := svc.Subscribe(SubscriberConfig{
		Topic:   topic,
		Channel: "channel-1",
	}, func(ctx context.Context, m Message) error {
		fmt.Println("ch1", m)
		wg.Done()
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Subscribe(SubscriberConfig{
		Topic:   topic,
		Channel: "channel-2",
	}, func(ctx context.Context, m Message) error {
		fmt.Println("ch2", m)
		wg.Done()
		return nil
	})

	if err != nil {
		t.Error(err)
		return
	}

	if err := svc.Publish(topic, Message{UUID: "id1", Type: "test", Payload: "any data"}); err != nil {
		t.Error(err)
		return
	}
	wg.Wait()
	// close subsribeer then publish topic
	svc.Close()
	wg.Add(total)
	// resubscribe with same topic and channel
	svc = NewRabitmq(RabitmqConfiguration{
		URL: "amqp://guest:guest@riset.public.efishery.com:5672/",
	})
	defer svc.Close()
	if err := svc.Publish(topic, Message{UUID: "id1", Type: "test", Payload: "any data"}); err != nil {
		t.Error(err)
		return
	}

	// resubscribe with same topic and channel
	err = svc.Subscribe(SubscriberConfig{
		Topic:   topic,
		Channel: "channel-2",
	}, func(ctx context.Context, m Message) error {
		fmt.Println("ch2", m)
		wg.Done()
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}

	// resubscribe with same topic and channel
	err = svc.Subscribe(SubscriberConfig{
		Topic:   topic,
		Channel: "channel-1",
	}, func(ctx context.Context, m Message) error {
		fmt.Println("ch1", m)
		wg.Done()
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	wg.Wait()
}
