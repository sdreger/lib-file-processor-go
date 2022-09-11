package lang

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

// Upsert adds a new language to DB if it doesn't exist.
// Returns an inserted ID, or an existing ID, if the language already exist.
func (s PostgresStore) Upsert(ctx context.Context, language string) (int64, error) {
	if language == "" {
		return 0, fmt.Errorf("the lang name should not be blank")
	}

	var languageID int64
	err := transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		selectStmt, err := tx.PrepareContext(txCtx, "SELECT id FROM ebook.languages WHERE name = $1")
		if err != nil {
			return err
		}
		defer s.closeResource(selectStmt)

		row := selectStmt.QueryRowContext(txCtx, language)
		err = row.Scan(&languageID)
		if err == nil {
			//log.Printf("[INFO] - Existing language ID: %d", languageID)
			return nil
		}
		if err != sql.ErrNoRows {
			return err
		}

		insertStmt, err := s.db.PrepareContext(txCtx, "INSERT INTO ebook.languages(name) VALUES ($1) RETURNING id")
		if err != nil {
			return err
		}
		defer s.closeResource(insertStmt)

		if err := insertStmt.QueryRowContext(txCtx, language).Scan(&languageID); err != nil {
			return err
		}
		s.logger.Printf("[INFO] - Stored language ID: %d", languageID)

		return nil
	})

	if err != nil {
		return 0, err
	}

	return languageID, nil
}

func (s PostgresStore) closeResource(rows io.Closer) {
	err := rows.Close()
	if err != nil {
		s.logger.Printf("[ERROR] - %v", err)
	}
}
