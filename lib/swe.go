package lib

import (
	"context"
	"fmt"
	"jenujari/go-sbc-webapp/config"

	"github.com/jenujari/go-swe-api/client"
	"github.com/jenujari/go-swe-api/proto"
)

func GetPing() (*proto.PingResponse, error) {

	cfg := config.GetConfig()

	c, err := client.NewEphServiceClient(cfg.SweGrpcConfig.Addr)
	if err != nil {
		return nil, fmt.Errorf("could not create client: %v", err)
	}
	defer c.Close()

	pingResp, err := c.Ping(context.Background())

	if err != nil {
		return nil, fmt.Errorf("could not ping server: %v", err)
	}

	return pingResp, nil
}
