package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/store"
)

var (
	_ store.ContactStore = (*contactStore)(nil)
)

type contactStore struct {
}

// CanSend implements [store.ContactStore].
func (c *contactStore) CanSend(ctx context.Context, query *dto.CanSendQuery) (bool, error) {
	panic("unimplemented")
}

// Create implements [store.ContactStore].
func (c *contactStore) Create(ctx context.Context, contact *model.Contact) (*model.Contact, error) {
	panic("unimplemented")
}

// Delete implements [store.ContactStore].
func (c *contactStore) Delete(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// Search implements [store.ContactStore].
func (c *contactStore) Search(ctx context.Context, filter *dto.ContactSearchFilter) ([]*model.Contact, error) {
	panic("unimplemented")
}

// Update implements [store.ContactStore].
func (c *contactStore) Update(ctx context.Context, updater *dto.UpdateContactCommand) (*model.Contact, error) {
	panic("unimplemented")
}

func NewContactStore() *contactStore {
	var store = new(contactStore)
	{
		//setting
	}
	return store
}
