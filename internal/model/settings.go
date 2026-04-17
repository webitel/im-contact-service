package model

import (
	"time"

	"github.com/google/uuid"
)

type UserFilter int

const (
	All UserFilter = iota
	Nobody
	SameIssuer
)

func (u *UserFilter) InFilter(from, to *Contact) bool {
	if from == nil || to == nil {
		return false
	}
	var (
		in bool
	)
	switch *u {
	case All:
		in = true
	case Nobody:
		in = false
	case SameIssuer:
		in = from.IssuerId == to.IssuerId
	}

	return in
}

type ContactSettings struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	ContactID        uuid.UUID  `json:"contact_id" db:"contact_id"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	AllowInvitesFrom UserFilter `json:"allow_invites_from" db:"allow_invites_from"`
}

type GetContactSettingsRequest struct {
	InitiatorContactID uuid.UUID
	ContactID          uuid.UUID
}

type UpdateContactSettingsRequest struct {
	InitiatorContactID uuid.UUID
	ContactID          uuid.UUID
	AllowInvitesFrom   *UserFilter
}

type CreateContactSettingsRequest struct {
	ContactID uuid.UUID
	Settings  *ContactSettings
}
