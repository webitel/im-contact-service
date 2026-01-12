package service

import (
	"context"

	"log/slog"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/store"
	"github.com/webitel/webitel-go-kit/pkg/errors"
)

var (
	ensureValidationError = errors.InvalidArgument("ensure bot validation violation!")
	updateValidationError = errors.InvalidArgument("update bot validation violation!")
	deleteValidationError = errors.InvalidArgument("delete bot validation violation!")
	createValidationError = errors.InvalidArgument("create bot validation violation!")
)

var (
	_ BotManager = (*baseBotManager)(nil)
)

type (
	BotManager interface {
		Create(ctx context.Context, bot *model.WebitelBot) (*model.WebitelBot, error)
		Search(ctx context.Context, filter *dto.SearchBotRequest) ([]*model.WebitelBot, bool, error)
		Update(ctx context.Context, updateCmd *dto.UpdateBotCommand) (*model.WebitelBot, error)
		Delete(ctx context.Context, deleteCmd *dto.DeleteBotCommand) error

		EnsureBot(ctx context.Context, ensureRequest *dto.EnsureBotRequest) (*model.WebitelBot, error)
	}

	baseBotManager struct {
		logger *slog.Logger
		botStore store.BotStore
	}
)

func NewBaseBotManager(logger *slog.Logger, botStore store.BotStore) *baseBotManager {
	return &baseBotManager{
		logger: logger,
		botStore: botStore,
	}
}

func (m *baseBotManager) Create(ctx context.Context, bot *model.WebitelBot) (*model.WebitelBot, error) {
	if bot.DomainId <= 0 || bot.FlowId <= 0 {
		return nil, createValidationError
	}

	return m.botStore.Create(ctx, bot)
}

func (m *baseBotManager) Search(ctx context.Context, filter *dto.SearchBotRequest) ([]*model.WebitelBot, bool, error) {
	var (
		err error
		botsList []*model.WebitelBot
		hasNext bool
	)

	if botsList, err = m.botStore.Search(ctx, filter); err != nil {
		return nil, hasNext, err
	}

	if len(botsList) > filter.Size {
		hasNext = true
		botsList = botsList[:filter.Size - 1]
	}

	return botsList, hasNext, nil
}

func (m *baseBotManager) Update(ctx context.Context, updateCmd *dto.UpdateBotCommand) (*model.WebitelBot, error) {
	if !isUpdateCommandValid(updateCmd) {
		m.logger.Warn("[BOT_MANAGER] update bot command validation violation!")
		return nil, updateValidationError
	}
	
	return m.botStore.Update(ctx, updateCmd)
}

func (m *baseBotManager) Delete(ctx context.Context, deleteCmd *dto.DeleteBotCommand) error {
	if !isDeleteCommandValid(deleteCmd) {
		return deleteValidationError
	}
	
	return m.botStore.Delete(ctx, deleteCmd)
}

func (m *baseBotManager) EnsureBot(ctx context.Context, ensureRequest *dto.EnsureBotRequest) (*model.WebitelBot, error) {
	var (
		err error
		bot *model.WebitelBot
	)

	if !isEnsureValid(ensureRequest) {
		m.logger.Warn("[BOT_MANAGER] ensure bot request params validation violation", "ensure_request", ensureRequest)
		return nil, ensureValidationError
	}

	searchBot := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: ensureRequest.DomainId,
			Size: 1,
		},
		FlowIds: []int64{int64(ensureRequest.FlowId)},
	}

	bots, err := m.botStore.Search(ctx, searchBot)
	if err != nil {
		m.logger.Error("[BOT_MANAGER] STORE ERROR SEARCHING BOT", "err", err, "domain_id", ensureRequest.DomainId, "flow_id", ensureRequest.FlowId)
		return nil, err
	}

	// FOUND!
	if len(bots) > 0 {
		bot = bots[0]
		return bot, nil
	}

	bot = &model.WebitelBot{
		BaseModel: model.BaseModel{
			DomainId: ensureRequest.DomainId,
		},
		FlowId: ensureRequest.FlowId,
	}

	// OTHERWISE TRY TO CREATE NEW BOT! 
	if bot, err = m.botStore.Create(ctx, bot); err != nil {
		m.logger.Error("[BOT_MANAGER] STORE ERROR CREATING BOT", "err", err, "domain_id", ensureRequest.DomainId, "flow_id", ensureRequest.FlowId)
		return nil, err
	}

	return bot, nil
}

func isUpdateCommandValid(updateCmd *dto.UpdateBotCommand) bool {
	if updateCmd == nil {
		return false
	}

	if updateCmd.Id == uuid.Nil || updateCmd.DomainId <= 0 {
		return false
	}

	return true
}

func isEnsureValid(ensureRequest *dto.EnsureBotRequest) bool {
	if ensureRequest == nil {
		return false
	}

	if ensureRequest.DomainId <= 0 || ensureRequest.FlowId <= 0 {
		return false
	}

	return true
}


func isDeleteCommandValid(deleteCmd *dto.DeleteBotCommand) bool {
	if deleteCmd == nil {
		return false
	}

	if deleteCmd.DomainId <= 0 {
		return false
	} 

	if deleteCmd.Id == nil || deleteCmd.FlowId == nil {
		return false
	}

	return true
}
