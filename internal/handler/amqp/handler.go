package amqp

import (
	"context"

	"github.com/webitel/im-contact-service/internal/domain/events"
)

type DomainEventsHandler interface {
	DeleteByDomain(ctx context.Context, domainId int) error
	DeleteBotByFlowID(ctx context.Context, flowID string) error
}

type MessageHandler struct {
	service DomainEventsHandler
}

func NewMessageHandler(svc DomainEventsHandler) *MessageHandler {
	return &MessageHandler{service: svc}
}

func (h *MessageHandler) OnDomainDeleted(ctx context.Context, event events.DomainDeleted) error {
	return h.service.DeleteByDomain(ctx, event.DomainID())
}

func (h *MessageHandler) OnFlowSchemaDelete(ctx context.Context, event events.FlowSchemaDeleted) error {
	return h.service.DeleteBotByFlowID(ctx, event.FlowID)
}
