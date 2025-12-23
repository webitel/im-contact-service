package cmd

import (
	"github.com/webitel/im-contact-service/config"
	"go.uber.org/fx"
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
	)
}
