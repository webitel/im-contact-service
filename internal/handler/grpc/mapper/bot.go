package mapper

import (
	"github.com/google/uuid"
	impb "github.com/webitel/im-contact-service/gen/go/api/v1"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
)

func SearchBotsRequest2BotsFilter(request *impb.SearchBotsRequest) *dto.SearchBotRequest {
	var ids uuid.UUIDs = make(uuid.UUIDs, 0, len(request.GetIds()))
	if len(request.GetIds()) > 0 {
		for _, id := range request.Ids {
			if parsedId, err := uuid.Parse(id); err == nil {
				ids = append(ids, parsedId)
			}
		}
	}
	
	return &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: int(request.GetDomainId()),
			Page: int(request.GetPage()),
			Size: int(request.GetSize()),
			Sort: request.GetSort(),
			Q: request.GetQ(),
		},
		FlowIds: request.GetFlowIds(),
		DisplayNames: request.GetDisplayNames(),
		Ids: ids,
	}
}

func BotModelList2SearchResponse(bots []*model.WebitelBot, hasNext bool) *impb.SearchBotsResponse {
	if bots == nil {
		return nil
	}
	
	webitelBotsList := make([]*impb.WebitelBot, 0, len(bots))
	
	for _, bot := range bots {
		webitelBot := &impb.WebitelBot{
			Id: bot.Id.String(),
			DomainId: int32(bot.DomainId),
			CreatedAt: bot.CreatedAt.UnixMilli(),
			UpdatedAt: bot.UpdatedAt.UnixMilli(),
			FlowId: int64(bot.FlowId),
			DisplayName: bot.DisplayName,
		}

		webitelBotsList = append(webitelBotsList, webitelBot)
	}

	return &impb.SearchBotsResponse{
		Next: hasNext,
		WebitelBotList: webitelBotsList,
	}
}

func BotsCreateRequest2CreateModel(req *impb.CreateBotsRequest) *model.WebitelBot {
	return &model.WebitelBot{
		BaseModel: model.BaseModel{
			DomainId: int(req.GetDomainId()),
		},
		FlowId: int(req.GetDomainId()),
		DisplayName: req.GetDisplayName(),
	}
}

func BotsModel2CreateResponse(bot *model.WebitelBot) *impb.CreateBotsResponse {
	if bot == nil {
		return nil
	}
	
	return &impb.CreateBotsResponse{
		WebitelBot: &impb.WebitelBot{
			Id: bot.Id.String(),
			DomainId: int32(bot.DomainId),
			CreatedAt: bot.CreatedAt.UnixMilli(),
			UpdatedAt: bot.UpdatedAt.UnixMilli(),
			FlowId: int64(bot.FlowId),
			DisplayName: bot.DisplayName,
		},
	}
}


func BotsUpdateRequest2UpdateCommand(req *impb.UpdateBotsRequest) *dto.UpdateBotCommand {
	var (
		id, _ = uuid.Parse(req.Id)
		flowId = int(req.GetFlowId())
	)
	
	return &dto.UpdateBotCommand{
		Id: id,
		DomainId: int(req.GetDomainId()),
		FlowId: &flowId,
		DisplayName: req.DisplayName,
	}
}

func BotsModel2UpdateResponse(bot *model.WebitelBot) *impb.UpdateBotsResponse {
	if bot == nil {
		return nil
	}

	return &impb.UpdateBotsResponse{
		WebitelBot: &impb.WebitelBot{
			Id: bot.Id.String(),
			DomainId: int32(bot.DomainId),
			CreatedAt: bot.CreatedAt.UnixMilli(),
			UpdatedAt: bot.UpdatedAt.UnixMilli(),
			FlowId: int64(bot.FlowId),
			DisplayName: bot.DisplayName,
		},
	}
}

func BotsDeleteRequest2DeleteCommand(req *impb.DeleteBotsRequest) *dto.DeleteBotCommand {
	var (
		id *uuid.UUID
		flowId *int
	)

	if req.Id != nil {
		idPtr, _ := uuid.Parse(*req.Id)
		id = &idPtr
	}

	if req.FlowId != nil {
		flowPtr := int(*req.FlowId)
		flowId = &flowPtr
	}
	
	return &dto.DeleteBotCommand{
		DomainId: int(req.GetDomainId()),
		Id: id,
		FlowId: flowId,
	}
}

func BotsEnsureRequest2EnsureDTO(req *impb.EnsureBotRequest) *dto.EnsureBotRequest {
	return &dto.EnsureBotRequest{
		DomainId: int(req.GetDomainId()),
		FlowId: int(req.GetFlowId()),
	}
}

func BotModel2EnsureResponse(bot *model.WebitelBot) *impb.EnsureBotResponse {
	if bot == nil {
		return nil
	}
	
	return &impb.EnsureBotResponse{
		WebitelBot: &impb.WebitelBot{
			Id: bot.Id.String(),
			DomainId: int32(bot.DomainId),
			CreatedAt: bot.CreatedAt.UnixMilli(),
			UpdatedAt: bot.UpdatedAt.UnixMilli(),
			FlowId: int64(bot.FlowId),
			DisplayName: bot.DisplayName,
		},
	}
}