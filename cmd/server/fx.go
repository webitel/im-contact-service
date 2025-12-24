package server

import (
	"go.uber.org/fx"

	"github.com/webitel/im-contact-service/config"
	grpcsrv "github.com/webitel/im-contact-service/infra/server/grpc"
	grpchandler "github.com/webitel/im-contact-service/internal/handler/grpc"
)

func NewApp(cfg *config.Config) *fx.App {
	return fx.New(
		fx.Provide(
			func() *config.Config { return cfg },
			ProvideLogger,
			ProvideSD,
			ProvidePubSub,
			ProvideNewDBConnection,
		),
		grpcsrv.Module,
		grpchandler.Module,
	)
}
