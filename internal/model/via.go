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

type CreateViaCommunicationCommand struct {
	ContactID     uuid.UUID
	Iss           string
	Sub           string
	Disable       bool
	DisableReason *string
	Metadata      map[string]any
	Via           string
}

func (c *CreateViaCommunicationCommand) Validate() error {
	if c == nil {
		return errors.InvalidArgument("received nil pointer call for create contact request", errors.WithID("model.via.validate"))
	}

	if c.ContactID == uuid.Nil && (c.Iss == "" || c.Sub == "") {
		return errors.InvalidArgument("contact id or pair (iss + sub) is required", errors.WithID("model.via.validate"))
	}

	if strings.Trim(c.Via, " ") == "" {
		return errors.InvalidArgument("via is required", errors.WithID("model.via.validate"))
	}

	if c.ContactID != uuid.Nil {
		return nil
	}

	if c.Iss == "" || c.Sub == "" {
		return errors.InvalidArgument("both values in pair (iss + sub) are required", errors.WithID("model.via.validate"))
	}

	return nil
}

func (c *CreateViaCommunicationCommand) GetContactIDPtr() *uuid.UUID {
	if c.ContactID == uuid.Nil {
		return nil
	}

	return &c.ContactID
}

func (c *CreateViaCommunicationCommand) GetIssPtr() *string {
	if c.Iss == "" {
		return nil
	}

	return &c.Iss
}

func (c *CreateViaCommunicationCommand) GetSubPtr() *string {
	if c.Sub == "" {
		return nil
	}

	return &c.Sub
}
