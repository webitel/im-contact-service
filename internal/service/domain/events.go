package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/model"
)

const (
	ContactCreatedTopic = "contact.created"
	ContactUpdatedTopic = "contact.updated"
	ContactDeletedTopic = "contact.deleted"
)

type ContactCreatedEvent struct {
	ContactID     uuid.UUID         `json:"contact_id"`
	Name          string            `json:"name"`
	Username      string            `json:"username"`
	Type          model.ContactType `json:"type"`
	ApplicationID uuid.UUID         `json:"application_id"`
	IssuerID      uuid.UUID         `json:"issuer_id"`
	OccurredAt    uint64            `json:"occurred_at"`
}

// NewContactCreatedEvent creates a new event from the contact model
func NewContactCreatedEvent(m *model.Contact) *ContactCreatedEvent {
	return &ContactCreatedEvent{
		ContactID:     m.Id,
		Name:          m.Name,
		Username:      m.Username,
		Type:          m.Type,
		ApplicationID: m.ApplicationId,
		IssuerID:      m.IssuerId,
		OccurredAt:    m.CreatedAt,
	}
}

type ContactUpdatedEvent struct {
	ContactID  uuid.UUID         `json:"contact_id"`
	Name       string            `json:"name"`
	Username   string            `json:"username"`
	Type       model.ContactType `json:"type"`
	OccurredAt uint64            `json:"occurred_at"`
}

// NewContactUpdatedEvent creates an update event from the updated model
func NewContactUpdatedEvent(m *model.Contact) *ContactUpdatedEvent {
	return &ContactUpdatedEvent{
		ContactID:  m.Id,
		Name:       m.Name,
		Username:   m.Username,
		Type:       m.Type,
		OccurredAt: m.UpdatedAt,
	}
}

type ContactDeletedEvent struct {
	ContactID  uuid.UUID `json:"contact_id"`
	OccurredAt int64     `json:"occurred_at"`
}

// NewContactDeletedEvent creates a deletion event with current timestamp
func NewContactDeletedEvent(id uuid.UUID) *ContactDeletedEvent {
	return &ContactDeletedEvent{
		ContactID:  id,
		OccurredAt: time.Now().UnixMilli(),
	}
}
