package pubsub

import (
	"context"
	"sync"
)

type ChannelPubSub[T any] struct {
	mu     sync.RWMutex
	topics map[string][]chan T
}

func NewChannelPubSub[T any]() *ChannelPubSub[T] {
	return &ChannelPubSub[T]{topics: make(map[string][]chan T)}
}

func (p *ChannelPubSub[T]) Publish(ctx context.Context, topic string, msg T) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if subs, ok := p.topics[topic]; ok {
		for _, ch := range subs {
			ch <- msg
		}
	}
	return nil
}

func (p *ChannelPubSub[T]) Subscribe(ctx context.Context, topic string) (<-chan T, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	ch := make(chan T)
	p.topics[topic] = append(p.topics[topic], ch)

	return ch, nil
}

func (p *ChannelPubSub[T]) Unsubscribe(ctx context.Context, topic string, ch <-chan T) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if subs, ok := p.topics[topic]; ok {
		for i, subCh := range subs {
			if subCh == ch {
				close(subCh)
				p.topics[topic] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}

	return nil
}
