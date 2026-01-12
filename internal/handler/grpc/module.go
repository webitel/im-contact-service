package grpc

import (
	"go.uber.org/fx"

	impb "github.com/webitel/im-contact-service/gen/go/api/v1"
	grpcsrv "github.com/webitel/im-contact-service/infra/server/grpc"
)

var Module = fx.Module("grpc",
	fx.Provide(
		NewContactService,
		NewBotsService,
	),
	fx.Invoke(
		RegisterContactService,
		RegisterBotService,
	),
)

func RegisterContactService(server *grpcsrv.Server, service *ContactService, lc fx.Lifecycle) error {
	impb.RegisterContactsServer(server.Server, service)
	return nil
}

func RegisterBotService(server *grpcsrv.Server, service *BotsServer, lc fx.Lifecycle) error {
	impb.RegisterBotsServer(server.Server, service)
	return nil
}
