package mapper

import (
	impb "github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/model"
)

//go:generate goverter gen github.com/webitel/im-contact-service/internal/handler/grpc/mapper

// goverter:converter
// goverter:matchIgnoreCase
// goverter:extend ConvertInt64ToInt
// goverter:extend github.com/google/uuid:Parse
type PrivacyInMapper interface {
	ConvertCanSendRequest(*impb.CanSendRequest) (*model.CanSendRequest, error)
	ConvertCanInviteRequest(*impb.CanInviteRequest) (*model.CanInviteRequest, error)
}