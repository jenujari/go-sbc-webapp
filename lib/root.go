package lib

import (
	"context"

	"jenujari/go-sbc-webapp/config"
)

var services map[string]any

func init() {
	cfg := config.GetConfig()

	services = make(map[string]any)

	sweClient, _ := NewSweGrpcClient()
	dbService, err := NewDBService(context.Background(), cfg)
	if err != nil {
		config.GetLogger().Println("database initialization failed", err)
	}

	services["sweClient"] = sweClient
	services["planetShadbalaService"] = NewPlanetShadbalaService(sweClient)
	services["webData"] = GetGlobalWebData(cfg)
	services["config"] = cfg
	services["db"] = dbService
}

func GetAllServices() map[string]any {
	return services
}
