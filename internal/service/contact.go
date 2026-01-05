package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/domain/events"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/store"
	"github.com/webitel/webitel-go-kit/pkg/errors"
)

// Contacter defines the primary API for managing contacts.
type Contacter interface {
	Search(ctx context.Context, filter *dto.ContactSearchFilter) ([]*model.Contact, error)
	Create(ctx context.Context, input *model.Contact) (*model.Contact, error)
	Update(ctx context.Context, input *dto.UpdateContactCommand) (*model.Contact, error)
	Delete(ctx context.Context, input *dto.DeleteContactCommand) error
	CanSend(ctx context.Context, query *dto.CanSendQuery) error
}

// EventPublisher defines the contract for publishing domain events.
// Note: We removed the 'topic string' argument because the event knows its own topic.
type EventPublisher interface {
	Publish(ctx context.Context, event events.Event) error
}

type ContactService struct {
	store     store.ContactStore
	publisher EventPublisher
}

// NewContactService creates a new ContactService instance.
func NewContactService(store store.ContactStore, publisher EventPublisher) Contacter {
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

	event := events.NewContactCreated(out)
	if err := s.publisher.Publish(ctx, event); err != nil {
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

	event := events.NewContactUpdated(out)
	if err := s.publisher.Publish(ctx, event); err != nil {
		return out, err
	}

	return out, nil
}

// Delete removes a contact and publishes a ContactDeletedEvent.
func (s *ContactService) Delete(ctx context.Context, input *dto.DeleteContactCommand) error {
	if input.Id == uuid.Nil {
		return errors.InvalidArgument("id is required")
	}
	if input.DomainId == 0 {
		return errors.InvalidArgument("domainId is required")
	}

	if err := s.store.Delete(ctx, input); err != nil {
		return err
	}

	// Event handles its own topic and timestamp logic.
	return s.publisher.Publish(ctx, events.NewContactDeleted(input.Id))
}

// CanSend checks if a message can be sent to/from a contact.
func (s *ContactService) CanSend(ctx context.Context, query *dto.CanSendQuery) error {
	if query == nil {
		return errors.InvalidArgument("query is required")
	}

	peers, err := s.store.Search(ctx, &dto.ContactSearchFilter{Ids: []uuid.UUID{query.From, query.To}})
	if err != nil {
		return err
	}

	switch len(peers) {
	case 0:
		return errors.NotFound("no contacts found for the provided IDs")
	case 1:
		if query.From != query.To {
			return errors.NotFound("no contacts found for the provided IDs")
		}
	case 2:
		return nil
	default:
		return errors.InvalidArgument("too many contacts found for the provided IDs")
	}

	return nil
}

// validateCreate performs business rules validation for new contacts.
func (s *ContactService) validateCreate(input *model.Contact) error {
	if input == nil {
		return errors.InvalidArgument("input is nil")
	}

	if input.Username == "" {
		return errors.InvalidArgument("username is required")
	}

	if input.IssuerId == "" {
		return errors.InvalidArgument("issuerId is required")
	}

	if input.Type == "" {
		return errors.InvalidArgument("contact type is required")
	}
	return nil
}
