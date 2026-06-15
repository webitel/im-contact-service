package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	otelsemconv "go.opentelemetry.io/otel/semconv/v1.38.0"
	"go.uber.org/fx"

	"github.com/webitel/webitel-go-kit/infra/discovery"
	otelsdk "github.com/webitel/webitel-go-kit/infra/otel/sdk"
	"github.com/webitel/webitel-go-kit/infra/profiler"
	"github.com/webitel/webitel-go-kit/pkg/depenlog"
	"github.com/webitel/webitel-go-kit/pkg/errors"
	"github.com/webitel/webitel-go-kit/pkg/logger"
	"github.com/webitel/webitel-go-kit/pkg/semconv"

	"github.com/webitel/im-contact-service/config"
	"github.com/webitel/im-contact-service/infra/db/pg"
	"github.com/webitel/im-contact-service/internal/model"

	_ "github.com/webitel/webitel-go-kit/infra/discovery/consul" // register consul discovery driver
	// -------------------- plugin(s) -------------------- //
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/log/otlp"
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/log/stdout"
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/metric/otlp"
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/metric/stdout"
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/trace/otlp"
	_ "github.com/webitel/webitel-go-kit/infra/otel/sdk/trace/stdout"
)

// ProvideLogger builds the service's unified logger via depenlog and returns
// both the *slog.Logger (slog.SetDefault, for existing consumers) and the
// logger.Logger (for the kit adapters). depenlog.New installs the logger
// process-wide: slog.SetDefault plus grpc-go's global logger (UseGRPC).
//
// When OTel logging is enabled, the OTel SDK is configured (resource, metrics,
// log bridge) and logs are routed through depenlog.WithHandler(otelBridge) so
// the OTel LoggerProvider owns the record schema and trace correlation. The
// bridge is only wired when a log exporter is actually configured, so
// metrics-only OTel still falls back to depenlog's plain console/file handler.
func ProvideLogger(cfg *config.Config, lc fx.Lifecycle) (*slog.Logger, logger.Logger, error) {
	logSettings := cfg.Log

	if !logSettings.Console && !logSettings.Otel && logSettings.File == "" {
		logSettings.Console = true
	}

	depCfg := depenlog.Config{
		Level:   logSettings.Level,
		JSON:    logSettings.JSON,
		File:    logSettings.File,
		Console: logSettings.Console,
	}

	var opts []depenlog.Option

	if logSettings.Otel {
		service := resource.NewSchemaless(
			otelsemconv.ServiceName(model.ServiceName),
			otelsemconv.ServiceVersion(model.Version),
			otelsemconv.ServiceInstanceID(discovery.GenerateInstanceID(model.ServiceName)),
			otelsemconv.ServiceNamespace(model.ServiceNamespace),
		)
		otelBridge := otelslog.NewHandler("slog")

		metricExporter, err := otlpmetricgrpc.New(context.Background())
		if err != nil {
			return nil, nil, fmt.Errorf("create otlp metric exporter: %w", err)
		}

		reader := metric.NewPeriodicReader(metricExporter)

		shutdown, err := otelsdk.Configure(
			context.Background(),
			otelsdk.WithResource(service),
			otelsdk.WithLogBridge(
				func() {
					opts = append(opts, depenlog.WithHandler(otelBridge))
				},
			),
			otelsdk.WithMetricOptions(metric.WithReader(reader)),
		)
		if err != nil {
			return nil, nil, err
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return shutdown(ctx)
			},
		})
	}

	kit, err := depenlog.New(depCfg, opts...)
	if err != nil {
		return nil, nil, err
	}

	return slog.Default(), kit, nil
}

func ProvideSD(cfg *config.Config, log *slog.Logger, lc fx.Lifecycle) (discovery.DiscoveryProvider, error) {
	provider, err := discovery.DefaultFactory.CreateProvider(
		discovery.ProviderConsul,
		log,
		cfg.Consul.Addr,
		discovery.WithHeartbeat[discovery.DiscoveryProvider](true),
		discovery.WithTimeout[discovery.DiscoveryProvider](time.Second*30),
	)
	if err != nil {
		return nil, err
	}

	si := new(discovery.ServiceInstance)
	{
		si.Id = discovery.GenerateInstanceID(model.ServiceName)
		si.Name = model.ServiceName
		si.Version = model.Version
		si.Metadata = map[string]string{
			"commit":         model.Commit,
			"commitDate":     model.CommitDate,
			"branch":         model.Branch,
			"buildTimestamp": model.BuildTimestamp,
		}
		si.Endpoints = []string{(&url.URL{Scheme: "grpc", Host: cfg.Service.Addr}).String()}
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := provider.Register(ctx, si); err != nil {
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := provider.Deregister(ctx, si); err != nil {
				return err
			}

			return nil
		},
	})

	return provider, nil
}

func ProvideNewDBConnection(cfg *config.Config, l *slog.Logger, lc fx.Lifecycle) (*pg.PgxDB, error) {
	config := pg.ConnectionConfig{
		DSN:         cfg.Postgres.DSN,
		OTELEnabled: cfg.Log.Otel,
	}

	db, err := pg.New(context.Background(), l, config)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			db.Master().Close()

			return nil
		},
	})

	return db, err
}

func ProvideProfiler(cfg *config.Config) profiler.Config {
	return profiler.Config{
		Addr:                 cfg.Profiler.Addr,
		MutexProfileFraction: cfg.Profiler.MutexFraction,
		BlockProfileRate:     cfg.Profiler.BlockRate,
	}
}

func ProvideRuntimeMetrics(cfg *config.Config, logger *slog.Logger, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			if !cfg.Log.Otel {
				logger.Info("OTEL disabled, skipping setting up runtime metrics")

				return nil
			}

			logger.Info("starting collecting OpenTelemetry runtime metrics")

			if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second * 5)); err != nil {
				logger.Error("starting collecting OpenTelemetry runtime metrics", semconv.ErrorKey, err)

				return errors.Internal("starting collecting otel runtime metrics", errors.WithCause(err), errors.WithID("server.providers.provide_runtime_metrics"))
			}

			return nil
		},
	})
}
