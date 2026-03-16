package grpc

import (
	"context"
	"log/slog"

	impb "github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/handler/grpc/mapper"
	"github.com/webitel/im-contact-service/internal/handler/grpc/mapper/generated"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/service"
)



type ContactPrivacyServer struct {
	impb.UnimplementedContactPrivacyServer

	logger  *slog.Logger
	handler service.ContactPrivacyService
	inMapper mapper.PrivacyInMapper
}

type PrivacyService interface {
	CanSend(context.Context, *model.CanSendRequest) error
	CanInvite(context.Context, *model.CanInviteRequest) error
}




func NewPrivacyServer(handler service.ContactPrivacyService, logger *slog.Logger) *ContactPrivacyServer{
	return &ContactPrivacyServer{handler: handler, logger: logger, inMapper: &generated.PrivacyInMapperImpl{}}
}

func (c *ContactPrivacyServer) CanSend(ctx context.Context, request *impb.CanSendRequest) (*impb.CanSendResponse, error) {
	converted, err := c.inMapper.ConvertCanSendRequest(request)
	if err != nil {
		return nil, err
	}

	err = c.handler.CanSend(ctx, converted)
	if err != nil {
		return nil, err
	}

	return &impb.CanSendResponse{Can: true}, nil
	
}
func (c *ContactPrivacyServer) CanInvite(ctx context.Context, request *impb.CanInviteRequest) (*impb.CanInviteResponse, error) {
	converted, err := c.inMapper.ConvertCanInviteRequest(request)
	if err != nil {
		return nil, err
	}

	err = c.handler.CanInvite(ctx, converted)
	if err != nil {
		return nil, err
	}

	return &impb.CanInviteResponse{Can: true}, nil
	
}

