package store

import (
	"context"

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
)
