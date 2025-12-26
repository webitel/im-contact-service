package amqp

import (
	"github.com/ThreeDotsLabs/watermill/message"
	pubsubadapter "github.com/webitel/im-contact-service/internal/adapter/pubsub"
)

func RegisterHandlers(router *message.Router, subProvider *pubsubadapter.SubscriberProvider, h *MessageHandler) error {
	subscriptions := []struct {
		topic      string
		queueGroup string
		handler    message.NoPublishHandlerFunc
	}{
		{
			topic:      "domain.deleted",
			queueGroup: "wbt.directory.domain_deleted",
			handler:    bind(h.OnDomainDeleted),
		},
	}

	for _, s := range subscriptions {
		sub, err := subProvider.Build(
			s.queueGroup,
			"webitel.admin",
			s.topic,
		)
		if err != nil {
			return err
		}

		router.AddConsumerHandler(
			s.queueGroup+"_handler",
			s.topic,
			sub,
			s.handler,
		)
	}

	return nil
}
