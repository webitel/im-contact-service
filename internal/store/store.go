package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
)

type (
	Store interface{}

	ContactStore interface {
		Search(ctx context.Context, filter *dto.ContactSearchFilter) ([]*model.Contact, error)
		Create(ctx context.Context, contact *model.Contact) (*model.Contact, error)
		Update(ctx context.Context, updater *dto.UpdateContactCommand) (*model.Contact, error)
		Delete(ctx context.Context, id uuid.UUID) error
		CanSend(ctx context.Context, query *dto.CanSendQuery) (bool, error)
	}
)
