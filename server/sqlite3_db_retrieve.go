package server

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *SQLiteDB) getCreds(ctx context.Context) (map[string]string, error) {
	creds := map[string]string{}
	if err := s.TransactContext(ctx, func(ctx context.Context, tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, `
			SELECT username, password
			FROM creds
		`)
		if err != nil {
			return fmt.Errorf("server: sqlite_db_retrieve: failed to retrieve creds rows: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var username, password string
			if err := rows.Scan(
				&username,
				&password,
			); err != nil {
				return fmt.Errorf("server: sqlite_db_retrieve: failed to scan creds row: %w", err)
			}

			creds[username] = password
		}

		if err := rows.Err(); err != nil {
			return fmt.Errorf("server: sqlite_db_retrieve: failed to scan last creds row: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("server: sqlite_db_retrieve: getCreds transaction failed: %w", err)
	}

	return creds, nil
}
