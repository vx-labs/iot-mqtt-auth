package identity

import (
	"errors"

	"github.com/vx-labs/iot-mqtt-auth/types"
)

var ErrAuthenticationFailed = errors.New("authentication failed")

type Identity struct {
	ID       string
	Tenant   string
	Domain   string
	Provider string
}
type Provider interface {
	CanHandle(*types.ProtocolContext, *types.TransportContext) bool
	Authenticate(*types.ProtocolContext, *types.TransportContext) (*Identity, error)
}
