package lib

import (
	"context"
	"fmt"
	"jenujari/go-sbc-webapp/config"

	"github.com/jenujari/go-swe-api/client"
	"github.com/jenujari/go-swe-api/proto"
)

type SweGrpcClient interface {
	Ping(ctx context.Context) (*proto.PingResponse, error)
	GetPos(ctx context.Context, datetime string, planet string) (*proto.PosResponse, error)
	Tithy(ctx context.Context, timestamp string) (*proto.TithyResponse, error)
}

type SweGrpcClientImpl struct {
	client *client.EphServiceClient
}

func NewSweGrpcClient() (SweGrpcClient, error) {
	cfg := config.GetConfig()
	c, err := client.NewEphServiceClient(cfg.SweGrpcConfig.Addr)
	if err != nil {
		return nil, fmt.Errorf("could not create client: %v", err)
	}

	sweImpl := new(SweGrpcClientImpl)

	sweImpl.client = c

	return sweImpl, nil
}

func (c *SweGrpcClientImpl) Ping(ctx context.Context) (*proto.PingResponse, error) {
	return c.client.Ping(ctx)
}

func (c *SweGrpcClientImpl) GetPos(ctx context.Context, datetime string, planet string) (*proto.PosResponse, error) {
	return c.client.GetPos(ctx, datetime, planet)
}

func (c *SweGrpcClientImpl) Tithy(ctx context.Context, timestamp string) (*proto.TithyResponse, error) {
	return c.client.Tithy(ctx, timestamp)
}
