package amqp

import (
	"github.com/ThreeDotsLabs/watermill/message"
	pubsubadapter "github.com/webitel/im-contact-service/internal/adapter/pubsub"
	"github.com/webitel/im-contact-service/internal/domain/events"
)

// Exchange names
const (
	WebitelGoExchange = "webitel"
)

// Queue names
const (
	DomainDeletedQueue = "im_contacts.domain_delete"
)

func RegisterHandlers(
	router *message.Router,
	subProvider *pubsubadapter.SubscriberProvider,
	h *MessageHandler,
) error {
	subscriptions := []struct {
		topic   string
		queue   string
		handler message.NoPublishHandlerFunc
	}{
		{
			topic:   events.DomainDeletedTopic,
			queue:   DomainDeletedQueue,
			handler: bind(h.OnDomainDeleted),
		},
	}

	for _, s := range subscriptions {
		sub, err := subProvider.Build(
			s.queue,
			WebitelGoExchange,
			s.topic,
		)
		if err != nil {
			return err
		}

		router.AddConsumerHandler(
			s.queue+"_handler",
			s.topic,
			sub,
			s.handler,
		)
	}

	return nil
}
