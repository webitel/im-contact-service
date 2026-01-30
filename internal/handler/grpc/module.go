package grpc

import (
	"go.uber.org/fx"

	impb "github.com/webitel/im-contact-service/gen/go/contact/v1"
	grpcsrv "github.com/webitel/im-contact-service/infra/server/grpc"
)

var Module = fx.Module("grpc",
	fx.Provide(
		NewContactService,
	),
	fx.Invoke(
		RegisterContactService,		
	),
)

func RegisterContactService(server *grpcsrv.Server, service *ContactService, lc fx.Lifecycle) error {
	impb.RegisterContactsServer(server.Server, service)
	return nil
}
