package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/webitel/webitel-go-kit/pkg/errors"

	impb "github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/handler/grpc/mapper"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/service"
	"github.com/webitel/im-contact-service/internal/utils"
)

var _ impb.ContactsServer = &ContactServer{}

type ContactServer struct {
	impb.UnimplementedContactsServer

	logger  *slog.Logger
	handler service.ContactService
	inMapper mapper.ContactInConverter
}





func NewContactService(handler service.ContactService, logger *slog.Logger) *ContactServer{
	return &ContactServer{handler: handler, logger: logger}
}

func (c *ContactServer) SearchContact(ctx context.Context, request *impb.SearchContactRequest) (*impb.ContactList, error) {
	ids := utils.Map(request.GetIds(), func(id string) uuid.UUID {
		parsed, err := uuid.Parse(id)
		if err != nil {
			return uuid.Nil
		}

		return parsed
	})
	page, size := ParsePagination(request.GetPage(), request.GetSize())

	contacts, err := c.handler.Search(ctx, &model.ContactSearchRequest{
		Page:     page,
		Size:     size, // + 1,
		Q:        &request.Q,
		Sort:     request.GetSort(),
		Fields:   request.GetFields(),
		Apps:     request.GetAppId(),
		Issuers:  request.GetIssId(),
		Types:    request.GetType(),
		Subjects: request.GetSubjects(),
		DomainID: int(request.GetDomainId()),
		IDs:      ids,
		OnlyBots: request.OnlyBots,
	})

	if err != nil {
		return nil, err
	}

	result := &impb.ContactList{
		Page:     page,
		Size:     size,
		Contacts: make([]*impb.Contact, 0, len(contacts)),
	}

	result.Contacts = utils.Map(contacts, func(contact *model.Contact) *impb.Contact {
		marshaledContact, _ := mapper.MarshalContact(contact)
		return marshaledContact
	})
	result.Contacts, result.Next = ResolvePaging(int(size), result.Contacts)

	return result, nil
}

func (c *ContactServer) CreateContact(ctx context.Context, request *impb.CreateContactRequest) (*impb.Contact, error) {
	timeNow := time.Now()

	contact, err := c.handler.Create(ctx, &model.Contact{
		BaseModel: model.BaseModel{
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			DomainID:  int(request.GetDomainId()),
		},
		IssuerId:      request.GetIssId(),
		ApplicationId: request.GetAppId(),
		Type:          request.GetType(),
		Name:          request.GetName(),
		Username:      request.GetUsername(),
		Metadata:      request.GetMetadata(),
		SubjectId:     request.GetSubject(),
		IsBot: request.GetIsBot(),
	})
	if err != nil {
		return nil, err
	}

	return mapper.MarshalContact(contact)
}

func (c *ContactServer) UpdateContact(ctx context.Context, request *impb.UpdateContactRequest) (*impb.Contact, error) {
	contactId, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, errors.New("invalid contact id", errors.WithCause(err))
	}

	updatedContact, err := c.handler.Update(ctx, &model.UpdateContactRequest{
		ID:       contactId,
		Name:     &request.Name,
		Username: &request.Username,
		Metadata: request.GetMetadata(),
		Subject:  request.GetSubject(),
		DomainID: int(request.GetDomainId()),
	})
	if err != nil {
		return nil, err
	}

	return mapper.MarshalContact(updatedContact)
}

func (c *ContactServer) DeleteContact(ctx context.Context, request *impb.DeleteContactRequest) (*impb.Contact, error) {
	contactId, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, errors.New("invalid contact id", errors.WithCause(err))
	}

	err = c.handler.Delete(ctx, &model.DeleteContactRequest{
		ID:       contactId,
		DomainID: int(request.GetDomainId()),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *ContactServer) Upsert(ctx context.Context, req *impb.CreateContactRequest) (*impb.Contact, error) {
	var (
		contact = &model.Contact{
			BaseModel: model.BaseModel{
				DomainID: int(req.GetDomainId()),
			},
			IssuerId:      req.GetIssId(),
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

func (c *ContactServer) Patch(ctx context.Context, request *impb.PatchContactRequest) (*impb.Contact, error) {
	contactPartialUpdateCmd := mapper.MapPatchContactRequestToPartialUpdateContactCommand(request)
	contact, err := c.handler.PartialUpdate(ctx, contactPartialUpdateCmd)
	if err != nil {
		return nil, err
	}

	return mapper.MarshalContact(contact)
}
