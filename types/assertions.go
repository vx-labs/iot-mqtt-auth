package types

type TransportAssertion func(c *TransportContext) bool
type ProtocolAssertion func(c *ProtocolContext) bool

func MustBeEncrypted() TransportAssertion {
	return func(c *TransportContext) bool {
		return c.Encrypted
	}
}

func MustUseStaticSharedKey(key string) ProtocolAssertion {
	return func(c *ProtocolContext) bool {
		if c.Username == "vx:psk" {
			return key == c.Password
		}
		return false
	}
}
