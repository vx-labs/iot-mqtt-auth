package api

import (
	"github.com/vx-labs/iot-mqtt-auth/types"
	"crypto/x509"
)

type authOptions struct {
	ProtocolContext  *types.ProtocolContext
	TransportContext *types.TransportContext
}

type AuthOpt func(o *authOptions)

func WithTransportContext(encrypted bool, remoteAddr string, cert *x509.Certificate) AuthOpt {
	return func(o *authOptions) {
		o.TransportContext.RemoteAddress = remoteAddr
		o.TransportContext.Encrypted = encrypted
		if cert != nil {
			o.TransportContext.X509Certificate = cert.Raw
		}
	}
}
func WithProtocolContext(username string, password string) AuthOpt {
	return func(o *authOptions) {
		o.ProtocolContext.Username = username
		o.ProtocolContext.Password = password
	}
}

func getOpts(opts []AuthOpt) *authOptions {
	o := &authOptions{
		ProtocolContext:  &types.ProtocolContext{},
		TransportContext: &types.TransportContext{},
	}
	for _, f := range opts {
		f(o)
	}
	return o
}
