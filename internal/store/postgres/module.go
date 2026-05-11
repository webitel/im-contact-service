package postgres

import (
	"go.uber.org/fx"

	"github.com/webitel/im-contact-service/internal/store"
)

var Module = fx.Module("store",
	fx.Provide(

		fx.Annotate(
			NewContactStore,
			fx.As(new(store.ContactStore))),

		fx.Annotate(
			NewSettingsStore,
			fx.As(new(store.SettingsStore)),
		),
		fx.Annotate(newViaStore, fx.As(new(store.ViaStore))),
	))
