package server

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/webitel/webitel-go-kit/infra/health"

	grpcsrv "github.com/webitel/im-contact-service/infra/server/grpc"
)

func registerHealth(
	h *health.Registry,
	srv *grpcsrv.Server,
	pool *pgxpool.Pool,
) {
	h.Critical("grpc", health.ListenerCheck(srv.Listener()))
	h.Informational("postgres", pool.Ping)
}
