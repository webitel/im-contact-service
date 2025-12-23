package service

import (
	infrapubsub "github.com/webitel/im-contact-service/infra/pubsub"
	"github.com/webitel/im-contact-service/infra/pubsub/factory"
	"github.com/webitel/im-contact-service/internal/adapter/pubsub"
	"go.uber.org/fx"
)

var Module = fx.Module("service",
	fx.Provide(
		// Register the service
		NewContactService,

		func(p infrapubsub.Provider) (EventPublisher, error) {
			wmPub, err := p.GetFactory().BuildPublisher(&factory.PublisherConfig{
				Exchange: factory.ExchangeConfig{
					Name:    "im.contacts",
					Type:    "topic",
					Durable: true,
				},
			})
			if err != nil {
				return nil, err
			}

			return pubsub.NewWatermillPublisher(wmPub), nil
		},
	),
)
