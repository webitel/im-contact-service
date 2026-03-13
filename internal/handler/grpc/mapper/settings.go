package mapper

import (
	"time"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
)

//go:generate goverter gen github.com/webitel/im-contact-service/internal/handler/grpc/mapper

// goverter:converter
// goverter:matchIgnoreCase
// goverter:extend github.com/google/uuid:Parse
// goverter:extend time:UnixMilli
// goverter:extend ConvertInUserFilter
type SettingsInConverter interface {
	ConvertGetSettingsRequest(*contact.GetContactSettingsRequest) (*dto.GetContactSettingsRequest, error)
	// goverter:useZeroValueOnPointerInconsistency
	ConvertUpdateSettingsRequest(*contact.UpdateContactSettingsRequest) (*dto.UpdateContactSettingsRequest, error)
}

// goverter:converter
// goverter:matchIgnoreCase
// goverter:ignoreUnexported
// goverter:extend ConvertUUID
// goverter:extend ConvertTimeToInt64
// goverter:extend ConvertOutUserFilter
type SettingsOutConverter interface {
	ConvertSettings(*model.ContactSettings) (*contact.Settings, error)
}


func ConvertUUID(id uuid.UUID) string {
	return id.String()
}


func ConvertTimeToInt64(in time.Time) int64 {
	return in.UnixMilli()
}


func ConvertInUserFilter(in contact.UserFilter) model.UserFilter {
	return model.UserFilter(in)
}
func ConvertOutUserFilter(in model.UserFilter) contact.UserFilter {
	return contact.UserFilter(in)
}

