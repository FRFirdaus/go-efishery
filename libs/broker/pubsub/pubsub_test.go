package pubsub

import (
	"context"
	"testing"
)

func TestPubSubSvc(t *testing.T) {

	svc := NewNsq(NsqConfiguration{
		URL: "riset.public.efishery.com:4150",
	})

	Init(svc)

	closed := make(chan (bool), 2)
	topic := "TestPubSubSvc_#ephemeral"

	err := Subscribe(SubscriberConfig{
		Topic:   topic,
		Channel: "TestPubSubSvc_#ephemeral",
	}, func(ctx context.Context, m Message) error {
		t.Log(m)
		closed <- true
		return nil
	})

	if err != nil {
		t.Error(err)
		return
	}

	if err := Publish(topic, Message{UUID: "id1", Type: "test", Payload: "any data"}); err != nil {
		t.Error(err)
		return
	}

	<-closed

	if err := Close(); err != nil {
		t.Error(err)
		return
	}
}
