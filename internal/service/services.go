package service

import (
	pubsubadapter "github.com/webitel/im-contact-service/internal/adapter/pubsub"
	"github.com/webitel/im-contact-service/internal/handler/amqp"
	"go.uber.org/fx"
)

var Module = fx.Module("service",
    fx.Provide(
        pubsubadapter.NewPublisherProvider,
        func(pp *pubsubadapter.PublisherProvider) (EventPublisher, error) {
            wmPub, err := pp.Build("im.contacts")
            if err != nil {
                return nil, err
            }
            return pubsubadapter.NewEventDispatcher(wmPub), nil
        },

        pubsubadapter.NewSubscriberProvider,
        amqp.NewMessageHandler,
        amqp.NewWatermillRouter,

        fx.Annotate(
            NewContactService,
            fx.As(new(Contacter)),
            fx.As(new(amqp.DomainDeletedEventHandler)),
        ),

        fx.Annotate(
            NewBaseBotManager,
            fx.As(new(BotManager)),
        ),
    ),

    fx.Invoke(amqp.RegisterHandlers),
)
