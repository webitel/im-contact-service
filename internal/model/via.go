package model

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/webitel/webitel-go-kit/pkg/errors"
)

const MaxCharactersInDisableReason int = 255

type ViaCommunication struct {
	ContactID     uuid.UUID      `db:"contact_id" json:"contact_id"`
	Via           string         `db:"via" json:"via"`
	Disable       bool           `db:"disable" json:"disable"`
	DisableReason *string        `db:"disable_reason" json:"disable_reason"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updated_at"`
	Metadata      map[string]any `db:"metadata" json:"metadata"`
}

func (communication *ViaCommunication) CreatedAtUTCUnix() int64 {
	if communication == nil {
		return 0
	}

	return max(communication.CreatedAt.UTC().UnixMilli(), 0)
}

func (communication *ViaCommunication) UpdatedAtUTCUnix() int64 {
	if communication == nil {
		return 0
	}

	return max(communication.UpdatedAt.UTC().UnixMilli(), 0)
}

func (communication *ViaCommunication) Validate() error {
	if communication == nil {
		return errors.InvalidArgument("received nil pointer communication call", errors.WithID("model.communication.validate"))
	}

	if communication.ContactID == uuid.Nil {
		return errors.InvalidArgument("communication contact ID is required", errors.WithID("model.communication.validate"))
	}

	if communication.Via == "" || strings.Trim(communication.Via, " ") == "" {
		return errors.InvalidArgument("communication via is required", errors.WithID("model.communication.validate"))
	}

	if communication.DisableReason != nil && utf8.RuneCountInString(*communication.DisableReason) > MaxCharactersInDisableReason {
		return errors.InvalidArgument(
			fmt.Sprintf("disable reason has more characters than allowed number of %d", MaxCharactersInDisableReason),
		)
	}

	return nil
}

func (communication *ViaCommunication) AvailableFields() []string {
	return []string{"contact_id", "via", "disable", "disable_reason", "created_at", "updated_at", "metadata"}
}

func (communication *ViaCommunication) DefaultFields() []string {
	return []string{"contact_id", "via", "disable", "disable_reason", "created_at", "updated_at", "metadata"}
}

func (communication *ViaCommunication) TableName() string { return "im_contact.via" }

type SearchViaCommunicationsFilter struct {
	Sort       string
	Limit      int
	Page       int
	Fields     []string
	ContactIDs []uuid.UUID
	Vias       []string
	Disabled   *bool
}

func (searchCommunicationsFilter *SearchViaCommunicationsFilter) Validate() error {
	if searchCommunicationsFilter == nil {
		return errors.InvalidArgument("received nil pointer dereference for search communication filter", errors.WithID("model.communication.validate"))
	}

	return nil
}

type CommunicationViaPartialUpdateCmd struct {
	ViaCommunication

	Fields []string
}

func (communicationPartialUpdateCmd *CommunicationViaPartialUpdateCmd) Validate() error {
	if communicationPartialUpdateCmd == nil {
		return errors.InvalidArgument("received nil pointer dereference call for communication partial update cmd", errors.WithID("model.communication.validate"))
	}

	if communicationPartialUpdateCmd.ContactID == uuid.Nil {
		return errors.InvalidArgument("contact id is required in communication partial update cmd", errors.WithID("model.communication.validate"))
	}

	if communicationPartialUpdateCmd.Via == "" {
		return errors.InvalidArgument("via is required in communication partial update cmd", errors.WithID("model.communication.validate"))
	}

	if len(communicationPartialUpdateCmd.Fields) < 1 {
		return errors.InvalidArgument("fields to be updated is required", errors.WithID("model.communication.validate"))
	}

	return nil
}
