package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/service/domain"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/store"
	"github.com/webitel/webitel-go-kit/pkg/errors"
)

// Contacter defines the primary API for managing contacts.
type Contacter interface {
	Search(ctx context.Context, filter *dto.ContactSearchFilter) ([]*model.Contact, error)
	Create(ctx context.Context, input *model.Contact) (*model.Contact, error)
	Update(ctx context.Context, input *dto.UpdateContactCommand) (*model.Contact, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CanSend(ctx context.Context, query *dto.CanSendQuery) (bool, error)
}

// EventPublisher defines the contract for publishing domain events.
type EventPublisher interface {
	Publish(ctx context.Context, topic string, event any) error
}

type ContactService struct {
	store     store.ContactStore
	publisher EventPublisher
}

// NewContactService creates a new ContactService instance.
func NewContactService(store store.ContactStore, publisher EventPublisher) *ContactService {
	return &ContactService{
		store:     store,
		publisher: publisher,
	}
}

// Search retrieves contacts based on the provided filter.
func (s *ContactService) Search(ctx context.Context, filter *dto.ContactSearchFilter) ([]*model.Contact, error) {
	if filter == nil {
		return nil, errors.InvalidArgument("filter is required")
	}
	return s.store.Search(ctx, filter)
}

// Create persists a new contact and publishes a ContactCreatedEvent.
func (s *ContactService) Create(ctx context.Context, input *model.Contact) (*model.Contact, error) {
	if err := s.validateCreate(input); err != nil {
		return nil, err
	}

	out, err := s.store.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	event := domain.NewContactCreatedEvent(out)
	if err := s.publisher.Publish(ctx, domain.ContactCreatedTopic, event); err != nil {
		return out, err
	}

	return out, nil
}

// Update modifies an existing contact and publishes a ContactUpdatedEvent.
func (s *ContactService) Update(ctx context.Context, input *dto.UpdateContactCommand) (*model.Contact, error) {
	if input == nil || input.Id == uuid.Nil {
		return nil, errors.InvalidArgument("input with a valid ID is required")
	}

	out, err := s.store.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	event := domain.NewContactUpdatedEvent(out)
	if err := s.publisher.Publish(ctx, domain.ContactUpdatedTopic, event); err != nil {
		return out, err
	}

	return out, nil
}

// Delete removes a contact and publishes a ContactDeletedEvent.
func (s *ContactService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.InvalidArgument("id is required")
	}

	if err := s.store.Delete(ctx, id); err != nil {
		return err
	}

	event := domain.NewContactDeletedEvent(id)
	return s.publisher.Publish(ctx, domain.ContactDeletedTopic, event)
}

// CanSend checks if a message can be sent to/from a contact.
func (s *ContactService) CanSend(ctx context.Context, query *dto.CanSendQuery) (bool, error) {
	if query == nil {
		return false, errors.InvalidArgument("query is required")
	}

	return s.store.CanSend(ctx, query)
}

// validateCreate performs business rules validation for new contacts.
func (s *ContactService) validateCreate(input *model.Contact) error {
	if input == nil {
		return errors.InvalidArgument("input is nil")
	}
	if input.Username == "" {
		return errors.InvalidArgument("username is required")
	}
	if input.IssuerId == uuid.Nil {
		return errors.InvalidArgument("issuerId is required")
	}

	switch input.Type {
	case model.Webitel, model.User, model.Bot:
		return nil
	default:
		return errors.InvalidArgument("invalid contact type")
	}
}
