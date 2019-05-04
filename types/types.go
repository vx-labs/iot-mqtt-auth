package types

//go:generate protoc --go_out=plugins=grpc:. types.proto

func (c *TransportContext) Ensure(a ...TransportAssertion) (bool, string) {
	var tenant string
	for _, assertion := range a {
		success, t := assertion(c)
		if !success {
			return false, ""
		}
		tenant = t
	}
	return true, tenant
}

func (c *ProtocolContext) Ensure(a ...ProtocolAssertion) (bool, string) {
	var tenant string
	for _, assertion := range a {
		success, t := assertion(c)
		if !success {
			return false, ""
		}
		tenant = t
	}
	return true, tenant
}
