package dto

import "github.com/google/uuid"

type (
	SearchBotRequest struct {
		BaseFilter
		
		FlowIds []int64
		DisplayNames []string
		Ids uuid.UUIDs
	}

	UpdateBotCommand struct {
		Id uuid.UUID
		DomainId int
		
		FlowId *int
		DisplayName *string
	}

	DeleteBotCommand struct {
		Id *uuid.UUID
		DomainId int
		FlowId *int
	}

	EnsureBotRequest struct {
		DomainId int `json:"domain_id"`
		FlowId   int `json:"flow_id"`
	}
)

func NewUpdateBotCommand(id uuid.UUID, domainId, flowId int, displayName string) *UpdateBotCommand {
	var cmd = new(UpdateBotCommand)
	{
		cmd.Id = id
		cmd.DomainId = domainId
	}

	if flowId <= 0 {
		cmd.FlowId = nil
	} else {
		cmd.FlowId = &flowId
	}

	if displayName == "" {
		cmd.DisplayName = nil
	} else {
		cmd.DisplayName = &displayName
	}

	return cmd
}