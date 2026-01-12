package store

import (
	"context"

	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
)

type (
	Store        interface{}
	ContactStore interface {
		Search(ctx context.Context, filter *dto.ContactSearchFilter) ([]*model.Contact, error)
		Create(ctx context.Context, contact *model.Contact) (*model.Contact, error)
		Update(ctx context.Context, updater *dto.UpdateContactCommand) (*model.Contact, error)
		Delete(ctx context.Context, command *dto.DeleteContactCommand) error
		ClearByDomain(ctx context.Context, domainId int) error
	}

	BotStore interface {
		Create(ctx context.Context, bot *model.WebitelBot) (*model.WebitelBot, error)
		Search(ctx context.Context, filter *dto.SearchBotRequest) ([]*model.WebitelBot, error)
		Update(ctx context.Context, updateCmd *dto.UpdateBotCommand) (*model.WebitelBot, error)
		Delete(ctx context.Context, deleteCmd *dto.DeleteBotCommand) error		
	}
)
