package server

import (
	"context"
	"database/sql"
)

type DB interface {
	// insert
	insertCreds(ctx context.Context, cred credential) error
	insertTestData(ctx context.Context, testData staticTestData) error

	// retrieve
	getCreds(ctx context.Context) (map[string]string, error)

	// general
	migrate(ctx context.Context) error
	TransactContext(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) error) error
	Close() error
}
