package cmd

import (
	"go.uber.org/fx"

	"github.com/webitel/im-contact-service/config"
	grpchandler "github.com/webitel/im-contact-service/internal/handler/grpc"
)

func NewApp(cfg *config.Config) *fx.App {
	return fx.New(
		fx.Provide(
			func() *config.Config { return cfg },
			ProvideLogger,
			ProvideGrpcServer,
			ProvideSD,
			ProvidePubSub,
			ProvideNewDBConnection,
		),
		fx.Invoke(
			StartGrpcServer,
		),
		grpchandler.Module,
	)
}
