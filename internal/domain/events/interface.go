package events

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	Topic() string
	OccurredAt() time.Time
	EntityID() uuid.UUID
}

// Base provides shared logic for all domain events.
type Base struct {
	ID        uuid.UUID `json:"-"`
	TopicName string    `json:"-"`
	Timestamp time.Time `json:"occurred_at"`
}

func (b Base) Topic() string         { return b.TopicName }
func (b Base) OccurredAt() time.Time { return b.Timestamp }
func (b Base) EntityID() uuid.UUID   { return b.ID }
