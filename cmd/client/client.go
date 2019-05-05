package main

import (
	"github.com/vx-labs/iot-mqtt-auth/api"
	"github.com/sirupsen/logrus"
	"context"
)

func main() {
	ctx := context.Background()
	c, err := api.New("localhost:7994")
	if err != nil {
		logrus.Fatal(err)
	}
	status,  tenant, err := c.Authenticate(ctx, api.WithProtocolContext("test", "test"), api.WithTransportContext(true, "127.0.0.1", nil))
	if err != nil {
		logrus.Fatal(err)
	}
	if status {
		logrus.Infof("authentication successful, user's tenant is '%v'", tenant)
	} else {
		logrus.Infof("authentication failed")
	}
}
