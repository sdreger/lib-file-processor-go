package category

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

// UpsertAll adds only new categories to DB, and returns the list of all storedData IDs from the input list.
// The order in_book the input list matters. It should contain storedData hierarchy: [0] > [1] > [2].
// Where the [0] element is the parent of [1] element, [1] is the parent of [2], and so on.
func (s PostgresStore) UpsertAll(ctx context.Context, categories []string) ([]int64, error) {
	if len(categories) == 0 {
		return []int64{}, nil
	}

	var ret []int64
	err := transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		existingCategories := make(map[string]storedData)

		stmt, err := tx.PrepareContext(txCtx, "SELECT id, name, parent_id FROM ebook.categories WHERE name = ANY($1)")
		if err != nil {
			return err
		}
		defer s.closeResource(stmt)

		result, err := stmt.QueryContext(txCtx, pq.Array(categories))
		if err != nil {
			return err
		}
		defer s.closeResource(result)

		for result.Next() {
			var cat storedData
			err := result.Scan(&cat.ID, &cat.name, &cat.parentID)
			if err != nil {
				return err
			}
			existingCategories[cat.name] = cat
		}

		insertStmt, err := tx.PrepareContext(txCtx, "INSERT INTO ebook.categories(name, parent_id) VALUES ($1, $2) RETURNING id")
		if err != nil {
			return err
		}
		defer s.closeResource(insertStmt)

		for i, cat := range categories {
			if existingCategory, ok := existingCategories[cat]; ok {
				ret = append(ret, existingCategory.ID)
				continue
			}

			var parentIDNullable sql.NullInt64
			if i > 0 {
				parentIDNullable = getNullableInt64(existingCategories[categories[i-1]].ID)
			}

			var lastInsertId int64
			err = insertStmt.QueryRowContext(txCtx, cat, parentIDNullable).Scan(&lastInsertId)
			if err != nil {
				return err
			}

			existingCategories[cat] = storedData{
				ID:       lastInsertId,
				name:     cat,
				parentID: parentIDNullable,
			}
			ret = append(ret, lastInsertId)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	s.logger.Printf("[INFO] - Stored category IDs: %d", ret)

	return ret, nil
}

// ReplaceBookCategories removes all records from the book-category join table for the particular book.
// And adds new records for all categories from the input slice.
func (s PostgresStore) ReplaceBookCategories(ctx context.Context, bookID int64, categoryIDs []int64) error {
	if bookID == 0 || len(categoryIDs) == 0 {
		return fmt.Errorf("there is no bookID: %q or categoryIDs: %v", bookID, categoryIDs)
	}

	return transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		deleteStmt, err := tx.PrepareContext(txCtx, "DELETE FROM ebook.book_category WHERE book_id = $1")
		if err != nil {
			return err
		}
		defer s.closeResource(deleteStmt)

		_, err = deleteStmt.ExecContext(txCtx, bookID)
		if err != nil {
			return err
		}

		insertStmt, err := tx.PrepareContext(txCtx, "INSERT INTO ebook.book_category(book_id, category_id) VALUES ($1, $2)")
		if err != nil {
			return err
		}
		defer s.closeResource(insertStmt)

		for _, categoryID := range categoryIDs {
			_, err := insertStmt.ExecContext(txCtx, bookID, categoryID)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func getNullableInt64(val int64) sql.NullInt64 {
	nullInt64 := sql.NullInt64{Int64: val}
	if val > 0 {
		nullInt64.Valid = true
	}
	return nullInt64
}

func (s PostgresStore) closeResource(rows io.Closer) {
	err := rows.Close()
	if err != nil {
		s.logger.Printf("[ERROR] - %v", err)
	}
}
