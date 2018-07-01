package identity

import (
	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	"github.com/vx-labs/iot-mqtt-auth/types"
	"github.com/vx-labs/iot-mqtt-config"
)

type staticVaultProvider struct {
	value string
	api   *vault.Client
}

func NewStaticVaultProvider(api *vault.Client, value string) Provider {
	return &staticVaultProvider{
		value: value,
		api:   api,
	}
}

func (s *staticVaultProvider) CanHandle(app *types.ProtocolContext, _ *types.TransportContext) bool {
	return app.GetUsername() == s.value
}
func (s *staticVaultProvider) Authenticate(app *types.ProtocolContext, _ *types.TransportContext) (*Identity, error) {
	authConfig, err := config.Authentication(s.api)
	if err != nil {
		logrus.Errorf("failed to fetch tokens list from Vaut: %v", err)
		return &Identity{Provider: "static-vault"}, ErrAuthenticationFailed
	}
	pw := app.GetPassword()
	for _, token := range authConfig.StaticTokens {
		if pw == token {
			return &Identity{
				Domain:   "_local",
				Provider: "static-vault",
				ID:       "_anonymous",
				Tenant:   "_default",
			}, nil
		}
	}
	return &Identity{Provider: "static-vault"}, ErrAuthenticationFailed
}
