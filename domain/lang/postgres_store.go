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
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) PostgresStore {
	return PostgresStore{db: db}
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
		defer closeResource(selectStmt)

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
		defer closeResource(insertStmt)

		if err := insertStmt.QueryRowContext(txCtx, language).Scan(&languageID); err != nil {
			return err
		}
		log.Printf("[INFO] - Stored language ID: %d", languageID)

		return nil
	})

	if err != nil {
		return 0, err
	}

	return languageID, nil
}

func closeResource(rows io.Closer) {
	err := rows.Close()
	if err != nil {
		log.Printf("[ERROR] - %v", err)
	}
}
