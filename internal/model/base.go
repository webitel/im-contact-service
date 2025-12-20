package model

import "github.com/google/uuid"

type BaseModel struct {
	Id       uuid.UUID `json:"id"`
	DomainId int       `json:"domain_id"`

	CreatedAt uint64 `json:"created_at"`
	UpdatedAt uint64 `json:"updated_at"`
}
