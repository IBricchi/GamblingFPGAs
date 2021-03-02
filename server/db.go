package server

import (
	"context"
	"database/sql"
)

type DB interface {
	// general
	migrate(ctx context.Context) error
	TransactContext(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) error) error
	Close() error
}
