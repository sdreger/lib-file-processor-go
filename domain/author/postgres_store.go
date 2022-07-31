package author

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

// UpsertAll adds new authors from the input slice, existing authors are ignored.
// Returns both new and existing IDs for all authors from the input slice.
func (s PostgresStore) UpsertAll(ctx context.Context, authors []string) ([]int64, error) {
	if len(authors) == 0 {
		return []int64{}, nil
	}

	var ret []int64
	err := transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		existingAuthors := make(map[string]int64)
		stmt, err := tx.PrepareContext(txCtx, "SELECT id, name FROM ebook.authors WHERE name = ANY ($1)")
		if err != nil {
			return err
		}
		defer closeResource(stmt)

		rows, err := stmt.QueryContext(txCtx, pq.Array(authors))
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
			existingAuthors[name] = ID
		}

		stmt, err = tx.PrepareContext(txCtx, "INSERT INTO ebook.authors(name) VALUES ($1) RETURNING id")
		if err != nil {
			return err
		}
		defer closeResource(stmt)

		for _, author := range authors {
			if existingID, ok := existingAuthors[author]; ok {
				ret = append(ret, existingID)
				continue
			}
			var lastInsertId int64
			if err := stmt.QueryRowContext(txCtx, author).Scan(&lastInsertId); err != nil {
				return err
			}
			ret = append(ret, lastInsertId)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// ReplaceBookAuthors removes all records from the book-author join table for the particular book.
// And adds new records for all authors from the input slice.
func (s PostgresStore) ReplaceBookAuthors(ctx context.Context, bookID int64, authorIDs []int64) error {
	if bookID == 0 || len(authorIDs) == 0 {
		return fmt.Errorf("there is no bookID: %q or authorIDs: %v", bookID, authorIDs)
	}

	return transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		deleteStmt, err := tx.PrepareContext(txCtx, "DELETE FROM ebook.book_author WHERE book_id = $1")
		if err != nil {
			return err
		}
		defer closeResource(deleteStmt)

		_, err = deleteStmt.ExecContext(txCtx, bookID)
		if err != nil {
			return err
		}

		insertStmt, err := tx.PrepareContext(txCtx, "INSERT INTO ebook.book_author(book_id, author_id) VALUES ($1, $2)")
		if err != nil {
			return err
		}
		defer closeResource(insertStmt)

		for _, authorID := range authorIDs {
			_, err := insertStmt.ExecContext(txCtx, bookID, authorID)
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
