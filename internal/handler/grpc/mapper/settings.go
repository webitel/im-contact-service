package mapper

import (
	"github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/model"
)

//go:generate goverter gen github.com/webitel/im-contact-service/internal/handler/grpc/mapper

// goverter:converter
// goverter:matchIgnoreCase
// goverter:extend github.com/google/uuid:Parse
// goverter:extend time:UnixMilli
// goverter:extend ConvertInUserFilter
type SettingsInConverter interface {
	ConvertGetSettingsRequest(*contact.GetContactSettingsRequest) (*model.GetContactSettingsRequest, error)
	// goverter:useZeroValueOnPointerInconsistency
	ConvertUpdateSettingsRequest(*contact.UpdateContactSettingsRequest) (*model.UpdateContactSettingsRequest, error)
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




