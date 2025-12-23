package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/handler/grpc/mapper"

	"github.com/webitel/webitel-go-kit/pkg/errors"

	"github.com/webitel/im-contact-service/gen/go/api/v1"
	grpc_srv "github.com/webitel/im-contact-service/infra/server/grpc"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/service"
	"github.com/webitel/im-contact-service/internal/service/dto"
)

var _ impb.ContactsServer = &ContactService{}

type ContactService struct {
	impb.UnimplementedContactsServer

	logger  *slog.Logger
	handler service.Contacter
}

func NewContactService(handler service.Contacter, logger *slog.Logger) *ContactService {
	return &ContactService{handler: handler, logger: logger}
}

func RegisterContactService(server *grpc_srv.Server, service *ContactService) error {
	impb.RegisterContactsServer(server.Server, service)
	return nil
}

func (c *ContactService) SearchContact(ctx context.Context, request *impb.SearchContactRequest) (*impb.ContactList, error) {
	contacts, err := c.handler.Search(ctx, &dto.ContactSearchFilter{
		Page:    request.GetPage(),
		Size:    request.GetSize(),
		Q:       request.GetQ(),
		Sort:    request.GetSort(),
		Fields:  request.GetFields(),
		Apps:    request.GetAppId(),
		Issuers: request.GetIssId(),
		Types:   request.GetType(),
	})
	if err != nil {
		return nil, err
	}

	result := &impb.ContactList{
		Page:     request.GetPage(),
		Size:     request.GetSize(),
		Contacts: make([]*impb.Contact, len(contacts)),
	}
	if len(contacts) > int(request.GetSize()) {
		result.Next = true
		contacts = contacts[:request.GetSize()-1]
	}

	for _, contact := range contacts {
		marshaledContact, err := mapper.MarshalContact(contact)
		if err != nil {
			return nil, err
		}
		result.Contacts = append(result.Contacts, marshaledContact)
	}

	return result, nil
}

func (c *ContactService) CreateContact(ctx context.Context, request *impb.CreateContactRequest) (*impb.Contact, error) {
	timeNow := time.Now()

	contact, err := c.handler.Create(ctx, &model.Contact{
		BaseModel: model.BaseModel{
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		},
		IssuerId:      request.GetIssId(),
		ApplicationId: request.GetAppId(),
		Type:          request.GetType(),
		Name:          request.GetName(),
		Username:      request.GetUsername(),
		Metadata:      request.GetMetadata(),
	})
	if err != nil {
		return nil, err
	}

	return mapper.MarshalContact(contact)
}

func (c *ContactService) UpdateContact(ctx context.Context, request *impb.UpdateContactRequest) (*impb.Contact, error) {
	contactId, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, errors.New("invalid contact id", errors.WithCause(err))
	}

	updatedContact, err := c.handler.Update(ctx, &dto.UpdateContactCommand{
		Id:       contactId,
		Name:     request.GetName(),
		Username: request.GetUsername(),
		Metadata: request.GetMetadata(),
	})
	if err != nil {
		return nil, err
	}

	return mapper.MarshalContact(updatedContact)
}

func (c *ContactService) DeleteContact(ctx context.Context, request *impb.DeleteContactRequest) (*impb.Contact, error) {
	contactId, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, errors.New("invalid contact id", errors.WithCause(err))
	}
	err = c.handler.Delete(ctx, contactId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *ContactService) CanSend(ctx context.Context, request *impb.CanSendRequest) (*impb.CanSendResponse, error) {
	return &impb.CanSendResponse{
		Can: true,
	}, nil
}
