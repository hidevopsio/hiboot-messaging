package pubsub

import (
	"context"
	"sync"
)

type ChannelPubSub[T any] struct {
	mu     sync.RWMutex
	topics map[string]chan T
}

func NewChannelPubSub[T any]() *ChannelPubSub[T] {
	return &ChannelPubSub[T]{topics: make(map[string]chan T)}
}

func (p *ChannelPubSub[T]) Publish(ctx context.Context, topic string, msg T) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if ch, ok := p.topics[topic]; ok {
		ch <- msg
	}
	return nil
}

func (p *ChannelPubSub[T]) Subscribe(ctx context.Context, topic string) (<-chan T, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.topics[topic]; !ok {
		p.topics[topic] = make(chan T)
	}

	return p.topics[topic], nil
}

func (p *ChannelPubSub[T]) Unsubscribe(ctx context.Context, topic string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ch, ok := p.topics[topic]; ok {
		close(ch)
		delete(p.topics, topic)
	}

	return nil
}
