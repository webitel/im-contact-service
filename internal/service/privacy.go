package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/store"
	"github.com/webitel/webitel-go-kit/pkg/errors"
)

type contactPrivacyService struct {
	logger        *slog.Logger
	settingsStore store.SettingsStore
	contactStore  store.ContactStore
}


func NewContactPrivacyService(log *slog.Logger, settingsStore store.SettingsStore, contactStore store.ContactStore) (ContactPrivacyService, error) {
	return &contactPrivacyService{logger: log, settingsStore: settingsStore, contactStore: contactStore}, nil
}

type ValidationFunc func(from, to *model.Contact, toSettings *model.ContactSettings) error

var (
	inviteValidators = []ValidationFunc{
		ensureSharedDomain,
		denyBotToBotCommunication,
		validateAllowInviteFrom,
	}
	sendValidators = []ValidationFunc{
		ensureSharedDomain,
		denyBotToBotCommunication,
	}
)

func (s *contactPrivacyService) CanSend(ctx context.Context, request *model.CanSendRequest) error {
	if request == nil {
		return errors.InvalidArgument("request required")
	}

	from, to := request.From, request.To
	if from == uuid.Nil {
		return errors.InvalidArgument("from required")
	}
	if to == uuid.Nil {
		return errors.InvalidArgument("to required")
	}

	fromContact, toContact, err := s.findContactPair(ctx, from, to)
	if err != nil {
		return err
	}

	settingsTo, err := s.settingsStore.Get(ctx, to)
	if err != nil {
		return err
	}

	err = s.validateCanSend(fromContact, toContact, settingsTo)
	if err != nil {
		return errors.Forbidden("sending forbidden", errors.WithCause(err))
	}

	return nil
}

func (s *contactPrivacyService) CanInvite(ctx context.Context, request *model.CanInviteRequest) error {
	if request == nil {
		return errors.InvalidArgument("request required")
	}

	from, to := request.From, request.To
	if from == uuid.Nil {
		return errors.InvalidArgument("from required")
	}
	if to == uuid.Nil {
		return errors.InvalidArgument("to required")
	}

	fromContact, toContact, err := s.findContactPair(ctx, from, to)
	if err != nil {
		return err
	}

	settingsTo, err := s.settingsStore.Get(ctx, to)
	if err != nil {
		return err
	}

	err = s.validateCanInvite(fromContact, toContact, settingsTo)
	if err != nil {
		return errors.Forbidden("inviting forbidden", errors.WithCause(err))
	}

	return nil
}

func (s *contactPrivacyService) findContactPair(ctx context.Context, from, to uuid.UUID) (fromContact, toContact *model.Contact, err error) {
	contacts, err := s.contactStore.Search(ctx, &model.ContactSearchRequest{IDs: []uuid.UUID{from, to}})
	if err != nil {
		return nil, nil, err
	}

	foundLen := len(contacts)
	if foundLen <=0 || foundLen > 2 {
		return nil, nil, errors.InvalidArgument("invalid contacts found")
	}

	for _, contact := range contacts {
		if fromContact != nil && toContact != nil {
			break
		}
		if contact == nil {
			continue
		}

		
		if contact.ID == from {
			fromContact = contact
		}

		if contact.ID == to {
			toContact = contact
		}
	}

	if fromContact == nil || toContact == nil {
		return nil, nil, errors.InvalidArgument("can't find contacts to validate")
	}

	return fromContact, toContact, nil

}

func (s *contactPrivacyService) checkValidationRules(from, to *model.Contact, toSettings *model.ContactSettings, validators []ValidationFunc) error {
	for _, validate := range validators {
		err := validate(from, to, toSettings)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *contactPrivacyService) validateCanInvite(fromContact, toContact *model.Contact, toSettings *model.ContactSettings) error {
	return s.checkValidationRules(fromContact, toContact, toSettings, inviteValidators)
}

func (s *contactPrivacyService) validateCanSend(fromContact, toContact *model.Contact, toSettings *model.ContactSettings) error {
	return s.checkValidationRules(fromContact, toContact, toSettings, sendValidators)
}

func validateAllowInviteFrom(from, to *model.Contact, toSettings *model.ContactSettings) error {
	if toSettings == nil {
		return errors.InvalidArgument("receiver settings required")
	}

	userFilter := toSettings.AllowInvitesFrom

	allow := userFilter.InFilter(from, to)

	if !allow {
		return errors.Forbidden("receiver privacy settings forbid sending invites")
	}

	return nil
}

func ensureSharedDomain(from *model.Contact, to *model.Contact, _ *model.ContactSettings) error {
	if from == nil || to == nil {
		return errors.InvalidArgument("contacts required")
	}

	if from.DomainID != to.DomainID {
		return errors.InvalidArgument("contacts should share a domain")
	}

	return nil

}

func denyBotToBotCommunication(from *model.Contact, to *model.Contact, _ *model.ContactSettings) error {
	if from == nil || to == nil {
		return errors.InvalidArgument("contacts required")
	}

	if from.IsBot && to.IsBot {
		return errors.Forbidden("bot to bot communication is forbidden")
	}

	return nil

}