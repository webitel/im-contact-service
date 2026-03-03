package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/store/queries"
)

type (
	Store        interface{}
	ContactStore interface {
		Search(ctx context.Context, filter *dto.ContactSearchFilter) ([]*model.Contact, error)
		Create(ctx context.Context, contact *model.Contact) (*model.Contact, error)
		Update(ctx context.Context, updater *dto.UpdateContactCommand) (*model.Contact, error)
		Delete(ctx context.Context, command *dto.DeleteContactCommand) error
		ClearByDomain(ctx context.Context, domainId int) error
		Upsert(ctx context.Context, contact *model.Contact) (*model.Contact, bool, error)
		PartialUpdate(ctx context.Context, query queries.Query) (*model.Contact, error)
	}
	SettingsStore interface {
		Get(ctx context.Context, contactID uuid.UUID) (*model.ContactSettings, error)
		Update(ctx context.Context, command *dto.UpdateContactSettingsRequest) (*model.ContactSettings, error)
		Create(ctx context.Context, command *dto.CreateContactSettingsRequest) (*model.ContactSettings, error)	
	}

)
