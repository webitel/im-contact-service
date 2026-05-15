package model

import "github.com/google/uuid"

type Contact struct {
	BaseModel `json:"base_model"`

	IssuerID      string `json:"issuer_id" db:"issuer_id"`
	SubjectID     string `json:"subject_id" db:"subject_id"`
	ApplicationID string `json:"application_id" db:"application_id"`
	Type          string `json:"type" db:"type"`

	Name     string            `json:"name" db:"name"`
	Username string            `json:"username" db:"username"`
	Metadata map[string]string `json:"metadata" db:"metadata"`
	IsBot    bool              `jsob:"is_bot" db:"is_bot" json:"is_bot"`

	Via []*ViaCommunication `json:"via" db:"via"`
}

func (c *Contact) TableName() string { return "im_contact.contact" }

func (c *Contact) Equal(compare *Contact) bool {
	if c == nil && compare == nil {
		return true
	}

	return c.ApplicationID == compare.ApplicationID &&
		c.ID == compare.ID &&
		c.DomainID == compare.DomainID &&
		c.IssuerID == compare.IssuerID &&
		c.Name == compare.Name &&
		c.Type == compare.Type &&
		c.Username == compare.Username &&
		c.IsBot == compare.IsBot
}

func ContactAllowedFields() []string {
	return []string{
		"issuer_id", "application_id", "type", "name", "username", "metadata",
		"id", "domain_id", "created_at", "updated_at", "subject_id", "is_bot", "via",
	}
}

type (
	ContactSearchRequest struct {
		DomainID *int        `json:"domain_id"`
		IDs      []uuid.UUID `json:"ids"`
		Page     int32       `json:"page"`
		Size     int32       `json:"size"`
		Q        *string     `json:"q"`
		Sort     string      `json:"sort"`
		Fields   []string    `json:"fields"`
		Apps     []string    `json:"apps"`
		Issuers  []string    `json:"issuers"`
		Types    []string    `json:"types"`
		Subjects []string    `json:"subjects"`
		OnlyBots *bool       `json:"is_bot"`
	}

	UpdateContactRequest struct {
		ID       uuid.UUID         `json:"id"`
		DomainID int               `json:"domain_id"`
		Name     *string           `json:"name"`
		Username *string           `json:"username"`
		Metadata map[string]string `json:"metadata"`
		Subject  string            `json:"subject"`
	}

	PartialUpdateContactRequest struct {
		ID       uuid.UUID         `json:"id"`
		DomainID int               `json:"domain_id"`
		Name     string            `json:"name"`
		Username string            `json:"username"`
		Metadata map[string]string `json:"md"`
		Subject  string            `json:"sub"`
		Fields   []string
	}

	CanSendRequest struct {
		From uuid.UUID `json:"from"`
		To   uuid.UUID `json:"to"`
	}

	CanInviteRequest struct {
		From uuid.UUID
		To   uuid.UUID
	}

	DeleteContactRequest struct {
		DomainID int       `json:"domain_id"`
		ID       uuid.UUID `json:"id"`
	}

	CreateContactRequest struct {
		IssuerID      uuid.UUID         `json:"issuer_id"`
		SubjectID     string            `json:"subject_id"`
		ApplicationID uuid.UUID         `json:"application_id"`
		Type          string            `json:"type"`
		Name          string            `json:"name"`
		Username      string            `json:"username"`
		Metadata      map[string]string `json:"metadata"`
	}
)
