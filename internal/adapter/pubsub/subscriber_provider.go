package adapter

import (
	"github.com/ThreeDotsLabs/watermill/message"
	infrapubsub "github.com/webitel/im-contact-service/infra/pubsub"
	"github.com/webitel/im-contact-service/infra/pubsub/factory"
)

type SubscriberProvider struct {
	factory factory.Factory
}

func NewSubscriberProvider(p infrapubsub.Provider) *SubscriberProvider {
	return &SubscriberProvider{factory: p.GetFactory()}
}

func (sp *SubscriberProvider) Build(queue, exchange, routingKey string) (message.Subscriber, error) {
	return sp.factory.BuildSubscriber("im-contact-service", &factory.SubscriberConfig{
		Exchange: factory.ExchangeConfig{
			Name:    exchange,
			Type:    "topic",
			Durable: true,
		},
		Queue:      queue,
		RoutingKey: routingKey,
	})
}
