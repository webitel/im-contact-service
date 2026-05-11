package postgres

import (
	"context"
	"slices"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/webitel/webitel-go-kit/pkg/errors"

	"github.com/webitel/im-contact-service/infra/db/pg"
	"github.com/webitel/im-contact-service/internal/model"
)

type via struct {
	db *pg.PgxDB
}

func newViaStore(db *pg.PgxDB) *via {
	return &via{db: db}
}

func (communicationStore *via) Create(ctx context.Context, communication *model.ViaCommunication) (*model.ViaCommunication, error) {
	stmt, args := communicationStore.prepareCreateStmt(communication)

	rows, err := communicationStore.db.Master().Query(ctx, stmt, args)
	if err != nil {
		return nil, errors.Internal(
			"querying create contact communication stmt",
			errors.WithCause(err),
			errors.WithID("postgres.communication.create"),
			errors.WithValue("contact_id", communication.ContactID.String()),
			errors.WithValue("via", communication.Via),
		)
	}

	savedCommunication, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[model.ViaCommunication])
	if err != nil {
		if ok, rerr := pg.ErrorIntegrityViolation(err); ok {
			return nil, errors.Wrap(
				rerr,
				errors.WithID("postgres.communication.create"),
				errors.WithValue("contact_id", communication.ContactID.String()),
				errors.WithValue("via", communication.Via),
			)
		}

		return nil, errors.Internal(
			"collecting saved communication row",
			errors.WithCause(err),
			errors.WithID("postgres.communication.create"),
			errors.WithValue("contact_id", communication.ContactID.String()),
			errors.WithValue("via", communication.Via),
		)
	}

	return savedCommunication, nil
}

func (communicationStore *via) prepareCreateStmt(communication *model.ViaCommunication) (string, pgx.NamedArgs) {
	query := `
		insert into "im_contact"."via" (
			"contact_id", "via", "disable", "disable_reason", "metadata"
		)
		values (
			@ContactID, @Via,  @Disable, @DisableReason, @Metadata
		)
		returning "contact_id", "via", "disable","disable_reason", "metadata", "created_at", "updated_at"
	`

	args := pgx.NamedArgs{
		"ContactID":     communication.ContactID,
		"Via":           communication.Via,
		"Disable":       communication.Disable,
		"DisableReason": communication.DisableReason,
		"Metadata":      communication.Metadata,
	}

	return query, args
}

func (communicationStore *via) Update(ctx context.Context, communication *model.ViaCommunication) (*model.ViaCommunication, error) {
	stmt, args := communicationStore.prepareUpdateStmt(communication)

	rows, err := communicationStore.db.Master().Query(ctx, stmt, args)
	if err != nil {
		return nil, errors.Internal("executing update communication option", errors.WithCause(err), errors.WithID("postgres.communication.update"))
	}

	updated, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[model.ViaCommunication])
	if err != nil {
		if ok, rerr := pg.ErrorIntegrityViolation(err); ok {
			return nil, errors.Wrap(rerr, errors.WithID("postgres.communication.update"), errors.WithValue("contact_id", communication.ContactID.String()))
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.NotFound("zero records found for update", errors.WithCause(err), errors.WithID("postgres.communication.update"))
		}

		return nil, errors.Internal("collecting communication record row", errors.WithCause(err), errors.WithID("postgres.communication.update"), errors.WithValue("contact_id", communication.ContactID.String()))
	}

	return updated, nil
}

func (communicationStore *via) prepareUpdateStmt(communication *model.ViaCommunication) (string, pgx.NamedArgs) {
	stmt := `
		update "im_contact"."via"
		set
			"disable"=@Disable,
			"disable_reason"=@DisableReason,
			"metadata" = @Metadata
		where ("contact_id","via") = (@ContactID, @Via)
		returning
			"contact_id",
			"via",
			"disable",
			"disable_reason",
			"created_at",
			"updated_at",
			"metadata"
	`

	args := pgx.NamedArgs{
		"Via":           communication.Via,
		"Disable":       communication.Disable,
		"DisableReason": communication.DisableReason,
		"Metadata":      communication.Metadata,
		"ContactID":     communication.ContactID,
	}

	return stmt, args
}

func (communicationStore *via) PartialUpdate(ctx context.Context, updateCommand *model.CommunicationViaPartialUpdateCmd) (*model.ViaCommunication, error) {
	stmt, args, err := communicationStore.preparePartialUpdateStmt(updateCommand)
	if err != nil {
		return nil, err
	}

	rows, err := communicationStore.db.Master().Query(ctx, stmt, args...)
	if err != nil {
		return nil, errors.Internal("executing partial update query", errors.WithCause(err), errors.WithID("postgres.communication.partial_update"))
	}

	updated, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[model.ViaCommunication])
	if err != nil {
		if ok, rerr := pg.ErrorIntegrityViolation(err); ok {
			return nil, errors.Wrap(rerr, errors.WithID("postgres.communication.partial_update"), errors.WithValue("contact_id", updateCommand.ContactID.String()))
		}

		return nil, errors.Internal("collecting partial update query result records", errors.WithCause(err), errors.WithID("postgres.communication.partial_update"))
	}

	return updated, nil
}

func (communicationStore *via) preparePartialUpdateStmt(updateCommand *model.CommunicationViaPartialUpdateCmd) (string, []any, error) {
	communicationUpdateBuilder := sq.Update("im_contact.via").PlaceholderFormat(sq.Dollar)

	for _, field := range updateCommand.Fields {
		switch field {
		case "disable":
			communicationUpdateBuilder = communicationUpdateBuilder.Set("disable", updateCommand.Disable)
		case "disable_reason":
			communicationUpdateBuilder = communicationUpdateBuilder.Set("disable_reason", updateCommand.DisableReason)
		case "metadata":
			communicationUpdateBuilder = communicationUpdateBuilder.Set("metadata", updateCommand.Metadata)
		default:
			return "", nil, errors.InvalidArgument("unsupported field: "+field, errors.WithID("postgres.communication.prepare_partial_update_stmt"))
		}
	}

	communicationUpdateBuilder = communicationUpdateBuilder.Where(sq.Eq{"contact_id": updateCommand.ContactID})
	communicationUpdateBuilder = communicationUpdateBuilder.Where(sq.Eq{"via": updateCommand.Via})
	communicationUpdateBuilder = communicationUpdateBuilder.Suffix("returning contact_id, via, disable, disable_reason, created_at, updated_at, metadata")

	stmt, args, err := communicationUpdateBuilder.ToSql()
	if err != nil {
		return "", nil, errors.InvalidArgument("preparing partial update stmt", errors.WithCause(err), errors.WithID("postgres.communication.prepare_partial_update_stmt"))
	}

	return stmt, args, nil
}

func (communicationStore *via) Search(ctx context.Context, filter *model.SearchViaCommunicationsFilter) ([]*model.ViaCommunication, error) {
	stmt, args, err := communicationStore.prepareSearchStmt(filter)
	if err != nil {
		return nil, err
	}

	rows, err := communicationStore.db.Master().Query(ctx, stmt, args...)
	if err != nil {
		return nil, errors.Internal("executing search contact communication request", errors.WithCause(err), errors.WithID("postgres.communication.search"))
	}

	records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[model.ViaCommunication])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Internal("collecting search contact communication records", errors.WithCause(err), errors.WithID("postgres.communication.search"))
	}

	return records, nil
}

func (communicationStore *via) prepareSearchStmt(filter *model.SearchViaCommunicationsFilter) (string, []any, error) {
	const communicationaAllias = "c"

	communicationRecord := new(model.ViaCommunication)
	availableFields := communicationRecord.AvailableFields()

	validatedFields := slices.DeleteFunc(filter.Fields, func(f string) bool {
		return !slices.Contains(availableFields, f)
	})

	if len(validatedFields) == 0 {
		filter.Fields = communicationRecord.DefaultFields()
	}

	selectFields := make([]string, 0, len(filter.Fields))
	for _, field := range filter.Fields {
		selectFields = append(selectFields, Ident(communicationaAllias, field))
	}

	sb := sq.Select(selectFields...).From(communicationRecord.TableName() + " as " + communicationaAllias).PlaceholderFormat(sq.Dollar)
	sb = ApplyPaging(filter.Page, filter.Limit, sb)

	if sortingField, sortOperator := ExtractSortingOperator(filter.Sort); sortOperator != "" && sortingField != "" {
		sb = sb.OrderBy(sortingField + " " + sortOperator)
	}

	if len(filter.ContactIDs) > 0 {
		sb = sb.Where(
			sq.Eq{Ident(communicationaAllias, "contact_id"): filter.ContactIDs},
		)
	}

	if len(filter.Vias) > 0 {
		sb = sb.Where(sq.Eq{Ident(communicationaAllias, "via"): filter.Vias})
	}

	if disabled := filter.Disabled; disabled != nil {
		sb = sb.Where(sq.Eq{Ident(communicationaAllias, "disable"): *disabled})
	}

	stmt, args, err := sb.ToSql()
	if err != nil {
		return "", nil, errors.InvalidArgument(
			"preparing search communications stmt",
			errors.WithID("postgres.communication.prepare_search_stmt"),
			errors.WithCause(err),
		)
	}

	return stmt, args, nil
}
