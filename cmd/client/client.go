package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/vx-labs/iot-mqtt-auth/api"
)

func main() {
	ctx := context.Background()
	c, err := api.New("localhost:7994")
	if err != nil {
		logrus.Fatal(err)
	}
	identity, err := c.Authenticate(ctx, api.WithProtocolContext("test", "test"), api.WithTransportContext(true, "127.0.0.1", nil))
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("authentication successful, user's tenant is '%v'", identity.Tenant)

}
