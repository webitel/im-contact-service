package server

import (
	"go.uber.org/fx"

	"github.com/webitel/webitel-go-kit/infra/discovery"
	"github.com/webitel/webitel-go-kit/infra/profiler"

	"github.com/webitel/im-contact-service/config"
	"github.com/webitel/im-contact-service/infra/pubsub"
	grpcsrv "github.com/webitel/im-contact-service/infra/server/grpc"
	"github.com/webitel/im-contact-service/infra/tls"
	grpchandler "github.com/webitel/im-contact-service/internal/handler/grpc"
	"github.com/webitel/im-contact-service/internal/service"
	"github.com/webitel/im-contact-service/internal/store/postgres"
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
			ProvideNewDBConnection,
			ProvideProfiler,
		),
		fx.Invoke(func(_ discovery.DiscoveryProvider) error { return nil }),
		tls.Module,
		fx.Invoke(ProvideRuntimeMetrics),
		pubsub.Module,
		postgres.Module,
		service.Module,
		grpcsrv.Module,
		grpchandler.Module,
		profiler.Module,
	)
}
