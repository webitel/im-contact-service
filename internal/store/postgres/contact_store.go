 package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/webitel/im-contact-service/infra/db/pg"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/store"
)

var (
	_ store.ContactStore = (*contactStore)(nil)
)

type contactStore struct {
	db *pg.PgxDB
}

func NewContactStore(db *pg.PgxDB) *contactStore {
	return &contactStore{
		db: db,
	}
}

// Create implements [store.ContactStore].
func (c *contactStore) Create(ctx context.Context, contact *model.Contact) (*model.Contact, error) {
	var (
		query = `
			insert into im_contact.contact(
				domain_id, created_by, updated_by, issuer_id, 
				application_id, type, name, username, metadata
			)
			values(
				@domain_id, @created_by, @updated_by, @issuer_id,
				@application_id, @type, @name, @username, @metadata
			)
			returning
				id, domain_id, created_at, updated_at, created_by, updated_by,
				issuer_id, application_id, type, name, username, metadata
		`
		args = pgx.NamedArgs{
			"domain_id":      contact.DomainId,
			"created_by":     contact.CreatedBy,
			"updated_by":     contact.UpdatedBy,
			"issuer_id":      contact.IssuerId,
			"application_id": contact.ApplicationId,
			"type":           contact.Type,
			"name":           contact.Name,
			"username":       contact.Username,
			"metadata":       contact.Metadata,
		}
		result model.Contact
	)

	if err := pgxscan.Get(ctx, c.db.Master(), &result, query, args); err != nil {
		return nil, fmt.Errorf("failed to create contact: %v", err)
	}

	return &result, nil
}

// Delete implements [store.ContactStore].
func (c *contactStore) Delete(ctx context.Context, command *dto.DeleteContactCommand) error {
	var (
		query = `
			delete from im_contact.contact
			where domain_id = @domain_id and id = @id
		`
		args = pgx.NamedArgs{
			"id":        command.Id,
			"domain_id": command.DomainId,
		}
	)

	if _, err := c.db.Master().Exec(ctx, query, args); err != nil {
		return fmt.Errorf("failed to delete contact: %v", err)
	}
	return nil
}

// Search implements [store.ContactStore].
func (c *contactStore) Search(ctx context.Context, filter *dto.ContactSearchFilter) ([]*model.Contact, error) {
	if filter.Q != nil && *filter.Q != "" {
		*filter.Q += "%"
	}

	var (
		fields = strings.Join(filter.Fields, ",")
		size   = func(lim int32) int32 {
			if filter.Size < 1 {
				return 1
			}
			return filter.Size + 1
		}(filter.Size)
		sort = func(order string) string {
			if len(order) > 2 {
				if order[0] == '+' {
					return order[0:] + " asc"
				} else if order[0] == '-' {
					return order[0:] + " desc"
				}
			}

			return "created_at desc"
		}(filter.Sort)
		query = `
			select
		` + fields + `
			from im_contact.contact
			where domain_id = @domain_id
				and (@ids::uuid[] is null or id = any(@ids::uuid[]))
				and (@Q::text is null or username ilike @Q or name ilike @Q)
				and (@apps::text[] is null or application_id = any(@apps::text[]))
				and (@issuers::text[] is null or issuer_id = any(@issuers::text[]))
				and (@types::text[] is null or type = any(@types::text[]))
			order by ` + sort + ` limit ` + strconv.Itoa(int(size)) + ` offset ` + strconv.Itoa(int((filter.Page-1)*filter.Size))

		args = pgx.NamedArgs{
			"domain_id": filter.DomainId,
			"ids":       filter.Ids,
			"Q":         filter.Q,
			"apps":      filter.Apps,
			"issuers":   filter.Issuers,
			"types":     filter.Types,
		}
		contacts []*model.Contact
	)

	if err := pgxscan.Select(ctx, c.db.Master(), &contacts, query, args); err != nil {
		return nil, fmt.Errorf("error search contacts: %v", err)
	}
	return contacts, nil
}

// Update implements [store.ContactStore].
func (c *contactStore) Update(ctx context.Context, updater *dto.UpdateContactCommand) (*model.Contact, error) {
	var (
		query = `
			update im_contact.contact
			set
				name = coalesce(@name, name),
				username = coalesce(@username, username),
				metadata = coalesce(@metadata, metadata),
				updated_at = now()
			where domain_id = @domain_id
				and id = @id
			returning id, domain_id, created_at, updated_at, created_by, updated_by,
				issuer_id, application_id, type, name, username, metadata
		`
		args = pgx.NamedArgs{
			"id":        updater.Id,
			"domain_id": updater.DomainId,
			"name":      updater.Name,
			"username":  updater.Username,
			"metadata":  updater.Metadata,
		}
		result model.Contact
	)

	if err := pgxscan.Get(ctx, c.db.Master(), &result, query, args); err != nil {
		return nil, fmt.Errorf("error creating contact: %v", err)
	}
	return &result, nil
}
