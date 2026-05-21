package pg

import (
	"context"
	"log/slog"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/webitel/webitel-go-kit/pkg/errors"
)

type PgxDB struct {
	master *pgxpool.Pool
	logger *slog.Logger
}

type ConnectionConfig struct {
	DSN         string
	OTELEnabled bool
}

func New(ctx context.Context, logger *slog.Logger, config ConnectionConfig) (*PgxDB, error) {
	cfg, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, errors.InvalidArgument("parsing DSN", errors.WithCause(err), errors.WithID("pg.pg.new"))
	}

	if config.OTELEnabled {
		cfg.ConnConfig.Tracer = otelpgx.NewTracer(
			otelpgx.WithTrimSQLInSpanName(),
		)
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errors.Internal("creating DB connection pool", errors.WithCause(err), errors.WithID("pg.pg.new"))
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, errors.Internal("pinging database", errors.WithCause(err), errors.WithID("pg.pg.new"))
	}

	if config.OTELEnabled {
		err = otelpgx.RecordStats(
			dbpool,
			otelpgx.WithMinimumReadDBStatsInterval(10*time.Second),
		)
		if err != nil {
			logger.Error("[DB] starting recording pool statistic", "error", err)

			return nil, errors.Internal("starting recording pool statistic", errors.WithCause(err), errors.WithID("pg.pg.new"))
		}
	}

	db := new(PgxDB)
	{
		db.logger = logger
		db.master = dbpool
	}

	return db, nil
}

func (d *PgxDB) Master() *pgxpool.Pool {
	return d.master
}
