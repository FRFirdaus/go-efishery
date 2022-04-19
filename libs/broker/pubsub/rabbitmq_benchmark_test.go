package pubsub

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func doBenchmarkRabitmq(subscriber, publisher, msg int, b *testing.B) {
	for n := 0; n < b.N; n++ {

		topic := b.Name()
		wg := &sync.WaitGroup{}
		for i := 0; i < subscriber; i++ {
			go func() {

			}()
			svc := NewRabitmq(RabitmqConfiguration{
				URL: "amqp://guest:guest@riset.public.efishery.com:5672/",
			})
			defer svc.Close()

			err := svc.Subscribe(SubscriberConfig{
				Topic:   topic,
				Channel: fmt.Sprintf("%s_%d", "Subscriber", i),
			}, func(ctx context.Context, m Message) error {
				wg.Done()
				return nil
			})
			if err != nil {
				b.Error(err)
				return
			}
		}

		for i := 0; i < publisher; i++ {
			wg.Add(1 * subscriber * msg)
			go func() {
				svc := NewRabitmq(RabitmqConfiguration{
					URL: "amqp://guest:guest@riset.public.efishery.com:5672/",
				})
				defer svc.Close()
				for i := 0; i < msg; i++ {
					if err := svc.Publish(topic, Message{UUID: "id1", Type: "test", Payload: "any data"}); err != nil {
						b.Error(err)
						return
					}
				}
			}()
		}

		wg.Wait()
	}
}

func BenchmarkRabitmqSubscriber1Pub1(b *testing.B) { doBenchmarkRabitmq(1, 1, 100, b) }
