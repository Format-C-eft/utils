package db

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Client interface {
	Get(ctx context.Context, dest interface{}, serializer Squirrel) error
	Select(ctx context.Context, dest interface{}, serializer Squirrel) error
	Exec(ctx context.Context, sq Squirrel) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

func init() {
	api, err := pgxscan.NewDBScanAPI(
		dbscan.WithAllowUnknownColumns(true),
	)
	if err != nil {
		panic(fmt.Errorf("pgxscan.NewDBScanAPI err: %w", err))
	}

	pgxscan.DefaultAPI, err = pgxscan.NewAPI(api)
	if err != nil {
		panic(fmt.Errorf("pgxscan.NewAPI err: %w", err))
	}
}

type Squirrel interface {
	ToSql() (sql string, args []interface{}, err error)
}
