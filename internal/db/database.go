package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func newDB(pool *pgxpool.Pool) *Database {
	return &Database{
		pool: pool,
	}
}

func (db *Database) GetPool(_ context.Context) *pgxpool.Pool {
	return db.pool
}

func (db *Database) Get(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Get(ctx, db.pool, dest, query, args...)
}

func (db *Database) Select(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, db.pool, dest, query, args...)
}

func (db *Database) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, query, args...)
}

func (db *Database) ExecQueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return db.pool.QueryRow(ctx, query, args...)
}
