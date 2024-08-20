package main

import (
	"context"
	"encoding/json"
	"github.com/hidevopsio/hiboot-messaging/pubsub"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/idgen"
	"time"
)

type FooReceivedEvent struct {
	ID   string
	Name string
}

func (e *FooReceivedEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *FooReceivedEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}

func main() {
	// Example of creating a new Redis client and RedisPubSub instance

	fooReceivedEventPubSub := pubsub.NewChannelPubSub[*FooReceivedEvent]()
	ctx := context.Background()
	ch, err := fooReceivedEventPubSub.Subscribe(ctx, "foo")
	if err != nil {
		log.Errorf("Error subscribing: %v", err)
		return
	}

	// Simulate publishing a message after 2 seconds
	go func() {
		for {
			time.Sleep(5 * time.Second)
			var id string
			id, err = idgen.NextString()
			event := &FooReceivedEvent{
				ID:   id,
				Name: "Test Event",
			}
			err = fooReceivedEventPubSub.Publish(ctx, "foo", event)
			if err != nil {
				log.Error(err)
			}
		}
	}()

	// Start a goroutine to listen for messages
	for {
		select {
		case msg := <-ch:
			log.Infof("Received message: ID=%s, Name=%s", msg.ID, msg.Name)
		case <-ctx.Done():
			log.Infof("Context cancelled, stopping subscription")
			return
		}
	}
}
