package pubsub

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// Test one sequence subscribe and publish
// ensure publish is broadcasted
// ensure subscriber receive message
func TestNsqPubSub(t *testing.T) {

	svc := NewNsq(NsqConfiguration{
		URL: "riset.public.efishery.com:4150",
	})

	closed := make(chan (bool))
	topic := "TestNsqPubSub_#ephemeral"

	err := svc.Subscribe(SubscriberConfig{
		Topic:   topic,
		Channel: "TestNsqPubSub_#ephemeral",
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

	if err := svc.Close(); err != nil {
		t.Error(err)
		return
	}
}

// Test return error on subscriber handler
// ensure requeue on error
// ensure re consume after requeue
func TestErrNsqPubSub(t *testing.T) {

	svc := NewNsq(NsqConfiguration{
		URL: "riset.public.efishery.com:4150",
	})

	closed := make(chan (bool))
	topic := "TestErrNsqPubSub_#ephemeral"
	flagErr := false

	err := svc.Subscribe(SubscriberConfig{
		Topic:   topic,
		Channel: "TestErrNsqPubSub_#ephemeral",
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

// Test using timeout subscriber handler
// Ensure return error if proses break the timeout
// on timeout will requeue and re consume
func TestNsqPubSubTimeout(t *testing.T) {

	svc := NewNsq(NsqConfiguration{
		URL: "riset.public.efishery.com:4150",
	})
	defer svc.Close()
	var wg sync.WaitGroup
	total := 2
	wg.Add(total)
	topic := "TestNsqPubSubTimeout_#ephemeral"

	err := svc.Subscribe(SubscriberConfig{
		Topic:   topic,
		Channel: "TestNsqPubSubTimeout_#ephemeral",
		Timeout: 2,
	}, func(ctx context.Context, m Message) error {
		fmt.Println(m)
	loop:
		for {
			select {
			case <-ctx.Done():
				fmt.Println("break")
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

//  Check when subscriber down will always receive message
func TestNsqWhenSubscriberDown(t *testing.T) {

	svc := NewNsq(NsqConfiguration{
		URL: "riset.public.efishery.com:4150",
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
	svc = NewNsq(NsqConfiguration{
		URL: "riset.public.efishery.com:4150",
	})
	if err := svc.Publish(topic, Message{UUID: "id1", Type: "test", Payload: "any data"}); err != nil {
		t.Error(err)
		return
	}
	defer svc.Close()
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
