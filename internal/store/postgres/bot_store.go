package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/webitel/im-contact-service/infra/db/pg"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/store"
)

type botStore struct {
	db *pg.PgxDB
}

func NewBotStore(db *pg.PgxDB) *botStore {
	return &botStore{
		db: db,
	}
}

func (b *botStore) Create(ctx context.Context, bot *model.WebitelBot) (*model.WebitelBot, error) {
	var (
		query = `
			insert into im_contact.bot (
				domain_id, created_at, updated_at, flow_id, display_name
			)
			values (
				@DomainId, now(), now(), @FlowId, @DisplayName
			)
			returning
				id,
				created_at,
				updated_at
		`
		args = pgx.NamedArgs{
			"DomainId": bot.DomainId,
			"FlowId": bot.FlowId,
			"DisplayName": bot.DisplayName,
		}
	)

	if err := b.db.Master().QueryRow(ctx, query, args).Scan(&bot.Id, &bot.CreatedAt, &bot.UpdatedAt); err != nil {
		return nil, err
	}
	
	return bot, nil
}

func (b *botStore) Search(ctx context.Context, filter *dto.SearchBotRequest) ([]*model.WebitelBot, error) {
	limit := max(filter.Size, 1)
	offset := (filter.Page - 1) * filter.Size
	sortClause := store.ValidateAndFormatSort(filter.Sort, model.BotAllowedFields())

	var (
		query = fmt.Sprintf(`
			select
				id,
				domain_id,
				created_at,
				updated_at,
				flow_id,
				display_name
			from im_contact.bot
			where domain_id = @DomainId
				and (@Q::varchar is null or display_name ilike @Q)
				and (@FlowIds::int[] is null or flow_id = any(@FlowIds::int[]))
				and (@DisplayNames::text[] is null or display_name ilike any(@DisplayNames::text[]))
				and (@Ids::uuid[] is null or id = any(@Ids::uuid[]))
			order by %s 
			limit @Limit offset @Offset
		`, sortClause)

		args = pgx.NamedArgs{
			"DomainId": filter.DomainId,
			"Q": filter.GetQ(),
			"FlowIds": filter.FlowIds,
			"DisplayNames": filter.DisplayNames,
			"Ids": filter.Ids,
			"Limit": limit + 1,
			"Offset": offset,
		}

		bots = make([]*model.WebitelBot, 0)
	)

	rows, err := b.db.Master().Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		bot := new(model.WebitelBot)
		if err := rows.Scan(
			&bot.Id,
			&bot.DomainId,
			&bot.CreatedAt,
			&bot.UpdatedAt,
			&bot.FlowId,
			&bot.DisplayName,
		); err != nil {
			return nil, err
		}
		
		bots = append(bots, bot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bots, nil
}

func (b *botStore) Update(ctx context.Context, updateCmd *dto.UpdateBotCommand) (*model.WebitelBot, error) {
	var (
		query = `
			update im_contact.bot
			set
				updated_at = now(),
				flow_id = coalesce(@FlowId, flow_id),
				display_name = coalesce(@DisplayName, display_name)
			where domain_id = @DomainId
				and id = @Id
			returning
				id,
				domain_id,
				created_at,
				updated_at,
				flow_id,
				display_name
		`
		args = pgx.NamedArgs{
			"Id": updateCmd.Id,
			"FlowId": updateCmd.FlowId,
			"DisplayName": updateCmd.DisplayName,
			"DomainId": updateCmd.DomainId,
		}
		bot model.WebitelBot
	)
	
	if err := b.db.Master().QueryRow(ctx, query, args).Scan(
		&bot.Id, &bot.DomainId, &bot.CreatedAt, &bot.UpdatedAt, &bot.FlowId, &bot.DisplayName,
		); err != nil {
		return nil, err
	} 

	return &bot, nil
}

func (b *botStore) Delete(ctx context.Context, deleteCmd *dto.DeleteBotCommand) error {
	var (
		query = `
			delete from im_contact.bot
			where domain_id = @DomainId
				and (@Id::uuid is null or id = @Id::uuid)
				and (@FlowId::int is null or flow_id = @FlowId::int)
		`
		args = pgx.NamedArgs{
			"DomainId": deleteCmd.DomainId,
			"Id": deleteCmd.Id,
			"FlowId": deleteCmd.FlowId,
		}
	)

	if _, err := b.db.Master().Exec(ctx, query, args); err != nil {
		return err
	}
	
	return nil
}

