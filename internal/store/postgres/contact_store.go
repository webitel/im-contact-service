package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"

	"github.com/webitel/webitel-go-kit/pkg/errors"

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
			"issuer_id":      contact.IssuerID,
			"subject_id":     contact.SubjectID,
			"application_id": contact.ApplicationID,
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
			if pgErr.Code == "23505" {
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

func (c *contactStore) Search(ctx context.Context, filter *model.ContactSearchRequest) ([]*model.Contact, error) {
	stmt, args, err := c.prepareContactSearchQuery(filter)
	if err != nil {
		return nil, err
	}

	rows, err := c.db.Master().Query(ctx, stmt, args...)
	if err != nil {
		return nil, errors.Internal("querying contact search query", errors.WithCause(err), errors.WithID("postgres.contact_store.search"))
	}

	contacts, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[model.Contact])
	if err != nil {
		return nil, errors.Internal(
			"collecting contact search query result",
			errors.WithCause(err),
			errors.WithID("postgres.contact_store.search"),
			errors.WithValue("details", pg.ExtractPgErrorMap(err)),
		)
	}

	return contacts, nil
}

func (c *contactStore) prepareContactSearchQuery(filter *model.ContactSearchRequest) (string, []any, error) {
	const (
		contactAlias string = "c"
		viaAlias     string = "v"
	)

	const (
		linkVia int = 1 << iota
	)

	links := 0
	viaJoin := func(selectBuilder sq.SelectBuilder) sq.SelectBuilder {
		if links&linkVia != 0 {
			return selectBuilder
		}

		links |= linkVia

		selectBuilder = selectBuilder.LeftJoin(
			`lateral (
				select jsonb_agg(to_jsonb(v.*)) as via
				from "im_contact"."via" v
				where v.contact_id = c.id
			) v on true`,
		)

		return selectBuilder
	}

	searchFields := filter.Fields
	if len(searchFields) > 0 {
		if searchFields = store.SanitizeFields(searchFields, model.ContactAllowedFields()); len(searchFields) == 0 {
			return "", nil, errors.InvalidArgument("zero requested fields are allowed", errors.WithID("postgres.contact_store.prepare_contact_search_query"))
		}
	} else {
		searchFields = model.ContactAllowedFields()
	}

	contactSelect := (sq.SelectBuilder{}).From((*model.Contact)(nil).TableName() + " " + contactAlias).PlaceholderFormat(sq.Dollar)

	for _, field := range searchFields {
		switch field {
		case "via":
			contactSelect = viaJoin(contactSelect)
			contactSelect = contactSelect.Columns(Ident(viaAlias, field))
		default:
			contactSelect = contactSelect.Columns(Ident(contactAlias, field))
		}
	}

	contactSelect = ApplyPaging(int(filter.Page), int(filter.Size), contactSelect)

	if sortingField, sortOperator := ExtractSortingOperator(filter.Sort); sortOperator != "" && sortingField != "" {
		contactSelect = contactSelect.OrderBy(sortingField + " " + sortOperator)
	}

	if filter.DomainID != nil && *filter.DomainID > 0 {
		contactSelect = contactSelect.Where(sq.Eq{Ident(contactAlias, "domain_id"): *filter.DomainID})
	}

	if len(filter.IDs) > 0 {
		contactSelect = contactSelect.Where(sq.Eq{Ident(contactAlias, "id"): filter.IDs})
	}

	if q := filter.Q; q != nil && *q != "" {
		searchPattern := *q + "%"
		contactSelect = contactSelect.Where(sq.Or{
			sq.ILike{Ident(contactAlias, "username"): searchPattern},
			sq.ILike{Ident(contactAlias, "name"): searchPattern},
		})
	}

	if len(filter.Apps) > 0 {
		contactSelect = contactSelect.Where(sq.Eq{Ident(contactAlias, "application_id"): filter.Apps})
	}

	if len(filter.Issuers) > 0 {
		contactSelect = contactSelect.Where(sq.Eq{Ident(contactAlias, "issuer_id"): filter.Issuers})
	}

	if len(filter.Types) > 0 {
		contactSelect = contactSelect.Where(sq.Eq{Ident(contactAlias, "type"): filter.Types})
	}

	if len(filter.Subjects) > 0 {
		contactSelect = contactSelect.Where(sq.Eq{Ident(contactAlias, "subject"): filter.Subjects})
	}

	if filter.OnlyBots != nil {
		contactSelect = contactSelect.Where(sq.Eq{Ident(contactAlias, "is_bot"): *filter.OnlyBots})
	}

	stmt, args, err := contactSelect.ToSql()
	if err != nil {
		return "", nil, errors.New("building stmt for contact search", errors.WithCause(err), errors.WithCode(codes.FailedPrecondition), errors.WithID("postgres.contact_store.prepare_contact_search_query"))
	}

	return stmt, args, nil
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

func (c *contactStore) ClearByDomain(ctx context.Context, domainID int) error {
	var (
		query = `
			delete from im_contact.contact
			where domain_id = @domain_id
		`
		args = pgx.NamedArgs{
			"domain_id": domainID,
		}
	)

	if _, err := c.db.Master().Exec(ctx, query, args); err != nil {
		return fmt.Errorf("contactStore.ClearByDomain (id = %d): %w", domainID, err)
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
	sql, args, err := query.ToSQL()
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
	stmt := `
		with ins as (
			insert into "im_contact"."contact" (
				"domain_id", "issuer_id", "subject_id", "application_id", "type", "name", "username", "metadata"
			)
			values (
				@DomainID, @Iss, @Sub, @App, @Type, @Name, @Username, @Metadata
			)
			on conflict ("domain_id", "issuer_id", "subject_id")
			do update set
				"updated_at" = now(),
				"name" = excluded.name,
			 	"username" = excluded.username,
				"metadata" = excluded.metadata
			where
				("im_contact"."contact"."name", "im_contact"."contact"."username", "im_contact"."contact"."metadata")
				is distinct from
				(excluded.name, excluded.username, excluded.metadata)
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
		)
		select *
		from ins
	 	union all
		select
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
			false as is_insert from im_contact.contact
		where domain_id = @DomainID
		and issuer_id = @Iss
		and subject_id = @Sub
		and not exists (
			select 1
			from ins
		);
	`
	args := pgx.NamedArgs{
		"DomainID": contact.DomainID,
		"Iss":      contact.IssuerID,
		"Sub":      contact.SubjectID,
		"App":      contact.ApplicationID,
		"Type":     contact.Type,
		"Name":     contact.Name,
		"Username": contact.Username,
		"Metadata": contact.Metadata,
	}

	var (
		result   model.Contact
		isInsert bool
	)

	if err := c.db.Master().QueryRow(ctx, stmt, args).Scan(
		&result.ID,
		&result.DomainID,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.IssuerID,
		&result.ApplicationID,
		&result.SubjectID,
		&result.Type,
		&result.Name,
		&result.Username,
		&result.Metadata,
		&isInsert,
	); err != nil {
		if ok, rerr := pg.ErrorIntegrityViolation(err); ok {
			return nil, false, errors.Wrap(rerr, errors.WithID("postgres.contact_store.upsert"))
		}

		return nil, false, errors.Internal("performing contact upsert", errors.WithCause(err), errors.WithID("postgres.contact_store.upsert"))
	}

	return &result, isInsert, nil
}
