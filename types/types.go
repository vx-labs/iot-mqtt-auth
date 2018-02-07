package types

//go:generate protoc --go_out=plugins=grpc:. types.proto

func (c *TransportContext) Ensure(a ...TransportAssertion) bool {
	for _, assertion := range a {
		if !assertion(c) {
			return false
		}
	}
	return true
}

func (c *ProtocolContext) Ensure(a ...ProtocolAssertion) bool {
	for _, assertion := range a {
		if !assertion(c) {
			return false
		}
	}
	return true
}