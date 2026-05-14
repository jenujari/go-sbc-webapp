package lib

import (
	"context"
	"fmt"
	"net/url"

	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/sqls"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBService struct {
	Pool    *pgxpool.Pool
	Queries *sqls.Queries
}

func NewDBService(ctx context.Context, cfg *config.Config) (*DBService, error) {
	dsn, err := buildPostgresDSN(cfg.DBConfig)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &DBService{Pool: pool, Queries: sqls.New(pool)}, nil
}

func buildPostgresDSN(db config.DBConfig) (string, error) {
	if db.Host == "" || db.Port == "" || db.User == "" || db.Password == "" || db.Name == "" {
		return "", fmt.Errorf("database configuration is incomplete")
	}
	if db.SSLMode == "" {
		db.SSLMode = "disable"
	}

	dsn := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(db.User, db.Password),
		Host:   fmt.Sprintf("%s:%s", db.Host, db.Port),
		Path:   db.Name,
	}
	query := dsn.Query()
	query.Set("sslmode", db.SSLMode)
	dsn.RawQuery = query.Encode()

	return dsn.String(), nil
}
