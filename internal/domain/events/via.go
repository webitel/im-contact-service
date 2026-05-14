package events

import (
	"github.com/webitel/im-contact-service/internal/model"
)

const (
	ViaCreatedTopic string = "contact.via.created."
	ViaUpdatedTopic string = "contact.via.updated."
)

type ViaCreated struct {
	Base `json:",inline"`
}

func NewViaCreatedEvent(via *model.ViaCommunication) *ViaCreated {
	if via == nil {
		return nil
	}

	return &ViaCreated{
		Base: Base{
			ID:        via.ContactID,
			TopicName: ViaCreatedTopic + via.ContactID.String() + "." + via.Via,
			Timestamp: via.CreatedAt,
		},
	}
}

type ViaUpdated struct {
	Base `json:",inline"`

	Via           string  `json:"via"`
	Disable       bool    `json:"disable"`
	DisableReason *string `json:"disable_reason,omitempty"`
}

func NewViaUpdatedEvent(via *model.ViaCommunication) *ViaUpdated {
	if via == nil {
		return nil
	}

	return &ViaUpdated{
		Base: Base{
			ID:        via.ContactID,
			TopicName: ViaUpdatedTopic + via.ContactID.String() + "." + via.Via,
			Timestamp: via.UpdatedAt,
		},
		Via:           via.Via,
		Disable:       via.Disable,
		DisableReason: via.DisableReason,
	}
}
