package types

type TransportAssertion func(c *TransportContext) (bool, string)
type ProtocolAssertion func(c *ProtocolContext) (bool, string)

func (p ProtocolAssertion) Or(assertion ProtocolAssertion) ProtocolAssertion {
	return func(c *ProtocolContext) (bool, string) {
		left, tenant := p(c)
		if left {
			return left, tenant
		}
		right, tenant := assertion(c)
		if right {
			return right, tenant
		}
		return false, ""
	}
}

func (p TransportAssertion) Or(assertion TransportAssertion) TransportAssertion {
	return func(c *TransportContext) (bool, string) {
		left, tenant := p(c)
		if left {
			return left, tenant
		}
		right, tenant := assertion(c)
		if right {
			return right, tenant
		}
		return false, ""
	}
}

func AlwaysAllowTransport() TransportAssertion {
	return func(c *TransportContext) (bool, string) {
		return true, ""
	}
}
func MustBeEncrypted() TransportAssertion {
	return func(c *TransportContext) (bool, string) {
		return c.Encrypted, ""
	}
}

func MustUseStaticSharedKey(key string) ProtocolAssertion {
	return func(c *ProtocolContext) (bool, string) {
		if c.Username == "vx:psk" || c.Username == "vx_psk" {
			return key == c.Password, "_default"
		}
		return false, ""
	}
}

func MustUseDemoCredentials() ProtocolAssertion {
	return func(c *ProtocolContext) (bool, string) {
		if c.Username == "demo" && c.Password == "demo" {
			return true, "demo"
		}
		return false, ""
	}
}
