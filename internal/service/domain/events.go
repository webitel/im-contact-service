package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/model"
)

// DomainEvent is the contract for all events in the system.
type DomainEvent interface {
	Topic() string
	OccurredAt() int64
	EntityID() uuid.UUID
}

const (
	ContactCreatedTopic = "contact.created"
	ContactUpdatedTopic = "contact.updated"
	ContactDeletedTopic = "contact.deleted"
)

// --- ContactCreatedEvent ---
type ContactCreatedEvent struct {
	ContactID     uuid.UUID         `json:"contact_id"`
	Name          string            `json:"name"`
	Username      string            `json:"username"`
	Type          model.ContactType `json:"type"`
	ApplicationID uuid.UUID         `json:"application_id"`
	IssuerID      uuid.UUID         `json:"issuer_id"`
	Timestamp     uint64            `json:"occurred_at"`
}

func (e *ContactCreatedEvent) Topic() string       { return ContactCreatedTopic }
func (e *ContactCreatedEvent) OccurredAt() int64   { return int64(e.Timestamp) }
func (e *ContactCreatedEvent) EntityID() uuid.UUID { return e.ContactID }

func NewContactCreatedEvent(m *model.Contact) *ContactCreatedEvent {
	return &ContactCreatedEvent{
		ContactID:     m.Id,
		Name:          m.Name,
		Username:      m.Username,
		Type:          m.Type,
		ApplicationID: m.ApplicationId,
		IssuerID:      m.IssuerId,
		Timestamp:     m.CreatedAt,
	}
}

// --- ContactUpdatedEvent ---
type ContactUpdatedEvent struct {
	ContactID uuid.UUID         `json:"contact_id"`
	Name      string            `json:"name"`
	Username  string            `json:"username"`
	Type      model.ContactType `json:"type"`
	Timestamp uint64            `json:"occurred_at"`
}

func (e *ContactUpdatedEvent) Topic() string       { return ContactUpdatedTopic }
func (e *ContactUpdatedEvent) OccurredAt() int64   { return int64(e.Timestamp) }
func (e *ContactUpdatedEvent) EntityID() uuid.UUID { return e.ContactID }

func NewContactUpdatedEvent(m *model.Contact) *ContactUpdatedEvent {
	return &ContactUpdatedEvent{
		ContactID: m.Id,
		Name:      m.Name,
		Username:  m.Username,
		Type:      m.Type,
		Timestamp: m.UpdatedAt,
	}
}

// --- ContactDeletedEvent ---
type ContactDeletedEvent struct {
	ContactID uuid.UUID `json:"contact_id"`
	Timestamp int64     `json:"occurred_at"`
}

func (e *ContactDeletedEvent) Topic() string       { return ContactDeletedTopic }
func (e *ContactDeletedEvent) OccurredAt() int64   { return e.Timestamp }
func (e *ContactDeletedEvent) EntityID() uuid.UUID { return e.ContactID }

func NewContactDeletedEvent(id uuid.UUID) *ContactDeletedEvent {
	return &ContactDeletedEvent{
		ContactID: id,
		Timestamp: time.Now().UnixMilli(),
	}
}
