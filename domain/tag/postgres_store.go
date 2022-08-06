package tag

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
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) PostgresStore {
	return PostgresStore{db: db}
}

// UpsertAll adds new tags from the input slice, existing tags are ignored.
// Returns both new and existing IDs for all tags from the input slice.
func (s PostgresStore) UpsertAll(ctx context.Context, tags []string) ([]int64, error) {
	if len(tags) == 0 {
		return []int64{}, nil
	}

	var ret []int64
	err := transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		existingTags := make(map[string]int64)

		selectStmt, err := tx.PrepareContext(txCtx, "SELECT id, name FROM ebook.tags WHERE name = ANY($1)")
		if err != nil {
			return err
		}
		defer closeResource(selectStmt)

		rows, err := selectStmt.QueryContext(txCtx, pq.Array(tags))
		if err != nil {
			return err
		}
		defer closeResource(rows)

		for rows.Next() {
			var ID int64
			var name string
			err := rows.Scan(&ID, &name)
			if err != nil {
				return err
			}
			existingTags[name] = ID
		}

		insertStmt, err := tx.PrepareContext(txCtx, "INSERT INTO ebook.tags(name) VALUES ($1) RETURNING id")
		if err != nil {
			return err
		}
		defer closeResource(insertStmt)

		for _, tag := range tags {
			if existingID, ok := existingTags[tag]; ok {
				ret = append(ret, existingID)
				continue
			}
			var lastInsertId int64
			if err := insertStmt.QueryRowContext(txCtx, tag).Scan(&lastInsertId); err != nil {
				return err
			}
			ret = append(ret, lastInsertId)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] - Stored tag IDs: %d", ret)

	return ret, nil
}

// ReplaceBookTags removes all records from the book-tag join table for the particular book.
// And adds new records for all tags from the input slice.
func (s PostgresStore) ReplaceBookTags(ctx context.Context, bookID int64, tagIDs []int64) error {
	if bookID == 0 {
		return fmt.Errorf("there is no bookID: %q", bookID)
	}

	return transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		deleteStmt, err := tx.PrepareContext(txCtx, "DELETE FROM ebook.book_tag WHERE book_id = $1")
		if err != nil {
			return err
		}
		defer closeResource(deleteStmt)

		_, err = deleteStmt.ExecContext(txCtx, bookID)
		if err != nil {
			return err
		}

		insertStmt, err := tx.PrepareContext(txCtx, "INSERT INTO ebook.book_tag(book_id, tag_id) VALUES ($1, $2)")
		if err != nil {
			return err
		}
		defer closeResource(insertStmt)

		for _, tagID := range tagIDs {
			_, err := insertStmt.ExecContext(txCtx, bookID, tagID)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func closeResource(rows io.Closer) {
	err := rows.Close()
	if err != nil {
		log.Printf("[ERROR] - %v", err)
	}
}
