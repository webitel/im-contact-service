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
		NewPrivacyServer),
	fx.Invoke(
		RegisterContactService,
		RegisterContactSettingsService,
		RegisterContactPrivacyService,
	),
)

func RegisterContactService(server *grpcsrv.Server, service *ContactServer, lc fx.Lifecycle) error {
	impb.RegisterContactsServer(server.Server, service)
	return nil
}
func RegisterContactSettingsService(server *grpcsrv.Server, service *ContactSettingsServer, lc fx.Lifecycle) error {
	impb.RegisterContactSettingsServer(server.Server, service)
	return nil
}


func RegisterContactPrivacyService(server *grpcsrv.Server, service *ContactPrivacyServer, lc fx.Lifecycle) error {
	impb.RegisterContactPrivacyServer(server.Server, service)
	return nil
}