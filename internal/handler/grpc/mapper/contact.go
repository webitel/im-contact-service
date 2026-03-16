package mapper

import (
	"github.com/google/uuid"
	impb "github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/model"
)

//go:generate goverter gen github.com/webitel/im-contact-service/internal/handler/grpc/mapper

// goverter:converter
// goverter:matchIgnoreCase
// goverter:extend ConvertInt32ToInt
// goverter:extend github.com/google/uuid:Parse
type ContactInConverter interface {
	// goverter:map AppId Apps
	// goverter:map IssId Issuers
	// goverter:map Type Types
	ConvertSearchRequest(*impb.SearchContactRequest) (*model.ContactSearchRequest, error)
	ConvertUpdateRequest(*impb.UpdateContactRequest) (*model.UpdateContactRequest, error)
	// goverter:map FieldMask.Paths Fields
	// goverter:useZeroValueOnPointerInconsistency
	ConvertPartialUpdateRequest(*impb.PatchContactRequest) (*model.PartialUpdateContactRequest, error)
	ConvertDeleteRequest(*impb.DeleteContactRequest) (*model.DeleteContactRequest, error)
}


func MarshalContact(contact *model.Contact) (*impb.Contact, error) {
	if contact == nil {
		return nil, nil
	}
	return &impb.Contact{
		Id:        contact.ID.String(),
		IssId:     contact.IssuerId,
		AppId:     contact.ApplicationId,
		Type:      contact.Type,
		Name:      contact.Name,
		Username:  contact.Username,
		Metadata:  contact.Metadata,
		CreatedAt: contact.CreatedAt.UnixMilli(),
		UpdatedAt: contact.UpdatedAt.UnixMilli(),
		Subject: contact.SubjectId,
		DomainId: int32(contact.DomainID),
		IsBot: contact.IsBot,
	}, nil
}


func MapPatchContactRequestToPartialUpdateContactCommand(request *impb.PatchContactRequest) *model.PartialUpdateContactRequest {
	id, _ := uuid.Parse(request.Id)
	
	return &model.PartialUpdateContactRequest{
		ID: id,
		DomainID: int(request.DomainId),
		Name: request.Name,
		Username: request.Username,
		Metadata: request.Metadata,
		Subject: request.Subject,
		Fields: request.FieldMask.Paths,
	}
}
