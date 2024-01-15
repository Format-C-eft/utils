package db

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type connAdapterTx struct {
	client pgx.Tx
}

func (a *connAdapterTx) Get(ctx context.Context, dest interface{}, sqlizer Squirrel) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return fmt.Errorf("sq.ToSql err: %w", err)
	}

	return pgxscan.Get(ctx, a.client, dest, query, args...)
}

func (a *connAdapterTx) Select(ctx context.Context, dest interface{}, sq Squirrel) error {
	query, args, err := sq.ToSql()
	if err != nil {
		return fmt.Errorf("sq.ToSql err: %w", err)
	}

	return pgxscan.Select(ctx, a.client, dest, query, args...)
}

func (a *connAdapterTx) Exec(ctx context.Context, sq Squirrel) (pgconn.CommandTag, error) {
	query, args, err := sq.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("sq.ToSql err: %w", err)
	}

	return a.client.Exec(ctx, query, args...)
}

func (a *connAdapterTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return a.client.Begin(ctx)
}
