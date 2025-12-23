package grpc

import (
	"go.uber.org/fx"
)

var Module = fx.Module("grpc",
	fx.Provide(
		fx.Annotate(
			NewContactService,
		),
	),

	fx.Invoke(
		RegisterContactService,
	),
)
