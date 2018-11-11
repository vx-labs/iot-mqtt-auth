package identity

import (
	"crypto/sha1"
	"fmt"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	"github.com/vx-labs/iot-mqtt-auth/types"
	"github.com/vx-labs/iot-mqtt-config"
)

type staticVaultProvider struct {
	value string
	api   *vault.Client
}

func makeSessionID(tenant string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(tenant + time.Now().String()))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
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
			id, err := makeSessionID("_default")
			if err != nil {
				return nil, err
			}
			return &Identity{
				Domain:   "_local",
				Provider: "static-vault",
				ID:       id,
				Tenant:   "_default",
			}, nil
		}
	}
	return &Identity{Provider: "static-vault"}, ErrAuthenticationFailed
}
