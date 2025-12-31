package amqp

import (
	"context"

	"github.com/webitel/im-contact-service/internal/domain/events"
)

type DomainDeletedEventHandler interface {
	DeleteByDomain(ctx context.Context, domainId int) error
}

type MessageHandler struct {
	service DomainDeletedEventHandler
}

func NewMessageHandler(svc DomainDeletedEventHandler) *MessageHandler {
	return &MessageHandler{service: svc}
}

func (h *MessageHandler) OnDomainDeleted(ctx context.Context, event events.DomainDeleted) error {
	return h.service.DeleteByDomain(ctx, event.DomainID())
}
