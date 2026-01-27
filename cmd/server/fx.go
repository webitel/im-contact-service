package server

import (
	"github.com/webitel/im-contact-service/infra/pubsub"
	"github.com/webitel/im-contact-service/infra/tls"
	"github.com/webitel/im-contact-service/internal/service"
	"github.com/webitel/im-contact-service/internal/store/postgres"
	"github.com/webitel/webitel-go-kit/infra/discovery"
	"go.uber.org/fx"

	"github.com/webitel/im-contact-service/config"
	grpcsrv "github.com/webitel/im-contact-service/infra/server/grpc"
	grpchandler "github.com/webitel/im-contact-service/internal/handler/grpc"
)

func NewApp(cfg *config.Config) *fx.App {
	return fx.New(MainModule(cfg))
}

func MainModule(cfg *config.Config) fx.Option {
	return fx.Options(
		fx.Provide(
			func() *config.Config { return cfg },
			ProvideLogger,
			ProvideSD,
		),
		fx.Invoke(func(discovery discovery.DiscoveryProvider) error { return nil }),
		tls.Module,
		pubsub.Module,
		postgres.Module,
		service.Module,
		grpcsrv.Module,
		grpchandler.Module,
	)
}
