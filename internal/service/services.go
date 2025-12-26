package service

import (
	"go.uber.org/fx"

	"github.com/webitel/im-contact-service/internal/adapter/pubsub"
)

var Module = fx.Module("service",
	fx.Provide(
		pubsub.NewPublisherProvider,
		func(pp *pubsub.PublisherProvider) (EventPublisher, error) {
			// pp.Build return low level message.Publisher (Watermill)
			wmPub, err := pp.Build("im.contacts")
			if err != nil {
				return nil, err
			}
			return pubsub.NewEventDispatcher(wmPub), nil
		},
		fx.Annotate(
			NewContactService,
			fx.As(new(Contacter)),
		),
	),
)
