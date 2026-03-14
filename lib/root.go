package lib

import "jenujari/go-sbc-webapp/config"

var services map[string]any

func init() {
	cfg := config.GetConfig()

	services = make(map[string]any)

	sweClient, _ := NewSweGrpcClient()

	services["sweClient"] = sweClient
	services["webData"] = GetGlobalWebData(cfg)
	services["config"] = cfg
}

func GetAllServices() map[string]any {
	return services
}
