package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	domainDeletedTopic = "im.admin.domain.deleted"
)

type DomainDeleted struct {
	Id        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

var _ Event = (*DomainDeleted)(nil)

func (e DomainDeleted) Topic() string         { return domainDeletedTopic }
func (e DomainDeleted) OccurredAt() time.Time { return e.Timestamp }
func (e DomainDeleted) EntityID() uuid.UUID   { return e.Id }
