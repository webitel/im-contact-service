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
	DomainID  int       `json:"domain_id"`
	Timestamp time.Time `json:"timestamp"`
}

var _ Event = (*DomainDeleted)(nil)

func (e DomainDeleted) GetDomainID() int      { return e.DomainID }
func (e DomainDeleted) Topic() string         { return DomainDeletedTopic }
func (e DomainDeleted) OccurredAt() time.Time { return e.Timestamp }
func (e DomainDeleted) EntityID() uuid.UUID {
	// Placeholder to satisfy Event interface; actual DomainId is int
	return uuid.Nil
}
