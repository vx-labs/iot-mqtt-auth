package types

type TransportAssertion func(c *TransportContext) bool
type ProtocolAssertion func(c *ProtocolContext) bool

func (p ProtocolAssertion) Or(assertion ProtocolAssertion) ProtocolAssertion {
	return func(c *ProtocolContext) bool {
		return p(c) || assertion(c)
	}
}

func (p TransportAssertion) Or(assertion TransportAssertion) TransportAssertion {
	return func(c *TransportContext) bool {
		return p(c) || assertion(c)
	}
}

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
