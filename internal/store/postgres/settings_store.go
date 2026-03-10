package postgres

import (
	"context"
	"log/slog"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/infra/db/pg"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/store"
	"github.com/webitel/webitel-go-kit/pkg/errors"
)

var _ store.SettingsStore = (*SettingsStore)(nil)

type SettingsStore struct {
	logger *slog.Logger
	db     *pg.PgxDB
}

// Create implements [store.SettingsStore].
func (s *SettingsStore) Create(ctx context.Context, command *dto.CreateContactSettingsRequest) (*model.ContactSettings, error) {
	if command == nil {
		return nil, errors.InvalidArgument("create settings request is required")
	}
	if command.Settings == nil {
		return nil, errors.InvalidArgument("setting required to create row")
	}
	if command.ContactID == uuid.Nil {
		return nil, errors.InvalidArgument("contact id required to create settings")
	}
	_, err := s.db.Master().Exec(
		ctx,
		`INSERT INTO im_contact.contact_setting(contact_id, allow_invites_from) VALUES ($1, $2)`,
		command.ContactID,
		command.Settings.AllowInvitesFrom,
	)
	if err != nil {
		return nil, err
	}
	return command.Settings, nil
	
}


// Get implements [store.SettingsStore].
func (s *SettingsStore) Get(ctx context.Context, contactID uuid.UUID) (*model.ContactSettings, error) {
	if contactID == uuid.Nil {
		return nil, errors.InvalidArgument("contact id required to get settings")
	}
	var (
		settings model.ContactSettings
	)
	err := pgxscan.Get(
		ctx,
		s.db.Master(),
		&settings,
		"SELECT id, updated_at, contact_id, allow_invites_from FROM im_contact.contact_setting WHERE contact_id = $1",
		contactID,
	)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// Update implements [store.SettingsStore].
func (s *SettingsStore) Update(ctx context.Context, args *dto.UpdateContactSettingsRequest) (*model.ContactSettings, error) {
	if args == nil {
		return nil, errors.InvalidArgument("update settings request is required")
	}
	if args.Settings == nil {
		return nil, errors.InvalidArgument("setting required to update row")
	}
	if args.ContactID == uuid.Nil {
		return nil, errors.InvalidArgument("contact id required to update settings")
	}

	var updatedSettings model.ContactSettings
	err := pgxscan.Get(
		ctx,
		s.db.Master(),
		&updatedSettings,
		`UPDATE im_contact.contact_setting SET allow_invites_from= $1, updated_at = now() WHERE contact_id = $3 RETURNING id, updated_at, contact_id, allow_invites_from`,
		args.Settings.AllowInvitesFrom,
		args.ContactID,
	)
	if err != nil {
		return nil, err
	}
	return &updatedSettings, nil
}
