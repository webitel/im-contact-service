package amqp

import (
	"context"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	pubsubadapter "github.com/webitel/im-contact-service/internal/adapter/pubsub"
	"github.com/webitel/im-contact-service/internal/domain/events"
	"go.uber.org/fx"
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


func NewWatermillRouter(lc fx.Lifecycle, logger *slog.Logger) (*message.Router, error) {
    router, err := message.NewRouter(message.RouterConfig{}, watermill.NewSlogLogger(logger))
    if err != nil {
        return nil, err
    }

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            go func() {
                if err := router.Run(ctx); err != nil {
                    logger.Error("watermill router run error", "err", err)
                }
            }()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            return router.Close()
        },
    })

    return router, nil
}