package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/store/queries"
)

type (
	Store        interface{}
	ContactStore interface {
		Search(ctx context.Context, filter *model.ContactSearchRequest) ([]*model.Contact, error)
		Create(ctx context.Context, contact *model.Contact) (*model.Contact, error)
		Update(ctx context.Context, updater *model.UpdateContactRequest) (*model.Contact, error)
		Delete(ctx context.Context, command *model.DeleteContactRequest) error
		ClearByDomain(ctx context.Context, domainId int) error
		Upsert(ctx context.Context, contact *model.Contact) (*model.Contact, bool, error)
		PartialUpdate(ctx context.Context, query queries.Query) (*model.Contact, error)
		DeleteBotByFlowID(ctx context.Context, flowID string) error
	}
	SettingsStore interface {
		Get(ctx context.Context, contactID uuid.UUID) (*model.ContactSettings, error)
		Update(ctx context.Context, command *model.UpdateContactSettingsRequest) (*model.ContactSettings, error)
		Create(ctx context.Context, command *model.CreateContactSettingsRequest) (*model.ContactSettings, error)	
	}

)
    
