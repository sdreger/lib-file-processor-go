package book

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sdreger/lib-file-processor-go/db/transaction"
	"github.com/sdreger/lib-file-processor-go/domain/author"
	"github.com/sdreger/lib-file-processor-go/domain/category"
	"github.com/sdreger/lib-file-processor-go/domain/filetype"
	"github.com/sdreger/lib-file-processor-go/domain/lang"
	"github.com/sdreger/lib-file-processor-go/domain/publisher"
	"github.com/sdreger/lib-file-processor-go/domain/tag"
	"io"
	"log"
)

type PostgresStore struct {
	db             *sql.DB
	publisherStore publisher.Store
	languageStore  lang.Store
	authorStore    author.Store
	categoryStore  category.Store
	filetypeStore  filetype.Store
	tagStore       tag.Store
	logger         *log.Logger
}

func NewPostgresStore(db *sql.DB, publisherStore publisher.Store, languageStore lang.Store, authorStore author.Store,
	categoryStore category.Store, filetypeStore filetype.Store, tagStore tag.Store, logger *log.Logger) Store {
	return PostgresStore{
		db:             db,
		publisherStore: publisherStore,
		languageStore:  languageStore,
		authorStore:    authorStore,
		categoryStore:  categoryStore,
		filetypeStore:  filetypeStore,
		tagStore:       tagStore,
		logger:         logger,
	}
}

func (s PostgresStore) Find(ctx context.Context, req SearchRequest) (*StoredData, error) {

	var book StoredData
	err := transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		selectQuery := `SELECT books.id, books.title, books.subtitle, books.description,
		books.isbn10, books.isbn13, books.asin, books.pages, lang.name AS lang_name, pub.name AS pub_name,
        books.publisher_url, books.edition, books.pub_date,
        books.book_file_name, books.book_file_size, books.cover_file_name, books.created_at, books.updated_at,
        a.name AS author_name, c.name AS category_name, ft.name AS file_type_name, t.name AS tag_name
		FROM ebook.books
			LEFT JOIN ebook.publishers pub ON books.publisher_id = pub.id
			LEFT JOIN ebook.languages lang ON books.language_id = lang.id
			LEFT JOIN ebook.book_author ba on books.id = ba.book_id
			LEFT JOIN ebook.authors a on a.id = ba.author_id
			LEFT JOIN ebook.book_category bc on books.id = bc.book_id
			LEFT JOIN ebook.categories c on c.id = bc.category_id
			LEFT JOIN ebook.book_file_type bft on books.id = bft.book_id
			LEFT JOIN ebook.file_types ft on ft.id = bft.file_type_id
			LEFT JOIN ebook.book_tag bt on books.id = bt.book_id
			LEFT JOIN ebook.tags t on t.id = bt.tag_id
		WHERE (books.title = $1 AND books.edition = $2) 
			OR (books.isbn10 IS NOT NULL AND books.isbn10 = $3) 
			OR (books.isbn13 IS NOT NULL AND books.isbn13 = $4)
			OR (books.asin IS NOT NULL AND books.asin = $5)`
		selectStmt, err := tx.PrepareContext(txCtx, selectQuery)
		if err != nil {
			return err
		}
		defer s.closeResource(selectStmt)

		rows, err := selectStmt.QueryContext(txCtx, req.Title, req.Edition, req.ISBN10, req.ISBN13, req.ASIN)
		if err != nil {
			return err
		}
		defer s.closeResource(rows)
		bookData, err := scanBookData(rows)
		if err != nil {
			return err
		}
		book = bookData
		return nil
	})

	if err != nil {
		return nil, err
	}

	if book.IsEmpty() {
		return nil, nil
	}

	return &book, nil
}

func scanBookData(rows *sql.Rows) (StoredData, error) {
	var bookData StoredData
	for rows.Next() {
		var rowData dotProductRow
		err := rows.Scan(&rowData.ID, &rowData.Title, &rowData.Subtitle, &rowData.Description, &rowData.ISBN10,
			&rowData.ISBN13, &rowData.ASIN, &rowData.Pages, &rowData.Language, &rowData.Publisher,
			&rowData.PublisherURL, &rowData.Edition, &rowData.PubDate, &rowData.BookFileName, &rowData.BookFileSize,
			&rowData.CoverFileName, &rowData.CreatedAt, &rowData.UpdatedAt, &rowData.AuthorName, &rowData.CategoryName,
			&rowData.FileTypeName, &rowData.TagName)
		if err != nil {
			return StoredData{}, err
		}
		mapToStoredData(rowData, &bookData)
	}

	bookData.Authors = deduplicateMappedData(bookData.Authors)
	bookData.Categories = deduplicateMappedData(bookData.Categories)
	bookData.Formats = deduplicateMappedData(bookData.Formats)
	bookData.Tags = deduplicateMappedData(bookData.Tags)

	return bookData, nil
}

func (s PostgresStore) Add(ctx context.Context, parsedData ParsedData) (int64, error) {

	var bookID int64
	err := transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {

		// ---------- Store book relations data ----------
		relKeys, err := s.storeBookRelationData(txCtx, parsedData)
		if err != nil {
			return fmt.Errorf("can not store book relations data: %w", err)
		}

		// ---------- Store book record ----------
		insertQuery := `INSERT INTO ebook.books(title, subtitle, description, isbn10, isbn13, asin, pages, language_id, 
                        publisher_id, publisher_url, edition, pub_date, book_file_name, book_file_size, cover_file_name)
                		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`
		insertStmt, err := tx.PrepareContext(txCtx, insertQuery)
		if err != nil {
			return err
		}
		defer s.closeResource(insertStmt)

		isbn10, isbn13, asin := getBookIdentifiers(parsedData)
		bookIDRow := insertStmt.QueryRowContext(txCtx, parsedData.Title, getNullableString(parsedData.Subtitle),
			parsedData.Description, isbn10, isbn13, asin, parsedData.Pages, relKeys.languageID, relKeys.publisherID,
			parsedData.PublisherURL, parsedData.Edition, parsedData.PubDate, parsedData.BookFileName,
			parsedData.BookFileSize, parsedData.CoverFileName)
		bookStoreErr := bookIDRow.Scan(&bookID)
		if bookStoreErr != nil {
			return fmt.Errorf("can not store book: %w", bookStoreErr)
		}

		// ---------- Store book relation links ----------
		err = s.storeBookRelationLinks(txCtx, bookID, relKeys)
		if err != nil {
			return fmt.Errorf("can not store book relations: %w", err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	s.logger.Printf("[INFO] - Stored book ID: %d", bookID)

	return bookID, nil
}

func (s PostgresStore) Update(ctx context.Context, existingData *StoredData, parsedData *ParsedData) error {

	mapToParsedDate(existingData, parsedData)

	err := transaction.WithTransaction(ctx, s.db, func(txCtx context.Context, tx *sql.Tx) error {
		// ---------- Store book relations data ----------
		relKeys, err := s.storeBookRelationData(txCtx, *parsedData)
		if err != nil {
			return fmt.Errorf("can not update book relations data: %w", err)
		}

		// ---------- Update book record ----------
		updateQuery := `UPDATE ebook.books SET 
			title = $1, subtitle = $2, description = $3,
			isbn10 = $4, isbn13 = $5, asin = $6, pages = $7, 
			language_id = $8, publisher_id = $9, publisher_url = $10, edition = $11, pub_date = $12,
			book_file_name = $13, book_file_size = $14, cover_file_name = $15, updated_at = NOW()::timestamp
		WHERE id = $16`
		updateStmt, err := tx.PrepareContext(txCtx, updateQuery)
		if err != nil {
			return err
		}
		defer s.closeResource(updateStmt)

		isbn10, isbn13, asin := getBookIdentifiers(*parsedData)
		_, bookUpdateErr := updateStmt.ExecContext(txCtx, parsedData.Title, parsedData.Subtitle,
			parsedData.Description, isbn10, isbn13, asin, parsedData.Pages,
			relKeys.languageID, relKeys.publisherID, parsedData.PublisherURL, parsedData.Edition, parsedData.PubDate,
			parsedData.BookFileName, parsedData.BookFileSize, parsedData.CoverFileName, existingData.ID)
		if bookUpdateErr != nil {
			return fmt.Errorf("can not update book: %w", err)
		}

		// ---------- Store book relation links ----------
		err = s.storeBookRelationLinks(txCtx, existingData.ID, relKeys)
		if err != nil {
			return fmt.Errorf("can not update book relations: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}
	s.logger.Printf("[INFO] - Updated book ID: %d", existingData.ID)

	return nil
}

func (s PostgresStore) storeBookRelationData(txCtx context.Context, parsedData ParsedData) (relationKeys, error) {
	// ---------- Upsert publisher ----------
	publisherID, err := s.publisherStore.Upsert(txCtx, parsedData.Publisher)
	if err != nil {
		return relationKeys{}, fmt.Errorf("can not upsert publisher: %w", err)
	}

	// ---------- Upsert language ----------
	languageID, err := s.languageStore.Upsert(txCtx, parsedData.Language)
	if err != nil {
		return relationKeys{}, fmt.Errorf("can not upsert language: %w", err)
	}

	// ---------- Upsert authors ----------
	authorIDs, err := s.authorStore.UpsertAll(txCtx, parsedData.Authors)
	if err != nil {
		return relationKeys{}, fmt.Errorf("can not upsert authors: %w", err)
	}

	// ---------- Upsert categories ----------
	categoryIDs, err := s.categoryStore.UpsertAll(txCtx, parsedData.Categories)
	if err != nil {
		return relationKeys{}, fmt.Errorf("can not upsert categories: %w", err)
	}

	// ---------- Upsert file types ----------
	fileTypeIDs, err := s.filetypeStore.UpsertAll(txCtx, parsedData.Formats)
	if err != nil {
		return relationKeys{}, fmt.Errorf("can not upsert file types: %w", err)
	}

	// ---------- Upsert tags ----------
	tagIDs, err := s.tagStore.UpsertAll(txCtx, parsedData.Tags)
	if err != nil {
		return relationKeys{}, fmt.Errorf("can not upsert tags: %w", err)
	}

	return relationKeys{
		publisherID: publisherID,
		languageID:  languageID,
		authorIDs:   authorIDs,
		categoryIDs: categoryIDs,
		fileTypeIDs: fileTypeIDs,
		tagIDs:      tagIDs,
	}, nil
}

func (s PostgresStore) storeBookRelationLinks(txCtx context.Context, bookID int64, relKeys relationKeys) error {
	err := s.authorStore.ReplaceBookAuthors(txCtx, bookID, relKeys.authorIDs)
	if err != nil {
		return fmt.Errorf("can not update book-author relations: %w", err)
	}

	err = s.categoryStore.ReplaceBookCategories(txCtx, bookID, relKeys.categoryIDs)
	if err != nil {
		return fmt.Errorf("can not update book-category relations: %w", err)
	}

	err = s.filetypeStore.ReplaceBookFileTypes(txCtx, bookID, relKeys.fileTypeIDs)
	if err != nil {
		return fmt.Errorf("can not update book-fileType relations: %w", err)
	}

	err = s.tagStore.ReplaceBookTags(txCtx, bookID, relKeys.tagIDs)
	if err != nil {
		return fmt.Errorf("can not update book-tag relations: %w", err)
	}

	return nil
}

func getBookIdentifiers(parsedData ParsedData) (sql.NullString, sql.NullInt64, sql.NullString) {
	var isbn13 sql.NullInt64
	var isbn10, asin sql.NullString
	if parsedData.ISBN10 != "" {
		isbn10 = getNullableString(parsedData.ISBN10)
	}
	if parsedData.ISBN13 > 0 {
		isbn13 = getNullableInt64(parsedData.ISBN13)
	}
	if parsedData.ASIN != "" {
		asin = getNullableString(parsedData.ASIN)
	}
	return isbn10, isbn13, asin
}

func getNullableInt64(val int64) sql.NullInt64 {
	nullInt64 := sql.NullInt64{Int64: val}
	if val > 0 {
		nullInt64.Valid = true
	}
	return nullInt64
}

func getNullableString(val string) sql.NullString {
	nullString := sql.NullString{String: val}
	if val != "" {
		nullString.Valid = true
	}
	return nullString
}

func (s PostgresStore) closeResource(rows io.Closer) {
	err := rows.Close()
	if err != nil {
		s.logger.Printf("[ERROR] - %v", err)
	}
}
