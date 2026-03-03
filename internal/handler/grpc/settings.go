package grpc

import (
	"context"

	impb "github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/handler/grpc/mapper"
	"github.com/webitel/im-contact-service/internal/service/dto"
)

var _ impb.ContactSettingsServer = &ContactSettingsServer{}


type ContactSettingsService interface {
	Get(ctx context.Context, request *dto.GetContactSettingsRequest) (*model.ContactSettings, error)
	Update(ctx context.Context, request *dto.UpdateContactSettingsRequest) (*model.ContactSettings, error)
}

type ContactSettingsServer struct {
	impb.UnimplementedContactSettingsServer

	service ContactSettingsService
	inConverter mapper.SettingsInConverter
	outConverter mapper.SettingsOutConverter
}

// Get implements [contact.ContactSettingsServer].
func (c *ContactSettingsServer) Get(ctx context.Context, req *impb.GetContactSettingsRequest) (*impb.Settings, error) {
	converted, err := c.inConverter.ConvertGetSettingsRequest(req)
	if err != nil {
		return nil, err
	}
	settings, err := c.service.Get(ctx, converted)
	if err != nil {
		return nil, err
	}
	return c.outConverter.ConvertSettings(settings)
}

// Update implements [contact.ContactSettingsServer].
func (c *ContactSettingsServer) Update(ctx context.Context, req *impb.UpdateContactSettingsRequest) (*impb.Settings, error) {
	converted, err := c.inConverter.ConvertUpdateSettingsRequest(req)
	if err != nil {
		return nil, err
	}
	settings, err := c.service.Update(ctx, converted)
	if err != nil {
		return nil, err
	}
	return c.outConverter.ConvertSettings(settings)
	
}

