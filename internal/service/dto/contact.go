package dto

import "github.com/google/uuid"

type (
	ContactSearchFilter struct {
		Page    int32    `json:"page"`
		Size    int32    `json:"size"`
		Q       string   `json:"q"`
		Sort    string   `json:"sort"`
		Fields  []string `json:"fields"`
		Apps    []string `json:"apps"`
		Issuers []string `json:"issuers"`
		Types   []string `json:"types"`
	}

	UpdateContactCommand struct {
		Id       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Username string    `json:"username"`
		Metadata any       `json:"metadata"`
	}

	CanSendQuery struct {
		From string
		To   string
	}
)
