package publisher

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sdreger/lib-file-processor-go/db/transaction"
	"io"
	"log"
)

type PostgresStore struct {
	db     *sql.DB
	logger *log.Logger
}

func NewPostgresStore(db *sql.DB, logger *log.Logger) PostgresStore {
	return PostgresStore{
		db:     db,
		logger: logger,
	}
}

// Upsert adds a new publisher to DB if it doesn't exist.
// Returns an inserted ID, or an existing ID, if the publisher already exist.
func (s PostgresStore) Upsert(ctx context.Context, publisher string) (int64, error) {
	if publisher == "" {
		return 0, fmt.Errorf("the publisher name should not be blank")
	}

	var publisherID int64
	err := transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		selectStmt, err := tx.PrepareContext(txCtx, "SELECT id FROM ebook.publishers WHERE name = $1")
		if err != nil {
			return err
		}
		defer s.closeResource(selectStmt)

		row := selectStmt.QueryRowContext(txCtx, publisher)
		err = row.Scan(&publisherID)
		if err == nil {
			s.logger.Printf("[INFO] - Existing publisher ID: %d", publisherID)
			return nil
		}
		if err != sql.ErrNoRows {
			return err
		}

		insertStmt, err := tx.PrepareContext(txCtx, "INSERT INTO ebook.publishers(name) VALUES ($1) RETURNING id")
		if err != nil {
			return err
		}
		defer s.closeResource(insertStmt)

		if err := insertStmt.QueryRowContext(txCtx, publisher).Scan(&publisherID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	s.logger.Printf("[INFO] - Stored publisher ID: %d", publisherID)

	return publisherID, nil
}

func (s PostgresStore) closeResource(rows io.Closer) {
	err := rows.Close()
	if err != nil {
		s.logger.Printf("[ERROR] - %v", err)
	}
}
