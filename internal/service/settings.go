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

type ContactSettingsService struct {
	logger *slog.Logger
	settingsStore store.SettingsStore
}

func (s *ContactSettingsService) Get(ctx context.Context, req *dto.GetContactSettingsRequest) (*model.ContactSettings, error) {
	if req == nil {
		return nil, errors.InvalidArgument("get settings request is required")
	}
	if req.ContactID == uuid.Nil  {
		return nil, errors.InvalidArgument("contact id required to get settings")
	}
	return s.settingsStore.Get(ctx, req.ContactID)
}


func (s *ContactSettingsService) Update(ctx context.Context, request *dto.UpdateContactSettingsRequest) (*model.ContactSettings, error) {
	if request == nil {
		return nil, errors.InvalidArgument("update settings request is required")
	}
	if request.ContactID == uuid.Nil {
		return nil, errors.InvalidArgument("contact id required to update settings")
	}
	return s.settingsStore.Update(ctx, request)
}

func (s *ContactSettingsService) Create(ctx context.Context, request *dto.CreateContactSettingsRequest) (*model.ContactSettings, error) {
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