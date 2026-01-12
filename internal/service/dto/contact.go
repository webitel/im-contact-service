package dto

import (
	"github.com/google/uuid"
)

type (
	ContactSearchFilter struct {
		DomainId int         `json:"domain_id"`
		Ids      []uuid.UUID `json:"ids"`
		Page     int32       `json:"page"`
		Size     int32       `json:"size"`
		Q        *string     `json:"q"`
		Sort     string      `json:"sort"`
		Fields   []string    `json:"fields"`
		Apps     []string    `json:"apps"`
		Issuers  []string    `json:"issuers"`
		Types    []string    `json:"types"`
		Subjects []string    `json:"subjects"`
	}

	UpdateContactCommand struct {
		Id       uuid.UUID         `json:"id"`
		DomainId int               `json:"domain_id"`
		Name     *string           `json:"name"`
		Username *string           `json:"username"`
		Metadata map[string]string `json:"metadata"`
		Subject  string            `json:"subject"`
	}

	CanSendQuery struct {
		From uuid.UUID `json:"from"`
		To   uuid.UUID `json:"to"`
	}

	DeleteContactCommand struct {
		DomainId int       `json:"domain_id"`
		Id       uuid.UUID `json:"id"`
	}

	CreateContactCommand struct {
		IssuerId      uuid.UUID         `json:"issuer_id"`
		SubjectId     string            `json:"subject_id"`
		ApplicationId uuid.UUID         `json:"application_id"`
		Type          string            `json:"type"`
		Name          string            `json:"name"`
		Username      string            `json:"username"`
		Metadata      map[string]string `json:"metadata"`
	}
)
