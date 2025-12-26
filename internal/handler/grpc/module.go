package grpc

import (
	"context"

	"go.uber.org/fx"

	impb "github.com/webitel/im-contact-service/gen/go/api/v1"
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
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				impb.RegisterContactsServer(server.Server, service)
				return nil
			},
		})

	return nil
}
