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
	Delete(ctx context.Context, command *dto.DeleteContactCommand) error
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
	if input == nil {
		return nil, errors.InvalidArgument("input is nil")
	}

	if input.Username == "" {
		return nil, errors.InvalidArgument("username is required")
	}

	if input.IssuerId == uuid.Nil {
		return nil, errors.InvalidArgument("issuerId is required")
	}

	if input.ApplicationId == uuid.Nil {
		return nil, errors.InvalidArgument("applicationId is required")
	}

	out, err := s.store.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	// We pass only the event object. The dispatcher will call event.Topic() internally.
	if err := s.publisher.Publish(ctx, domain.NewContactCreatedEvent(out)); err != nil {
		return out, err
	}

	return out, nil
}

// Update modifies an existing contact and publishes a ContactUpdatedEvent.
func (s *ContactService) Update(ctx context.Context, input *dto.UpdateContactCommand) (*model.Contact, error) {
	if input == nil {
		return nil, errors.InvalidArgument("input is nil")
	}

	if input.Id == uuid.Nil {
		return nil, errors.InvalidArgument("id is required")
	}

	out, err := s.store.Update(ctx, input)
	if err != nil {
		return nil, err
	}

	if err := s.publisher.Publish(ctx, domain.NewContactUpdatedEvent(out)); err != nil {
		return out, err
	}

	return out, nil
}

func (s *ContactService) Delete(ctx context.Context, command *dto.DeleteContactCommand) error {
	if command.Id == uuid.Nil {
		return errors.InvalidArgument("id is required")
	}

	err := s.store.Delete(ctx, command)
	if err != nil {
		return err
	}

	// Event handles its own topic and timestamp logic.
	return s.publisher.Publish(ctx, domain.NewContactDeletedEvent(command.Id))
}
