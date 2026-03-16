package service

import (
	"context"

	pubsubadapter "github.com/webitel/im-contact-service/internal/adapter/pubsub"
	"github.com/webitel/im-contact-service/internal/handler/amqp"
	"github.com/webitel/im-contact-service/internal/model"
	"go.uber.org/fx"
)

type ContactSettingsService interface {
	Get(ctx context.Context, req *model.GetContactSettingsRequest) (*model.ContactSettings, error)
	Update(ctx context.Context, request *model.UpdateContactSettingsRequest) (*model.ContactSettings, error) 
	Create(ctx context.Context, request *model.CreateContactSettingsRequest) (*model.ContactSettings, error)
}

type ContactService interface {
	Search(ctx context.Context, filter *model.ContactSearchRequest) ([]*model.Contact, error)
	Create(ctx context.Context, input *model.Contact) (*model.Contact, error)
	Update(ctx context.Context, input *model.UpdateContactRequest) (*model.Contact, error)
	Delete(ctx context.Context, input *model.DeleteContactRequest) error
	CanSend(ctx context.Context, query *model.CanSendRequest) error
	Upsert(ctx context.Context, contact *model.Contact) (*model.Contact, error)
	PartialUpdate(ctx context.Context, cmd *model.PartialUpdateContactRequest) (*model.Contact, error)
    DeleteByDomain(ctx context.Context, domainId int) error 
    DeleteBotByFlowID(ctx context.Context, flowID string) error 
}

type ContactPrivacyService interface {
    CanSend(ctx context.Context, query *model.CanSendRequest) error 
    CanInvite(ctx context.Context, query *model.CanInviteRequest) error

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
        NewContactPrivacyService,
    ),

    fx.Invoke(amqp.RegisterHandlers),
)
