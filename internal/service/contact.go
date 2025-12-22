package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/store"
	"github.com/webitel/webitel-go-kit/pkg/errors"
)

type Contacter interface {
	Search(ctx context.Context, filter *dto.ContactSearchFilter) ([]*model.Contact, error)
	Create(ctx context.Context, input *model.Contact) (*model.Contact, error)
	Update(ctx context.Context, input *dto.UpdateContactCommand) (*model.Contact, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ContactService struct {
	store store.ContactStore
}

func NewContactService(store store.ContactStore) *ContactService {
	return &ContactService{
		store: store,
	}
}

func (s *ContactService) Search(ctx context.Context, filter *dto.ContactSearchFilter) ([]*model.Contact, error) {
	out, err := s.store.Search(ctx, filter)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *ContactService) Create(ctx context.Context, input *model.Contact) (*model.Contact, error) {
	if input == nil {
		return nil, errors.NewBadRequestError("contact.create.invalid", "input is nil")
	}

	if input.Username == "" {
		return nil, errors.NewBadRequestError("contact.create.invalid", "username is required")
	}

	if input.IssuerId == uuid.Nil {
		return nil, errors.NewBadRequestError("contact.create.invalid", "issuerId is required")
	}

	if input.ApplicationId == uuid.Nil {
		return nil, errors.NewBadRequestError("contact.create.invalid", "applicationId is required")
	}

	if !isValidContactType(input.Type) {
		return nil, errors.NewBadRequestError("contact.create.invalid", "type is invalid or missing")
	}

	out, err := s.store.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *ContactService) Update(ctx context.Context, input *dto.UpdateContactCommand) (*model.Contact, error) {
	if input == nil {
		return nil, errors.NewBadRequestError("contact.update.invalid", "input is nil")
	}

	if input.Id == uuid.Nil {
		return nil, errors.NewBadRequestError("contact.update.invalid", "id is required")
	}

	if input.Username == "" {
		return nil, errors.NewBadRequestError("contact.update.invalid", "username is cannot be empty")
	}

	out, err := s.store.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *ContactService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.NewBadRequestError("contact.delete.invalid", "id is required")
	}

	err := s.store.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func isValidContactType(t model.ContactType) bool {
	switch t {
	case model.Webitel, model.User, model.Bot:
		return true
	default:
		return false
	}
}
