package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/store"
	"github.com/webitel/webitel-go-kit/pkg/errors"
)

var _ ContactSettingsService = &contactSettingsService{}

type contactSettingsService struct {
	logger        *slog.Logger
	settingsStore store.SettingsStore
	contactStore  store.ContactStore
}

func NewContactSettingService(log *slog.Logger, store store.SettingsStore) (ContactSettingsService, error) {
	return &contactSettingsService{logger: log, settingsStore: store}, nil
}

func (s *contactSettingsService) Get(ctx context.Context, req *model.GetContactSettingsRequest) (*model.ContactSettings, error) {
	if req == nil {
		return nil, errors.InvalidArgument("get settings request is required")
	}
	if req.ContactID == uuid.Nil {
		return nil, errors.InvalidArgument("contact id required to get settings")
	}
	if req.InitiatorContactID != uuid.Nil && req.InitiatorContactID != req.ContactID {
		return nil, errors.Forbidden("contact can get only own settings")
	}
	return s.settingsStore.Get(ctx, req.ContactID)
}

func (s *contactSettingsService) Update(ctx context.Context, request *model.UpdateContactSettingsRequest) (*model.ContactSettings, error) {
	if request == nil {
		return nil, errors.InvalidArgument("update settings request is required")
	}
	if request.ContactID == uuid.Nil {
		return nil, errors.InvalidArgument("contact id required to update settings")
	}
	if request.InitiatorContactID != uuid.Nil && request.InitiatorContactID != request.ContactID {
		return nil, errors.Forbidden("contact can change only own settings")
	}
	return s.settingsStore.Update(ctx, request)
}

func (s *contactSettingsService) Create(ctx context.Context, request *model.CreateContactSettingsRequest) (*model.ContactSettings, error) {
	if request == nil {
		return nil, errors.InvalidArgument("update settings request is required")
	}
	if request.Settings == nil {
		return nil, errors.InvalidArgument("settings is required to create settings")
	}
	if request.ContactID == uuid.Nil {
		return nil, errors.InvalidArgument("contact id required to update settings")
	}
	return s.settingsStore.Create(ctx, request)
}
