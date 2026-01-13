package amqp

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
)

func bind[T any](fn func(context.Context, T) error) message.NoPublishHandlerFunc {
	return func(msg *message.Message) error {
		var event T
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}
		return fn(msg.Context(), event)
	}
}
