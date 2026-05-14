package grpc

import (
	"go.uber.org/fx"

	impb "github.com/webitel/im-contact-service/gen/go/contact/v1"
	grpcsrv "github.com/webitel/im-contact-service/infra/server/grpc"
)

var Module = fx.Module("grpc",
	fx.Provide(
		NewContactService,
		NewContactSettingsServer,
		NewPrivacyServer,
		newViaServer,
	),
	fx.Invoke(
		RegisterContactService,
		RegisterContactSettingsService,
		RegisterContactPrivacyService,
		RegisterViaServer,
	),
)

func RegisterContactService(server *grpcsrv.Server, service *ContactServer, _ fx.Lifecycle) error {
	impb.RegisterContactsServer(server.Server, service)

	return nil
}

func RegisterContactSettingsService(server *grpcsrv.Server, service *ContactSettingsServer, _ fx.Lifecycle) error {
	impb.RegisterContactSettingsServer(server.Server, service)

	return nil
}

func RegisterContactPrivacyService(server *grpcsrv.Server, service *ContactPrivacyServer, _ fx.Lifecycle) error {
	impb.RegisterContactPrivacyServer(server.Server, service)

	return nil
}

func RegisterViaServer(server *grpcsrv.Server, srv *ViaServer, _ fx.Lifecycle) error {
	impb.RegisterViasServer(server.Server, srv)

	return nil
}
