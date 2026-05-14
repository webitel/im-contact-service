package service

import (
	"context"
	"log/slog"

	"github.com/webitel/webitel-go-kit/pkg/errors"

	"github.com/webitel/im-contact-service/internal/domain/events"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/store"
)

type via struct {
	publisher          EventPublisher
	communicationStore store.ViaStore
	logger             *slog.Logger
}

func newCommunication(logger *slog.Logger, communicationStore store.ViaStore, publisher EventPublisher) *via {
	log := logger.With("component", "service.communication")

	return &via{logger: log, communicationStore: communicationStore, publisher: publisher}
}

func (communicationService *via) Create(ctx context.Context, communication *model.ViaCommunication) (*model.ViaCommunication, error) {
	log := communicationService.logger.With("operation", "create")

	if err := communication.Validate(); err != nil {
		log.Warn("validation error: ", "error", err)

		return nil, err
	}

	savedCommunication, err := communicationService.communicationStore.Create(ctx, communication)
	if err != nil {
		log.Error(
			"saving contact communication",
			"error", err,
			"contact_id", communication.ContactID.String(),
			"via", communication.Via,
		)

		return nil, errors.Wrap(err, errors.WithID("service.communication.create"))
	}

	if err = communicationService.publisher.Publish(ctx, events.NewViaCreatedEvent(savedCommunication)); err != nil {
		log.Error("publishing via created event", "error", err, "contact_id", communication.ContactID.String(), "via", communication.Via)

		return savedCommunication, errors.Internal("publishing via created event", errors.WithCause(err), errors.WithID("service.communication.create"))
	}

	return savedCommunication, nil
}

func (communicationService *via) Update(ctx context.Context, communication *model.ViaCommunication) (*model.ViaCommunication, error) {
	log := communicationService.logger.With("operation", "update")

	if err := communication.Validate(); err != nil {
		log.Warn("validation error: ", "error", err)

		return nil, err
	}

	updatetCommunication, err := communicationService.communicationStore.Update(ctx, communication)
	if err != nil {
		log.Error(
			"full updating contact communication",
			"error", err,
			"contact_id", communication.ContactID.String(),
			"via", communication.Via,
		)

		return nil, err
	}

	if err = communicationService.publisher.Publish(ctx, events.NewViaUpdatedEvent(updatetCommunication)); err != nil {
		log.Error("publishing via updated event", "error", err, "contact_id", updatetCommunication.ContactID, "new_via", updatetCommunication.Via, "old_via", communication.Via)

		return updatetCommunication, errors.Internal("publishing updaed via event", errors.WithCause(err), errors.WithID("service.communication.update"))
	}

	return updatetCommunication, nil
}

func (communicationService *via) PartialUpdate(ctx context.Context, updateCommand *model.CommunicationViaPartialUpdateCmd) (*model.ViaCommunication, error) {
	log := communicationService.logger.With("operation", "partial_update")

	if err := updateCommand.Validate(); err != nil {
		log.Warn("validating partial update command", "error", err)

		return nil, err
	}

	updated, err := communicationService.communicationStore.PartialUpdate(ctx, updateCommand)
	if err != nil {
		log.Error("partial updating contact communication", "error", err)

		return nil, err
	}

	if err = communicationService.publisher.Publish(ctx, events.NewViaUpdatedEvent(updated)); err != nil {
		log.Error("publishing via partially updated event", "error", err, "contact_id", updated.ContactID, "via", updated.Via)

		return updated, errors.Internal("publishing via partially updated event", errors.WithCause(err), errors.WithID("service.communication.partial_update"))
	}

	return updated, nil
}

func (communicationService *via) Search(ctx context.Context, filter *model.SearchViaCommunicationsFilter) ([]*model.ViaCommunication, error) {
	log := communicationService.logger.With("operation", "search")

	if err := filter.Validate(); err != nil {
		log.Warn("validating search communication option request", "error", err)

		return nil, err
	}

	records, err := communicationService.communicationStore.Search(ctx, filter)
	if err != nil {
		log.Error("retrieving records from store", "error", err)

		return nil, err
	}

	return records, nil
}
