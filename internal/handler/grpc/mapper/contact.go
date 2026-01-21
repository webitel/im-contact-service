package mapper

import (
	"github.com/google/uuid"
	impb "github.com/webitel/im-contact-service/gen/go/api/contact/v1"
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
	var (
		from = MapOneof2ConstPeerKind(request.From)
		to = MapOneof2ConstPeerKind(request.To)
	)

	canSendQuery := &dto.CanSendQuery{
		DomainId: int(request.GetDomainId()),
		From: from,
		To: to,
	}

	return canSendQuery
}

func MapOneof2ConstPeerKind(pbPeer *impb.CanSendRequest_Peer) model.Peer {
	if pbPeer == nil {
		return model.Peer{}
	}

	var peer model.Peer
	switch kind := pbPeer.Kind.(type) {
	case *impb.CanSendRequest_Peer_BotId:
		{
			peer.Id, _ = uuid.Parse(kind.BotId)
			peer.Kind = model.PeerBot 
		}
	case *impb.CanSendRequest_Peer_ContactId:
		{
			peer.Id, _ = uuid.Parse(kind.ContactId)
			peer.Kind = model.PeerContact
		}
	}

	return peer
}