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
	ConvertUpdateRequest(*impb.UpdateContactRequest) (*model.UpdateContactRequest, error)
	// goverter:map FieldMask.Paths Fields
	// goverter:useZeroValueOnPointerInconsistency
	ConvertPartialUpdateRequest(*impb.PatchContactRequest) (*model.PartialUpdateContactRequest, error)
	ConvertDeleteRequest(*impb.DeleteContactRequest) (*model.DeleteContactRequest, error)
}

func MarshalContact(contact *model.Contact) *impb.Contact {
	if contact == nil {
		return nil
	}

	return &impb.Contact{
		Id:        contact.ID.String(),
		IssId:     contact.IssuerID,
		AppId:     contact.ApplicationID,
		Type:      contact.Type,
		Name:      contact.Name,
		Username:  contact.Username,
		Metadata:  contact.Metadata,
		CreatedAt: contact.CreatedAt.UnixMilli(),
		UpdatedAt: contact.UpdatedAt.UnixMilli(),
		Subject:   contact.SubjectID,
		DomainId:  int32(contact.DomainID),
		IsBot:     contact.IsBot,
	}
}

func MapPatchContactRequestToPartialUpdateContactCommand(request *impb.PatchContactRequest) *model.PartialUpdateContactRequest {
	id, _ := uuid.Parse(request.GetId())

	return &model.PartialUpdateContactRequest{
		ID:       id,
		DomainID: int(request.GetDomainId()),
		Name:     request.GetName(),
		Username: request.GetUsername(),
		Metadata: request.GetMetadata(),
		Subject:  request.GetSubject(),
		Fields:   request.GetFieldMask().GetPaths(),
	}
}
