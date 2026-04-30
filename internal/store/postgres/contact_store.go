package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/webitel/webitel-go-kit/pkg/errors"
	"google.golang.org/grpc/codes"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/webitel/im-contact-service/infra/db/pg"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/store"
	"github.com/webitel/im-contact-service/internal/store/queries"
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
				application_id, type, name, username, metadata, is_bot
			)
			values(
				@domain_id, @issuer_id, @subject_id,
				@application_id, @type, @name, @username, @metadata, @is_bot
			)
			returning
				id, domain_id, created_at, updated_at, subject_id,
				issuer_id, application_id, type, name, username, metadata, is_bot
		`
		args = pgx.NamedArgs{
			"domain_id":      contact.DomainID,
			"issuer_id":      contact.IssuerId,
			"subject_id":     contact.SubjectId,
			"application_id": contact.ApplicationId,
			"type":           contact.Type,
			"name":           contact.Name,
			"username":       contact.Username,
			"metadata":       contact.Metadata,
			"is_bot":         contact.IsBot,
		}
		result *model.Contact
	)

	row, err := c.db.Master().Query(ctx, query, args)
	if err != nil {
		return nil, errors.Internal("executing create contact query", errors.WithID("postgres.contact_store.create"), errors.WithCause(err))
	}

	if result, err = pgx.CollectExactlyOneRow(row, pgx.RowToAddrOfStructByNameLax[model.Contact]); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, errors.New("conflict: contact already exists", errors.WithCause(err), errors.WithCode(codes.AlreadyExists), errors.WithID("postgres.contact_store.create"))
			}
		}
		return nil, errors.Internal("collecting create contact query result", errors.WithCause(err), errors.WithID("postgres.contact_store.create"))
	}

	return result, nil
}

// Delete implements [store.ContactStore].
func (c *contactStore) Delete(ctx context.Context, command *model.DeleteContactRequest) error {
	var (
		query = `
			delete from im_contact.contact
			where domain_id = @domain_id and id = @id
		`
		args = pgx.NamedArgs{
			"id":        command.ID,
			"domain_id": command.DomainID,
		}
	)

	if _, err := c.db.Master().Exec(ctx, query, args); err != nil {
		return fmt.Errorf("failed to delete contact: %v", err)
	}
	return nil
}

// Search implements [store.ContactStore].
func (c *contactStore) Search(ctx context.Context, filter *model.ContactSearchRequest) ([]*model.Contact, error) {
	if filter.Q != nil && *filter.Q != "" {
		*filter.Q += "%"
	} else {
		filter.Q = nil
	}

	selectFields := ""
	if len(filter.Fields) > 0 {
		selectFields = strings.Join(store.SanitizeFields(filter.Fields, model.ContactAllowedFields()), ",")
	} else {
		selectFields = strings.Join(model.ContactAllowedFields(), ",")
	}

	sortClause := store.ValidateAndFormatSort(filter.Sort, model.ContactAllowedFields())
	limit := max(filter.Size, 1)
	offset := max((filter.Page-1)*filter.Size, 0)

	var (
		query = fmt.Sprintf(`
        SELECT %s
        FROM im_contact.contact
        WHERE (@domain_id::int IS NULL OR domain_id = @domain_id)
            AND (@ids::uuid[] IS NULL OR id = ANY(@ids::uuid[]))
            AND (@Q::text IS NULL OR username ILIKE @Q OR name ILIKE @Q)
            AND (@apps::text[] IS NULL OR application_id = ANY(@apps::text[]))
            AND (@issuers::text[] IS NULL OR issuer_id = ANY(@issuers::text[]))
            AND (@types::text[] IS NULL OR type = ANY(@types::text[]))
			AND(@subjects::text[] IS NULL OR subject_id = ANY(@subjects::text[]))
			AND(@is_bot::BOOL IS NULL OR is_bot = @is_bot)
        ORDER BY %s
        LIMIT @limit OFFSET @offset`, selectFields, sortClause)

		args = pgx.NamedArgs{
			"domain_id": filter.DomainID,
			"ids":       arrayOrNull(filter.IDs),
			"Q":         filter.Q,
			"apps":      arrayOrNull(filter.Apps),
			"issuers":   arrayOrNull(filter.Issuers),
			"types":     arrayOrNull(filter.Types),
			"limit":     limit + 1,
			"offset":    offset,
			"subjects":  arrayOrNull(filter.Subjects),
			"is_bot":    filter.OnlyBots,
		}
		contacts []*model.Contact
	)

	rows, err := c.db.Master().Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("error search contacts: %v", err)
	}

	if contacts, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[model.Contact]); err != nil {
		return nil, fmt.Errorf("error search contacts: %v", err)
	}

	return contacts, nil
}

func arrayOrNull[T any](v []T) any {
	if len(v) > 0 {
		return v
	}
	return nil
}

// Update implements [store.ContactStore].
func (c *contactStore) Update(ctx context.Context, updater *model.UpdateContactRequest) (*model.Contact, error) {
	var (
		query = `
			update im_contact.contact
			set
				name = coalesce(@name, name),
				username = coalesce(@username, username),
				metadata = coalesce(@metadata, metadata),
				subject_id = coalesce(@subject, subject_id),
				updated_at = now()
			where domain_id = @domain_id
				and id = @id
			returning id, domain_id, created_at, updated_at, subject_id,
				issuer_id, application_id, type, name, username, metadata
		`
		args = pgx.NamedArgs{
			"id":        updater.ID,
			"domain_id": updater.DomainID,
			"name":      updater.Name,
			"username":  updater.Username,
			"metadata":  updater.Metadata,
			"subject":   updater.Subject,
		}
		result *model.Contact
	)

	rows, err := c.db.Master().Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("error creating contact: %v", err)
	}

	if result, err = pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[model.Contact]); err != nil {
		return nil, fmt.Errorf("error creating contact: %v", err)
	}

	return result, nil
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

func (c *contactStore) DeleteBotByFlowID(ctx context.Context, flowID string) error {
	var (
		query = `
			delete from im_contact.contact
			where is_bot = TRUE AND subject_id = @flow_id
		`
		args = pgx.NamedArgs{
			"flow_id": flowID,
		}
	)

	if _, err := c.db.Master().Exec(ctx, query, args); err != nil {
		return errors.Internal("error occurred while executing query", errors.WithCause(err))
	}
	return nil
}

func (c *contactStore) PartialUpdate(ctx context.Context, query queries.Query) (*model.Contact, error) {
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.db.Master().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	contact, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[model.Contact])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrUpdatedContactNotFound
		}

		return nil, err
	}

	return contact, nil
}

func (c *contactStore) Upsert(ctx context.Context, contact *model.Contact) (*model.Contact, bool, error) {
	var (
		query = `
			insert into im_contact.contact(
				domain_id,
				issuer_id,
				subject_id,
				application_id,
				type,
				name,
				username,
				metadata
			) values (
				@DomainId, @IssuerId, @SubjectId, @ApplicationId, @Type,
				@Name, @Username, @Metadata
			)
			on conflict (domain_id, issuer_id, subject_id)
			do update set
				updated_at = now(),
				name = excluded.name,
				username = excluded.username,
				metadata = excluded.metadata
			returning
				id,
				domain_id,
				created_at,
				updated_at,
				issuer_id,
				application_id,
				subject_id,
				type,
				name,
				username,
				metadata,
				(xmax = 0) as is_insert
		`
		args = pgx.NamedArgs{
			"DomainId":      contact.DomainID,
			"IssuerId":      contact.IssuerId,
			"SubjectId":     contact.SubjectId,
			"ApplicationId": contact.ApplicationId,
			"Type":          contact.Type,
			"Name":          contact.Name,
			"Username":      contact.Username,
			"Metadata":      contact.Metadata,
		}
		result   model.Contact
		isInsert bool
	)

	if err := c.db.Master().QueryRow(ctx, query, args).Scan(
		&result.ID,
		&result.DomainID,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.IssuerId,
		&result.ApplicationId,
		&result.SubjectId,
		&result.Type,
		&result.Name,
		&result.Username,
		&result.Metadata,
		&isInsert,
	); err != nil {
		return nil, false, err
	}

	return &result, isInsert, nil
}
