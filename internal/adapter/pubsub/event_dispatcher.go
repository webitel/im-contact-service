package adapter

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/webitel/webitel-go-kit/pkg/errors"

	"github.com/webitel/im-contact-service/internal/domain/events"
)

type EventDispatcher struct {
	publisher message.Publisher
}

func NewEventDispatcher(pub message.Publisher) *EventDispatcher {
	return &EventDispatcher{publisher: pub}
}

func (d *EventDispatcher) Publish(ctx context.Context, event events.Event) error {
	if event == nil {
		return errors.InvalidArgument("received nil pointer event", errors.WithID("pubsub.event_dispatcher.publish"))
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return errors.Internal("marshaling event", errors.WithCause(err), errors.WithID("pubsub.event_dispatcher.publish"))
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	msg.SetContext(ctx)

	topic := event.Topic()

	if err := d.publisher.Publish(topic, msg); err != nil {
		return errors.Internal("publishing event to topic", errors.WithCause(err), errors.WithID("pubsub.event_dispatcher.publish"), errors.WithValue("topic", topic))
	}

	return nil
}
