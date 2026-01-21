package postgres

import (
	"github.com/webitel/im-contact-service/internal/store"
	"go.uber.org/fx"
)

var Module = fx.Module("store",
	fx.Provide(

		fx.Annotate(
			NewContactStore,
			fx.As(new(store.ContactStore))),
		),
)
