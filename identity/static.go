package identity

import "github.com/vx-labs/iot-mqtt-auth/types"

type staticProvider struct {
	value string
	token string
}

func NewStaticProvider(value, token string) Provider {
	return &staticProvider{
		value: value,
		token: token,
	}
}

func (s *staticProvider) CanHandle(app *types.ProtocolContext, _ *types.TransportContext) bool {
	return app.GetUsername() == s.value
}
func (s *staticProvider) Authenticate(app *types.ProtocolContext, _ *types.TransportContext) (*Identity, error) {
	if app.GetPassword() == s.token {
		return &Identity{
			Domain:   "_local",
			Provider: "static",
			ID:       "_anonymous",
			Tenant:   "_default",
		}, nil
	}
	return &Identity{Provider: "static"}, ErrAuthenticationFailed
}
