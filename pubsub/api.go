package pubsub

import "context"

// PubSub is a generic interface for Pub/Sub systems
type PubSub[T any] interface {
	Publish(ctx context.Context, topic string, msg T) error
	Subscribe(ctx context.Context, topic string) (<-chan T, error)
	Unsubscribe(ctx context.Context, topic string) error
}
