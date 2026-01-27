package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/webitel/webitel-go-kit/pkg/errors"

	impb "github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/handler/grpc/mapper"
	"github.com/webitel/im-contact-service/internal/service"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/utils"
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

func (c *ContactService) SearchContact(ctx context.Context, request *impb.SearchContactRequest) (*impb.ContactList, error) {
	ids := utils.Map(request.GetIds(), func(id string) uuid.UUID {
		parsed, err := uuid.Parse(id)
		if err != nil {
			return uuid.Nil
		}

		return parsed
	})
	
	contacts, err := c.handler.Search(ctx, &dto.ContactSearchFilter{
		Page:    request.GetPage(),
		Size:    request.GetSize(),
		Q:       &request.Q,
		Sort:    request.GetSort(),
		Fields:  request.GetFields(),
		Apps:    request.GetAppId(),
		Issuers: request.GetIssId(),
		Types:   request.GetType(),
		Subjects: request.GetSubjects(),
		DomainId: int(request.GetDomainId()),
		Ids: ids,
	})

	if err != nil {
		return nil, err
	}

	result := &impb.ContactList{
		Page:     request.GetPage(),
		Size:     request.GetSize(),
		Contacts: make([]*impb.Contact, 0, len(contacts)),
	}

	if len(contacts) > int(request.GetSize()) {
		result.Next = true
		contacts = contacts[:request.GetSize()-1]
	}

	result.Contacts = utils.Map(contacts, func(contact *model.Contact) *impb.Contact {
		marshaledContact, _ := mapper.MarshalContact(contact)
		return marshaledContact
	})

	return result, nil
}

func (c *ContactService) CreateContact(ctx context.Context, request *impb.CreateContactRequest) (*impb.Contact, error) {
	timeNow := time.Now()

	contact, err := c.handler.Create(ctx, &model.Contact{
		BaseModel: model.BaseModel{
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			DomainId: int(request.GetDomainId()),
		},
		IssuerId:      request.GetIssId(),
		ApplicationId: request.GetAppId(),
		Type:          request.GetType(),
		Name:          request.GetName(),
		Username:      request.GetUsername(),
		Metadata:      request.GetMetadata(),
		SubjectId: request.GetSubject(),
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
		Name:     &request.Name,
		Username: &request.Username,
		Metadata: request.GetMetadata(),
		Subject: request.GetSubject(),
		DomainId: int(request.GetDomainId()),
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

	err = c.handler.Delete(ctx, &dto.DeleteContactCommand{
		Id: contactId,
		DomainId: int(request.GetDomainId()),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *ContactService) CanSend(ctx context.Context, request *impb.CanSendRequest) (*impb.CanSendResponse, error) {
	canSendQuery := mapper.CanSendRequest2Model(request)

	err := c.handler.CanSend(ctx, canSendQuery)
	if err != nil {
		return nil, err
	}

	return &impb.CanSendResponse{Can: true}, nil
}

func (c *ContactService) Upsert(ctx context.Context, req *impb.CreateContactRequest) (*impb.Contact, error) {
	var (
		contact = &model.Contact{
			BaseModel: model.BaseModel{
				DomainId: int(req.GetDomainId()),
			},
			IssuerId:  req.GetIssId(),
		ApplicationId: req.GetAppId(),
		Type:          req.GetType(),
		Name:          req.GetName(),
		Username:      req.GetUsername(),
		Metadata:      req.GetMetadata(),
		SubjectId:     req.GetSubject(),
		}
		err error
	)	

	if contact, err = c.handler.Upsert(ctx, contact); err != nil {
		return nil, err
	}
	
	return mapper.MarshalContact(contact)
}