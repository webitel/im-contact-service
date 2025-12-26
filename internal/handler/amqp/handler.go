package amqp

import (
	"context"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/domain/events"
)

type EventSubscriber interface {
	DeleteByDomain(ctx context.Context, id uuid.UUID) error
}

type MessageHandler struct {
	service EventSubscriber
}

func NewMessageHandler(svc EventSubscriber) *MessageHandler {
	return &MessageHandler{service: svc}
}

func (h *MessageHandler) OnDomainDeleted(ctx context.Context, event events.DomainDeleted) error {
	return h.service.DeleteByDomain(ctx, event.EntityID())
}
