package cmd

import (
	"context"
	"fmt"
	"log/slog"

	grpc_srv "github.com/webitel/im-contact-service/infra/server/grpc"
	"go.uber.org/fx"
)

func StartGrpcServer(lc fx.Lifecycle, srv *grpc_srv.Server, log *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info(fmt.Sprintf("listen grpc %s:%d", srv.Host(), srv.Port()))
				if err := srv.Listen(); err != nil {
					log.Error("grpc server error", "err", err)
				}
			}()
			return nil
		},
	})
}
