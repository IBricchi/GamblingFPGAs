package server

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// password should only be inserted in hashed form!
func (s *SQLiteDB) insertCreds(ctx context.Context, cred credential) error {
	if err := s.TransactContext(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO creds (username, password)
			VALUES (:username, :password)
		`,
			sql.Named("username", cred.username),
			sql.Named("password", cred.password),
		); err != nil {
			return fmt.Errorf("server: sqlite_db_insert: failed to insert credential into db: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("server: sqlite_db_insert: insertCreds transaction failed: %w", err)
	}
	return nil
}

// Inserts the staticTestData.Data field as a comma separated string
func (s *SQLiteDB) insertTestData(ctx context.Context, testData staticTestData) error {
	if err := s.TransactContext(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO tests (info, data)
			VALUES (:info, :data)
		`,
			sql.Named("info", testData.Info),
			sql.Named("data", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(testData.Data)), ","), "[]")),
		); err != nil {
			return fmt.Errorf("server: sqlite_db_insert: failed to insert test data into db: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("server: sqlite_db_insert: insertTestData transaction failed: %w", err)
	}
	return nil
}
