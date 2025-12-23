package model

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	Id       uuid.UUID `json:"id"`
	DomainId int       `json:"domain_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
