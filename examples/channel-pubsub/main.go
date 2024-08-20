package main

import (
	"context"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/idgen"
	"hiboot-messaging/pubsub"
	"time"
)

type FooReceivedEvent struct {
	ID string

	Name string
}

func main() {
	// Example of creating a new Redis client and RedisPubSub instance

	chnPubSub := pubsub.NewChannelPubSub[FooReceivedEvent]()
	ctx := context.Background()
	ch, err := chnPubSub.Subscribe(ctx, "foo")
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
			event := FooReceivedEvent{
				ID:   id,
				Name: "Test Event",
			}
			err = chnPubSub.Publish(ctx, "foo", event)
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
