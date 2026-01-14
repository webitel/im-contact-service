package grpc

import (
	"context"
	"log/slog"

	impb "github.com/webitel/im-contact-service/gen/go/api/v1"
	"github.com/webitel/im-contact-service/internal/handler/grpc/mapper"
	"github.com/webitel/im-contact-service/internal/service"
)
type BotsServer struct {
	impb.UnimplementedBotsServer

	logger *slog.Logger
	botManager service.BotManager
}

func NewBotsService(logger *slog.Logger, botManager service.BotManager) *BotsServer {
	return &BotsServer{
		logger: logger,
		botManager: botManager,
	}
}

func (s *BotsServer) Search(ctx context.Context,req *impb.SearchBotsRequest) (*impb.SearchBotsResponse, error) {
	filter := mapper.SearchBotsRequest2BotsFilter(req)
	
	botsList, hasNext, err := s.botManager.Search(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := mapper.BotModelList2SearchResponse(botsList, hasNext)

	return response, nil
}
func (s *BotsServer) Create(ctx context.Context, req *impb.CreateBotsRequest) (*impb.CreateBotsResponse, error) {
	model := mapper.BotsCreateRequest2CreateModel(req)
	
	model, err := s.botManager.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	response := mapper.BotsModel2CreateResponse(model)

	return response, nil
}

func (s *BotsServer) Update(ctx context.Context, req *impb.UpdateBotsRequest) (*impb.UpdateBotsResponse, error) {
	updateCmd := mapper.BotsUpdateRequest2UpdateCommand(req)

	bot, err := s.botManager.Update(ctx, updateCmd)
	if err != nil {
		return nil, err
	}

	response := mapper.BotsModel2UpdateResponse(bot)

	return response, nil
}

func (s *BotsServer) Delete(ctx context.Context, req *impb.DeleteBotsRequest) (*impb.DeleteBotsResponse, error) {
	deleteCmd := mapper.BotsDeleteRequest2DeleteCommand(req)

	if err := s.botManager.Delete(ctx, deleteCmd); err != nil {
		return nil, err
	}

	return &impb.DeleteBotsResponse{}, nil
}

func (s *BotsServer) EnsureBot(ctx context.Context, req *impb.EnsureBotRequest) (*impb.EnsureBotResponse, error) {
	ensureBotDto := mapper.BotsEnsureRequest2EnsureDTO(req)

	bot, err := s.botManager.EnsureBot(ctx, ensureBotDto)
	if err != nil {
		return nil, err
	}

	response := mapper.BotModel2EnsureResponse(bot)

	return response, nil
}



