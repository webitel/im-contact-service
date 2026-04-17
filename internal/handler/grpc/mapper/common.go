package mapper

import (
	"time"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/model"
)




func ConvertUUID(id uuid.UUID) string {
	return id.String()
}

func ConvertOptionalUUID(in *string) (uuid.UUID, error) {
	if in == nil {
		return uuid.Nil, nil
	}

	return uuid.Parse(*in)
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

func ConvertIntToInt64(in int) int64 {
	return int64(in)
}


func ConvertInt64ToInt(in int64) int {
	return int(in)
}


func ConvertInt32ToInt(in int32) int {
	return int(in)
}