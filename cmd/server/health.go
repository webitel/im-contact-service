package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/webitel/webitel-go-kit/infra/health"

	grpcsrv "github.com/webitel/im-contact-service/infra/server/grpc"
)

func registerHealth(
	lc fx.Lifecycle,
	h *health.Registry,
	srv *grpcsrv.Server,
	pool *pgxpool.Pool,
) {
	lc.Append(fx.StartHook(func() {
		h.Critical("grpc", health.ListenerCheck(srv.Listener()))
		h.Informational("postgres", pool.Ping)
	}))
}
