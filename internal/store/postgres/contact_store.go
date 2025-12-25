package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/webitel/im-contact-service/infra/db/pg"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/store"
)

var _ store.ContactStore = (*contactStore)(nil)

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
				domain_id, issuer_id, subject_id,
				application_id, type, name, username, metadata
			)
			values(
				@domain_id, @issuer_id, @subject_id,
				@application_id, @type, @name, @username, @metadata
			)
			returning
				id, domain_id, created_at, updated_at, subject_id,
				issuer_id, application_id, type, name, username, metadata
		`
		args = pgx.NamedArgs{
			"domain_id":      contact.DomainId,
			"issuer_id":      contact.IssuerId,
			"subject_id":     contact.SubjectId,
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

	selectFields := "*"
	if len(filter.Fields) > 0 {
		selectFields = strings.Join(store.SanitizeFields(filter.Fields, model.ContactAllowedFields()), ",")
	}

	sortClause := store.ValidateAndFormatSort(filter.Sort, model.ContactAllowedFields())
	limit := max(filter.Size, 1)
	offset := (filter.Page - 1) * filter.Size

	var (
		query = fmt.Sprintf(`
        SELECT %s
        FROM im_contact.contact
        WHERE domain_id = @domain_id
            AND (@ids::uuid[] IS NULL OR id = ANY(@ids::uuid[]))
            AND (@Q::text IS NULL OR username ILIKE @Q OR name ILIKE @Q)
            AND (@apps::text[] IS NULL OR application_id = ANY(@apps::text[]))
            AND (@issuers::text[] IS NULL OR issuer_id = ANY(@issuers::text[]))
            AND (@types::text[] IS NULL OR type = ANY(@types::text[]))
        ORDER BY %s
        LIMIT @limit OFFSET @offset`, selectFields, sortClause)

		args = pgx.NamedArgs{
			"domain_id": filter.DomainId,
			"ids":       filter.Ids,
			"Q":         filter.Q,
			"apps":      filter.Apps,
			"issuers":   filter.Issuers,
			"types":     filter.Types,
			"limit":     limit + 1,
			"offset":    offset,
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
			returning id, domain_id, created_at, updated_at, subject_id,
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

func (c *contactStore) ClearByDomain(ctx context.Context, domainId int) error {
	var (
		query = `
			delete from im_contact.contact
			where domain_id = @domain_id
		`
		args = pgx.NamedArgs{
			"domain_id": domainId,
		}
	)

	if _, err := c.db.Master().Exec(ctx, query, args); err != nil {
		return fmt.Errorf("contactStore.ClearByDomain (id = %d): %w", domainId, err)
	}
	return nil
}
