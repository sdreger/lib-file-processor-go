package author

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/sdreger/lib-file-processor-go/db/transaction"
	"io"
	"log"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{db: db}
}

func (s Store) UpsertAll(ctx context.Context, authors []string) ([]int64, error) {
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
	log.Printf("[INFO] - Stored author IDs: %d", ret)

	return ret, nil
}

func closeResource(rows io.Closer) {
	err := rows.Close()
	if err != nil {
		log.Printf("[ERROR] - %v", err)
	}
}
