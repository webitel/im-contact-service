package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	DomainDeletedTopic = "domains.delete.#"
)

type DomainDeleted struct {
	_         uuid.UUID
	DomainId  int       `json:"domain_id"`
	Timestamp time.Time `json:"timestamp"`
}

var _ Event = (*DomainDeleted)(nil)

func (e DomainDeleted) DomainID() int         { return e.DomainId }
func (e DomainDeleted) Topic() string         { return DomainDeletedTopic }
func (e DomainDeleted) OccurredAt() time.Time { return e.Timestamp }
func (e DomainDeleted) EntityID() uuid.UUID {
	// Placeholder to satisfy Event interface; actual DomainId is int
	return uuid.Nil
}
