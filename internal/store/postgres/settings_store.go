package postgres

import (
	"context"
	"encoding/json"
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
	settingsEncoded, err := encodeRules(command.Settings.Rules) 
	if err != nil {
		return nil, err
	}

	_, err = s.db.Master().Exec(
		ctx,
		`INSERT INTO im_contact.contact_setting(contact_id, rules) VALUES ($1, $2)`,
		command.ContactID,
		settingsEncoded,
	)
	if err != nil {
		return nil, err
	}
	return command.Settings, nil
	
}


func encodeRules(rules []model.SettingRule) ([]byte, error) {
	if len(rules) == 0 {
		return nil, errors.InvalidArgument("settings is required to encode settings")
	}
	var (
		result = map[string]any{}
	)


	for _, rule := range rules {
		result[rule.GetType()] = rule.GetValue()
	}


	return json.Marshal(result)
}



func decodeRules(encoded []byte) ([]model.SettingRule, error) {
	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, err
	}
	var rules []model.SettingRule
	for ruleName, value := range decoded {
		rule, err := model.BuildRule(ruleName)
		if err != nil {
			continue
		}
		if err := rule.SetValue(value); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// Get implements [store.SettingsStore].
func (s *SettingsStore) Get(ctx context.Context, contactID uuid.UUID) (*model.ContactSettings, error) {
	if contactID == uuid.Nil {
		return nil, errors.InvalidArgument("contact id required to get settings")
	}
	var (
		settings model.ContactSettings
		rulesEncoded []byte
	)
	err := s.db.Master().QueryRow(ctx, "SELECT id, updated_at, contact_id, rules FROM im_contact.contact_setting WHERE contact_id = $1", contactID).Scan(
		&settings.ID,
		&settings.UpdatedAt,
		&settings.ContactID,
		&rulesEncoded,
	)
	if err != nil {
		return nil, err
	}
	
	rules, err := decodeRules(rulesEncoded)
	if err != nil {
		return nil, err
	}
	settings.Rules = rules

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
	rules, err := encodeRules(args.Settings.Rules) 
	if err != nil {
		return nil, err
	}

	var updatedSettings model.ContactSettings
	err = pgxscan.Get(
		ctx,
		s.db.Master(),
		&updatedSettings,
		`UPDATE im_contact.contact_setting SET rules= $1, updated_at = now() WHERE contact_id = $3 RETURNING id, updated_at, contact_id, rules`,
		rules,
		args.ContactID,
	)
	if err != nil {
		return nil, err
	}
	return &updatedSettings, nil
}
