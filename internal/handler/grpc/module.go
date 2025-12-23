package grpc

import (
	impb "github.com/webitel/im-contact-service/gen/go/api/v1"
	grpcsrv "github.com/webitel/im-contact-service/infra/server/grpc"
	"go.uber.org/fx"
)

var Module = fx.Module("grpc",
	fx.Provide(
		NewContactService,
	),

	fx.Invoke(
		RegisterContactService,
	),
)

func RegisterContactService(server *grpcsrv.Server, service *ContactService) error {
	impb.RegisterContactsServer(server.Server, service)
	return nil
}
