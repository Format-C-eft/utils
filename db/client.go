package db

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type connAdapter struct {
	client *pgxpool.Pool
}

func NewClient(ctx context.Context, cfg Config) (Client, error) {
	config, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("pgx.ParseConfig err: %w", err)
	}

	config.ConnConfig.ConnectTimeout = cfg.ConnectTimeout
	config.MinConns = int32(cfg.MinConn)
	config.MaxConns = int32(cfg.MaxConn)
	config.HealthCheckPeriod = cfg.HealCheckDuration

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("pgx.Connect err: %w", err)
	}

	return &connAdapter{
		client: pool,
	}, nil
}

func (a *connAdapter) Get(ctx context.Context, dest interface{}, sqlizer Squirrel) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return fmt.Errorf("sq.ToSql err: %w", err)
	}

	return pgxscan.Get(ctx, a.client, dest, query, args...)
}

func (a *connAdapter) Select(ctx context.Context, dest interface{}, sq Squirrel) error {
	query, args, err := sq.ToSql()
	if err != nil {
		return fmt.Errorf("sq.ToSql err: %w", err)
	}

	return pgxscan.Select(ctx, a.client, dest, query, args...)
}

func (a *connAdapter) Exec(ctx context.Context, sq Squirrel) (pgconn.CommandTag, error) {
	query, args, err := sq.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("sq.ToSql err: %w", err)
	}

	return a.client.Exec(ctx, query, args...)
}

func (a *connAdapter) Begin(ctx context.Context) (pgx.Tx, error) {
	return a.client.Begin(ctx)
}
