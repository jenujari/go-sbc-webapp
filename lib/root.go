package lib

import (
	"context"

	"jenujari/go-sbc-webapp/config"
)

var app *App

func init() {
	a, err := NewApp(context.Background(), config.GetConfig())
	if err != nil {
		config.GetLogger().Println("app initialization incomplete", err)
	}
	app = a
}

func GetApp() *App {
	return app
}
