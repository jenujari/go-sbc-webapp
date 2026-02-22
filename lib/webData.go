package lib

import "jenujari/go-sbc-webapp/config"

type WebData map[string]any

func GetGlobalWebData(cfg *config.Config) WebData {
	return WebData{
		"appname": cfg.WebAppConfig.Appname,
	}
}
