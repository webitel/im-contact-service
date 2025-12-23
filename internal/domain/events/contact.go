package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/domain/model"
)

const (
	ContactCreatedTopic = "contact.created"
	ContactUpdatedTopic = "contact.updated"
	ContactDeletedTopic = "contact.deleted"
)

type ContactCreated struct {
	Base
	ContactID     uuid.UUID `json:"contact_id"`
	Name          string    `json:"name"`
	Username      string    `json:"username"`
	Type          string    `json:"type"`
	ApplicationID string    `json:"application_id"`
	IssuerID      string    `json:"issuer_id"`
}

func NewContactCreated(m *model.Contact) *ContactCreated {
	return &ContactCreated{
		Base: Base{
			ID:        m.Id,
			TopicName: ContactCreatedTopic,
			Timestamp: m.CreatedAt,
		},
		ContactID:     m.Id,
		Name:          m.Name,
		Username:      m.Username,
		Type:          m.Type,
		ApplicationID: m.ApplicationId,
		IssuerID:      m.IssuerId,
	}
}

type ContactUpdated struct {
	Base
	ContactID uuid.UUID `json:"contact_id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Type      string    `json:"type"`
}

func NewContactUpdated(m *model.Contact) *ContactUpdated {
	return &ContactUpdated{
		Base: Base{
			ID:        m.Id,
			TopicName: ContactUpdatedTopic,
			Timestamp: m.UpdatedAt,
		},
		ContactID: m.Id,
		Name:      m.Name,
		Username:  m.Username,
		Type:      m.Type,
	}
}

type ContactDeleted struct {
	Base
	ContactID uuid.UUID `json:"contact_id"`
}

func NewContactDeleted(id uuid.UUID) *ContactDeleted {
	return &ContactDeleted{
		Base: Base{
			ID:        id,
			TopicName: ContactDeletedTopic,
			Timestamp: time.Now().UTC(),
		},
		ContactID: id,
	}
}
