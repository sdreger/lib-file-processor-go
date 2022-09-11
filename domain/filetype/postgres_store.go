package filetype

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
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

// UpsertAll adds new file types from the input slice, existing file types are ignored.
// Returns both new and existing IDs for all file types from the input slice.
func (s PostgresStore) UpsertAll(ctx context.Context, fileTypes []string) ([]int64, error) {
	if len(fileTypes) == 0 {
		return []int64{}, nil
	}

	var ret []int64
	err := transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		existingFileTypes := make(map[string]int64)

		selectStmt, err := tx.PrepareContext(txCtx, "SELECT id, name FROM ebook.file_types WHERE name = ANY($1)")
		if err != nil {
			return err
		}
		defer s.closeResource(selectStmt)

		rows, err := selectStmt.QueryContext(txCtx, pq.Array(fileTypes))
		if err != nil {
			return err
		}
		defer s.closeResource(rows)

		for rows.Next() {
			var ID int64
			var name string
			err := rows.Scan(&ID, &name)
			if err != nil {
				return err
			}
			existingFileTypes[name] = ID
		}

		insertStmt, err := tx.PrepareContext(txCtx, "INSERT INTO ebook.file_types(name) VALUES ($1) RETURNING id")
		if err != nil {
			return err
		}
		defer s.closeResource(insertStmt)

		for _, fileType := range fileTypes {
			if existingID, ok := existingFileTypes[fileType]; ok {
				ret = append(ret, existingID)
				continue
			}
			var lastInsertId int64
			if err := insertStmt.QueryRowContext(txCtx, fileType).Scan(&lastInsertId); err != nil {
				return err
			}
			ret = append(ret, lastInsertId)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	s.logger.Printf("[INFO] - Stored fileType IDs: %d", ret)

	return ret, nil
}

// ReplaceBookFileTypes removes all records from the book-fileType join table for the particular book.
// And adds new records for all file types from the input slice.
func (s PostgresStore) ReplaceBookFileTypes(ctx context.Context, bookID int64, fileTypeIDs []int64) error {
	if bookID == 0 || len(fileTypeIDs) == 0 {
		return fmt.Errorf("there is no bookID: %q or fileTypeIDs: %v", bookID, fileTypeIDs)
	}

	return transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		deleteStmt, err := tx.PrepareContext(txCtx, "DELETE FROM ebook.book_file_type WHERE book_id = $1")
		if err != nil {
			return err
		}
		defer s.closeResource(deleteStmt)

		_, err = deleteStmt.ExecContext(txCtx, bookID)
		if err != nil {
			return err
		}

		insertStmt, err := tx.PrepareContext(txCtx, "INSERT INTO ebook.book_file_type(book_id, file_type_id) VALUES ($1, $2)")
		if err != nil {
			return err
		}
		defer s.closeResource(insertStmt)

		for _, fileTypeID := range fileTypeIDs {
			_, err := insertStmt.ExecContext(txCtx, bookID, fileTypeID)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s PostgresStore) closeResource(rows io.Closer) {
	err := rows.Close()
	if err != nil {
		s.logger.Printf("[ERROR] - %v", err)
	}
}
