package lib

import (
	"context"

	"jenujari/go-sbc-webapp/config"
)

// App is the composition root. Handlers take a typed *App instead of
// an untyped map[string]any service locator.
type App struct {
	Config         *config.Config
	SweClient      SweGrpcClient
	PlanetShadbala PlanetShadbalaService
	WebData        WebData
	DB             *DBService
	OHLC           *OHLCImporter
}

type ctxKey int

const appCtxKey ctxKey = iota

func WithApp(ctx context.Context, app *App) context.Context {
	return context.WithValue(ctx, appCtxKey, app)
}

func AppFromContext(ctx context.Context) (*App, bool) {
	app, ok := ctx.Value(appCtxKey).(*App)
	return app, ok && app != nil
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	sweClient, sweErr := NewSweGrpcClient()
	if sweErr != nil {
		config.GetLogger().Println("swe client initialization failed", sweErr)
	}

	dbService, dbErr := NewDBService(ctx, cfg)
	if dbErr != nil {
		config.GetLogger().Println("database initialization failed", dbErr)
	}

	return &App{
		Config:         cfg,
		SweClient:      sweClient,
		PlanetShadbala: NewPlanetShadbalaService(sweClient),
		WebData:        GetGlobalWebData(cfg),
		DB:             dbService,
		OHLC:           NewOHLCImporter(dbService),
	}, firstErr(sweErr, dbErr)
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
