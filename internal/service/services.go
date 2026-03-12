package service

import (
	"context"

	pubsubadapter "github.com/webitel/im-contact-service/internal/adapter/pubsub"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/handler/amqp"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"go.uber.org/fx"
)

type ContactSettingsService interface {
	Get(ctx context.Context, req *dto.GetContactSettingsRequest) (*model.ContactSettings, error)
	Update(ctx context.Context, request *dto.UpdateContactSettingsRequest) (*model.ContactSettings, error) 
	Create(ctx context.Context, request *dto.CreateContactSettingsRequest) (*model.ContactSettings, error)
}

type ContactService interface {
	Search(ctx context.Context, filter *dto.ContactSearchFilter) ([]*model.Contact, error)
	Create(ctx context.Context, input *model.Contact) (*model.Contact, error)
	Update(ctx context.Context, input *dto.UpdateContactCommand) (*model.Contact, error)
	Delete(ctx context.Context, input *dto.DeleteContactCommand) error
	CanSend(ctx context.Context, query *dto.CanSendQuery) error
	Upsert(ctx context.Context, contact *model.Contact) (*model.Contact, error)
	PartialUpdate(ctx context.Context, cmd *dto.PartialUpdateContactCommand) (*model.Contact, error)
    DeleteByDomain(ctx context.Context, domainId int) error 
    DeleteBotByFlowID(ctx context.Context, flowID string) error 
}

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
            fx.As(new(amqp.DomainEventsHandler)),
            fx.As(fx.Self()),
        ),
        NewContactSettingService,
    ),

    fx.Invoke(amqp.RegisterHandlers),
)
