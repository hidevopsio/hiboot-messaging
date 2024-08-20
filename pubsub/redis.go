package pubsub

import (
	"context"
	"encoding/json"
	"github.com/hidevopsio/hiboot-data/starter/redis"
	"github.com/hidevopsio/hiboot/pkg/log"
	goredis "github.com/redis/go-redis/v9"
	"reflect"
)

type RedisPubSub[T any] struct {
	pubSub *goredis.PubSub
}

func NewRedisPubSub[T any]() *RedisPubSub[T] {
	return &RedisPubSub[T]{}
}

func (p *RedisPubSub[T]) Publish(ctx context.Context, client *redis.Client, topic string, msg T) error {
	return client.Publish(ctx, topic, msg).Err()
}

func (p *RedisPubSub[T]) Subscribe(ctx context.Context, client *redis.Client, topic string) (<-chan T, error) {
	// Change the channel type to *T
	ch := make(chan T)
	p.pubSub = client.Subscribe(ctx, topic)

	// Check for errors during subscription
	_, err := p.pubSub.Receive(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		for msg := range p.pubSub.Channel() {
			// Create a new instance of the pointer type T
			var v T // This creates a pointer to a new instance of type T
			tType := reflect.TypeOf(v)

			if tType.Kind() == reflect.Ptr {
				// If T is already a pointer, return the zero value (which is nil) or a new instance
				v = reflect.New(tType.Elem()).Interface().(T)
			} else if tType.Kind() == reflect.Struct {
				// If T is a struct, create a new pointer to a struct and return it
				v = reflect.New(tType).Elem().Interface().(T)
			}

			// Unmarshal into the pointer
			e := json.Unmarshal([]byte(msg.Payload), v)
			if e != nil {
				log.Error(e)
			} else {
				// Send the pointer directly to the channel
				ch <- v
			}
		}
		close(ch)
	}()

	return ch, nil
}

func (p *RedisPubSub[T]) Unsubscribe(ctx context.Context, topic string) error {
	if p.pubSub == nil {
		return nil
	}
	return p.pubSub.Unsubscribe(ctx, topic)
}
