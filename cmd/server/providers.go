package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
	"go.uber.org/fx"

	"github.com/webitel/webitel-go-kit/infra/discovery"
	otelsdk "github.com/webitel/webitel-go-kit/infra/otel/sdk"
	"github.com/webitel/webitel-go-kit/infra/profiler"
	"github.com/webitel/webitel-go-kit/pkg/errors"
	"github.com/webitel/webitel-go-kit/pkg/logger"

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

func ProvideLogger(cfg *config.Config, lc fx.Lifecycle) (*slog.Logger, error) {
	logSettings := cfg.Log

	if !logSettings.Console && !logSettings.Otel && logSettings.File == "" {
		logSettings.Console = true
	}

	level := parseLevel(logSettings.Level)
	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handlers []slog.Handler

	if logSettings.Console {
		var h slog.Handler
		if logSettings.JSON {
			h = slog.NewJSONHandler(os.Stdout, opts)
		} else {
			h = slog.NewTextHandler(os.Stdout, opts)
		}

		handlers = append(handlers, h)
	}

	// File Handler
	if logSettings.File != "" {
		f, err := os.OpenFile(logSettings.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}

		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return f.Close()
			},
		})

		var h slog.Handler
		if logSettings.JSON {
			h = slog.NewJSONHandler(f, opts)
		} else {
			h = slog.NewTextHandler(f, opts)
		}

		handlers = append(handlers, h)
	}

	if logSettings.Otel {
		service := resource.NewSchemaless(
			semconv.ServiceName(model.ServiceName),
			semconv.ServiceVersion(model.Version),
			semconv.ServiceInstanceID(cfg.Service.ID),
			semconv.ServiceNamespace(model.ServiceNamespace),
		)
		otelHandler := otelslog.NewHandler("slog")

		metricExporter, err := otlpmetricgrpc.New(context.Background())
		if err != nil {
			return nil, fmt.Errorf("create otlp metric exporter: %w", err)
		}

		reader := metric.NewPeriodicReader(metricExporter)

		shutdown, err := otelsdk.Configure(
			context.Background(),
			otelsdk.WithResource(service),
			otelsdk.WithLogBridge(
				func() {
					handlers = append(handlers, otelHandler)
				},
			),
			otelsdk.WithMetricOptions(metric.WithReader(reader)),
		)
		if err != nil {
			return nil, err
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return shutdown(ctx)
			},
		})
	}

	var finalHandler slog.Handler

	switch len(handlers) {
	case 0:
		finalHandler = slog.NewTextHandler(os.Stdout, opts)
	case 1:
		finalHandler = handlers[0]
	default:
		finalHandler = MultiHandler(handlers...)
	}

	logger := slog.New(finalHandler)
	slog.SetDefault(logger)

	return logger, nil
}

func parseLevel(lvl string) slog.Level {
	switch lvl {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type multiHandler struct {
	handlers []slog.Handler
}

func MultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, hh := range h.handlers {
		if hh.Enabled(ctx, level) {
			return true
		}
	}

	return false
}

func (h *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, hh := range h.handlers {
		if hh.Enabled(ctx, r.Level) {
			_ = hh.Handle(ctx, r)
		}
	}

	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, hh := range h.handlers {
		newHandlers[i] = hh.WithAttrs(attrs)
	}

	return &multiHandler{handlers: newHandlers}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, hh := range h.handlers {
		newHandlers[i] = hh.WithGroup(name)
	}

	return &multiHandler{handlers: newHandlers}
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
		si.Id = cfg.Service.ID
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

func ProvideProfiler(cfg *config.Config, log *slog.Logger) (profiler.Config, logger.Logger) {
	return profiler.Config{
		Addr:                 cfg.Profiler.Addr,
		MutexProfileFraction: cfg.Profiler.MutexFraction,
		BlockProfileRate:     cfg.Profiler.BlockRate,
	}, logger.NewSlog(log)
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
				logger.Error("starting collecting OpenTelemetry runtime metrics", "error", err)

				return errors.Internal("starting collecting otel runtime metrics", errors.WithCause(err), errors.WithID("server.providers.provide_runtime_metrics"))
			}

			return nil
		},
	})
}
