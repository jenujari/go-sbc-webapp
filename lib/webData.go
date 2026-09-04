package lib

import "jenujari/go-sbc-webapp/config"

type WebData map[string]any

func GetGlobalWebData(cfg *config.Config) WebData {
	return WebData{
		"appname": cfg.WebAppConfig.Appname,
	}
}

func (d WebData) Clone() WebData {
	out := make(WebData, len(d)+8)
	for key, value := range d {
		out[key] = value
	}
	return out
}

func (a *App) PageData() WebData {
	if a == nil {
		return WebData{}
	}
	return a.WebData.Clone()
}
