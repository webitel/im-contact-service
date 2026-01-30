package mapper

import (
	"github.com/google/uuid"
	impb "github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
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
		Subject: contact.SubjectId,
	}, nil
}


func CanSendRequest2Model(request *impb.CanSendRequest) *dto.CanSendQuery {
	// Checked uuid validity in protobuf layer
	from,_:=uuid.Parse(request.From)
	to,_:=uuid.Parse(request.To)
	

	canSendQuery := &dto.CanSendQuery{
		DomainId: int(request.GetDomainId()),
		From: model.Peer{
			Id:from,
		},
		To:  model.Peer{
			Id:to,
		},
	}

	return canSendQuery
}