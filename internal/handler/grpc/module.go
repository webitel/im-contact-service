package grpc

import (
	impb "github.com/webitel/im-contact-service/gen/go/api/v1"
	grpc_srv "github.com/webitel/im-contact-service/infra/server/grpc"
	"go.uber.org/fx"
)

var Module = fx.Module("grpc",
	fx.Provide(
		fx.Annotate(
			NewContactService,
		),
	),

	fx.Invoke(
		RegisterContactService,
	),
)

func RegisterContactService(server *grpc_srv.Server, service *ContactService) error {
	impb.RegisterContactsServer(server.Server, service)
	return nil
}
