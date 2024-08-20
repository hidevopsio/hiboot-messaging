package main

import (
	"context"
	"encoding/json"
	"github.com/hidevopsio/hiboot-data/starter/redis"
	"github.com/hidevopsio/hiboot-messaging/pubsub"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/locale"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"github.com/hidevopsio/hiboot/pkg/utils/idgen"
)

type FooReceivedEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Status string `json:"status"`
}

func (e *FooReceivedEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *FooReceivedEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}

type DemoController struct {
	at.RestController
	at.RequestMapping `value:"/"`

	pubsub *pubsub.RedisPubSub[*FooReceivedEvent]
}

func newDemoController(pubsub *pubsub.RedisPubSub[*FooReceivedEvent]) *DemoController {
	return &DemoController{pubsub: pubsub}
}

func init() {
	app.Register(newDemoController, pubsub.NewRedisPubSub[*FooReceivedEvent])
}

func (c *DemoController) GetSub(redisClient *redis.Client) (err error) {
	go func() {
		for {
			ctx := context.Background()
			ch, _ := c.pubsub.Subscribe(ctx, redisClient, "foo")
			// Start a goroutine to listen for messages
			select {
			case msg := <-ch:
				log.Infof("Received message: ID=%v, Name=%v, Status: %v", msg.ID, msg.Name, msg.Status)
			case <-ctx.Done():
				log.Infof("Context cancelled, stopping subscription")
				return
			}
		}
	}()
	return
}

type PublishEvent struct {
	at.ResponseBody

	FooReceivedEvent
}

func (c *DemoController) GetPub(redisClient *redis.Client) (response *PublishEvent, err error) {
	var id string
	id, err = idgen.NextString()
	event := &FooReceivedEvent{
		ID:     id,
		Name:   "Test Event",
		Status: "test",
	}

	err = c.pubsub.Publish(context.Background(), redisClient, "foo", event)
	if err != nil {
		log.Error(err)
	}
	response = new(PublishEvent)
	response.ID = event.ID
	response.Name = event.Name
	response.Status = event.Status
	return
}

func main() {
	web.NewApplication().
		SetProperty(app.ProfilesInclude,
			actuator.Profile,
			locale.Profile,
			redis.Profile,
			logging.Profile).
		Run()
}
