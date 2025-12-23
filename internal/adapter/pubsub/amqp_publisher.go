package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type WatermillPublisher struct {
	publisher message.Publisher
}

func NewWatermillPublisher(pub message.Publisher) *WatermillPublisher {
	return &WatermillPublisher{publisher: pub}
}

func (p *WatermillPublisher) Publish(ctx context.Context, topic string, event any) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	msg.SetContext(ctx)

	if err := p.publisher.Publish(topic, msg); err != nil {
		return fmt.Errorf("watermill publish: %w", err)
	}
	return nil
}
