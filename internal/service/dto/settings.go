package dto

import (
	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/internal/domain/model"
)




type GetContactSettingsRequest struct {
	ContactID uuid.UUID
}

type UpdateContactSettingsRequest struct {
	ContactID uuid.UUID
	Settings   *model.ContactSettings
}


type CreateContactSettingsRequest struct {
	ContactID uuid.UUID
	Settings   *model.ContactSettings
}


