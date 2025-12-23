package mapper

import (
	impb "github.com/webitel/im-contact-service/gen/go/api/v1"
	"github.com/webitel/im-contact-service/internal/domain/model"
)

func MarshalContact(contact *model.Contact) (*impb.Contact, error) {
	if contact == nil {
		return nil, nil
	}
	return &impb.Contact{
		Id:        contact.Id.String(),
		IssId:     contact.IssuerId,
		AppId:     contact.ApplicationId,
		Type:      contact.Type,
		Name:      contact.Name,
		Username:  contact.Username,
		Metadata:  contact.Metadata,
		CreatedAt: contact.CreatedAt.UnixMilli(),
		UpdatedAt: contact.UpdatedAt.UnixMilli(),
	}, nil
}
