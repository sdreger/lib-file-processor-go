package book

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/sdreger/lib-file-processor-go/domain/author"
	"github.com/sdreger/lib-file-processor-go/domain/category"
	"github.com/sdreger/lib-file-processor-go/domain/filetype"
	"github.com/sdreger/lib-file-processor-go/domain/lang"
	"github.com/sdreger/lib-file-processor-go/domain/publisher"
	"github.com/sdreger/lib-file-processor-go/domain/tag"
	"log"
	"testing"
	"time"
)

const (
	findBookQuery = `SELECT books.id, books.title, books.subtitle, books.description,
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
		WHERE \(books.title = \$1 AND books.edition = \$2 AND pub.name = \$6\) 
			OR \(books.isbn10 IS NOT NULL AND books.isbn10 = \$3\) 
			OR \(books.isbn13 IS NOT NULL AND books.isbn13 = \$4\)
			OR \(books.asin IS NOT NULL AND books.asin = \$5\)`

	addBookQuery = `INSERT INTO ebook.books\(title, subtitle, description, isbn10, isbn13, asin, pages, 
						language_id, publisher_id, publisher_url, edition, pub_date, book_file_name, book_file_size,
						cover_file_name\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11, \$12, \$13,
						\$14, \$15\) RETURNING id`

	updateBookQuery = `UPDATE ebook.books SET 
			title = \$1, subtitle = \$2, description = \$3,
			isbn10 = \$4, isbn13 = \$5, asin = \$6, pages = \$7, 
			language_id = \$8, publisher_id = \$9, publisher_url = \$10, edition = \$11, pub_date = \$12,
			book_file_name = \$13, book_file_size = \$14, cover_file_name = \$15, updated_at = NOW\(\)::timestamp
		WHERE id = \$16`
)

func TestPostgresStore_Find(t *testing.T) {
	t.Log("Given the need to test book search.")
	db, mock := initMockDB(t)
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPublisherStore := publisher.NewMockStore(ctrl)
	mockLanguageStore := lang.NewMockStore(ctrl)
	mockAuthorStore := author.NewMockStore(ctrl)
	mockCategoryStore := category.NewMockStore(ctrl)
	mockFileTypeStore := filetype.NewMockStore(ctrl)
	mockTagStore := tag.NewMockStore(ctrl)
	store := NewPostgresStore(db, mockPublisherStore, mockLanguageStore, mockAuthorStore, mockCategoryStore,
		mockFileTypeStore, mockTagStore, log.Default())

	rows := sqlmock.NewRows([]string{
		"id", "title", "subtitle", "description", "isbn10", "isbn13", "asin", "pages", "lang_name", "pub_name",
		"publisher_url", "edition", "pub_date", "book_file_name", "book_file_size", "cover_file_name", "created_at",
		"updated_at", "author_name", "category_name", "file_type_name", "tag_name",
	})
	nowTime := time.Now()
	result := rows.AddRow(testBookID, testBookTitle, testBookSubtitle, testBookDescription, testBookISBN10,
		testBookISBN13, testBookASIN, testBookPages, testBookLanguage, testBookPublisher, testBookPublisherURL,
		testBookEdition, nowTime, testBookFileName, testBookFileSize, testBookCoverFileName, nowTime, nowTime,
		testBookAuthorName, testBookCategoryName, testBookFileTypeName, testBookTagName)

	mock.ExpectBegin()
	selectStmt := mock.ExpectPrepare(findBookQuery).WillBeClosed()
	selectStmt.ExpectQuery().
		WithArgs(testBookTitle, testBookEdition, testBookISBN10, testBookISBN13, testBookASIN, testBookPublisher).
		WillReturnRows(result).RowsWillBeClosed()
	mock.ExpectCommit()

	storedData, err := store.Find(context.Background(), SearchRequest{
		Title:     testBookTitle,
		Edition:   testBookEdition,
		ISBN10:    testBookISBN10,
		ISBN13:    testBookISBN13,
		ASIN:      testBookASIN,
		Publisher: testBookPublisher,
	})

	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to find a book: %v", failed, err)
	}
	if storedData == nil {
		t.Fatalf("\t\t%s\tShould be able to find a book", failed)
	}
	if storedData.ID != testBookID {
		t.Fatalf("\t\t%s\tShould get a %d book ID: %d", failed, storedData.ID, testBookID)
	}
	if storedData.Title != testBookTitle {
		t.Fatalf("\t\t%s\tShould get a %q book title: %q", failed, storedData.Title, testBookTitle)
	}
	if storedData.Subtitle != testBookSubtitle {
		t.Fatalf("\t\t%s\tShould get a %q book subtitle: %q", failed, storedData.Subtitle, testBookSubtitle)
	}
	if storedData.Description != testBookDescription {
		t.Fatalf("\t\t%s\tShould get a %q book description: %q", failed, storedData.Description, testBookDescription)
	}
	if storedData.ISBN10 != testBookISBN10 {
		t.Fatalf("\t\t%s\tShould get a %q book ISBN10: %q", failed, storedData.ISBN10, testBookISBN10)
	}
	if storedData.ISBN13 != testBookISBN13 {
		t.Fatalf("\t\t%s\tShould get a %d book ISBN13: %d", failed, storedData.ISBN13, testBookISBN13)
	}
	if storedData.ASIN != testBookASIN {
		t.Fatalf("\t\t%s\tShould get a %q book ASIN: %q", failed, storedData.ASIN, testBookASIN)
	}
	if storedData.Pages != testBookPages {
		t.Fatalf("\t\t%s\tShould get a %d book pages: %d", failed, storedData.Pages, testBookPages)
	}
	if storedData.Language != testBookLanguage {
		t.Fatalf("\t\t%s\tShould get a %q book language: %q", failed, storedData.Language, testBookLanguage)
	}
	if storedData.Publisher != testBookPublisher {
		t.Fatalf("\t\t%s\tShould get a %q book publisher: %q", failed, storedData.Publisher, testBookPublisher)
	}
	if storedData.PublisherURL != testBookPublisherURL {
		t.Fatalf("\t\t%s\tShould get a %q book publisher URL: %q", failed, storedData.PublisherURL,
			testBookPublisherURL)
	}
	if storedData.Edition != testBookEdition {
		t.Fatalf("\t\t%s\tShould get a %d book edition: %d", failed, storedData.Edition, testBookEdition)
	}
	if storedData.PubDate != nowTime {
		t.Fatalf("\t\t%s\tShould get a %q book publish date: %q", failed, storedData.PubDate, nowTime)
	}
	if storedData.BookFileName != testBookFileName {
		t.Fatalf("\t\t%s\tShould get a %q book file name: %q", failed, storedData.BookFileName, testBookFileName)
	}
	if storedData.BookFileSize != testBookFileSize {
		t.Fatalf("\t\t%s\tShould get a %d book file size: %d", failed, storedData.BookFileSize, testBookFileSize)
	}
	if storedData.CoverFileName != testBookCoverFileName {
		t.Fatalf("\t\t%s\tShould get a %q cover file name: %q", failed, storedData.CoverFileName, testBookCoverFileName)
	}
	if storedData.CreatedAt != nowTime {
		t.Fatalf("\t\t%s\tShould get a %q book create date: %q", failed, storedData.CreatedAt, nowTime)
	}
	if storedData.UpdatedAt != nowTime {
		t.Fatalf("\t\t%s\tShould get a %q book update date: %q", failed, storedData.UpdatedAt, nowTime)
	}
	if len(storedData.Authors) != 1 && storedData.Authors[0] != testBookAuthorName {
		t.Fatalf("\t\t%s\tShould get a %q book author: %q", failed, storedData.Authors, testBookAuthorName)
	}
	if len(storedData.Categories) != 1 && storedData.Categories[0] != testBookCategoryName {
		t.Fatalf("\t\t%s\tShould get a %q book category: %q", failed, storedData.Categories, testBookCategoryName)
	}
	if len(storedData.Tags) != 1 && storedData.Tags[0] != testBookTagName {
		t.Fatalf("\t\t%s\tShould get a %q book tag: %q", failed, storedData.Tags, testBookTagName)
	}
	if len(storedData.Formats) != 1 && storedData.Formats[0] != testBookFileTypeName {
		t.Fatalf("\t\t%s\tShould get a %q book format: %q", failed, storedData.Formats, testBookFileTypeName)
	}

	t.Logf("\t\t%s\tShould successfully find a book", succeed)
}

func TestPostgresStore_Add(t *testing.T) {
	t.Log("Given the need to test book adding.")
	db, mock := initMockDB(t)
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	parsedData := getTestParsedData()

	mockPublisherStore := publisher.NewMockStore(ctrl)
	mockPublisherStore.EXPECT().Upsert(gomock.Any(), gomock.Eq(parsedData.Publisher)).
		Return(testBookPublisherID, nil).Times(1)

	mockLanguageStore := lang.NewMockStore(ctrl)
	mockLanguageStore.EXPECT().Upsert(gomock.Any(), gomock.Eq(parsedData.Language)).
		Return(testBookLanguageID, nil).Times(1)

	mockAuthorStore := author.NewMockStore(ctrl)
	mockAuthorStore.EXPECT().UpsertAll(gomock.Any(), gomock.Eq(parsedData.Authors)).
		Return([]int64{testBookAuthorID}, nil).Times(1)
	mockAuthorStore.EXPECT().ReplaceBookAuthors(gomock.Any(), testBookID, []int64{testBookAuthorID}).
		Return(nil).Times(1)

	mockCategoryStore := category.NewMockStore(ctrl)
	mockCategoryStore.EXPECT().UpsertAll(gomock.Any(), gomock.Eq(parsedData.Categories)).
		Return([]int64{testBookCategoryID}, nil).Times(1)
	mockCategoryStore.EXPECT().ReplaceBookCategories(gomock.Any(), testBookID, []int64{testBookCategoryID}).
		Return(nil).Times(1)

	mockFileTypeStore := filetype.NewMockStore(ctrl)
	mockFileTypeStore.EXPECT().UpsertAll(gomock.Any(), gomock.Eq(parsedData.Formats)).
		Return([]int64{testBookFileTypeID}, nil).Times(1)
	mockFileTypeStore.EXPECT().ReplaceBookFileTypes(gomock.Any(), testBookID, []int64{testBookFileTypeID}).
		Return(nil).Times(1)

	mockTagStore := tag.NewMockStore(ctrl)
	mockTagStore.EXPECT().UpsertAll(gomock.Any(), gomock.Eq(parsedData.Tags)).
		Return([]int64{testBookTagID}, nil).Times(1)
	mockTagStore.EXPECT().ReplaceBookTags(gomock.Any(), testBookID, []int64{testBookTagID}).
		Return(nil).Times(1)

	store := NewPostgresStore(db, mockPublisherStore, mockLanguageStore, mockAuthorStore, mockCategoryStore,
		mockFileTypeStore, mockTagStore, log.Default())

	mock.ExpectBegin()

	rowsAdd := sqlmock.NewRows([]string{"id"})
	resultAdd := rowsAdd.AddRow(testBookID)

	addStmt := mock.ExpectPrepare(addBookQuery).WillBeClosed()
	addStmt.ExpectQuery().WithArgs(testBookTitle, testBookSubtitle, testBookDescription, testBookISBN10,
		testBookISBN13, testBookASIN, testBookPages, testBookLanguageID, testBookPublisherID, testBookPublisherURL,
		testBookEdition, testPublishDate, testBookFileName, testBookFileSize, testBookCoverFileName).
		WillReturnRows(resultAdd).RowsWillBeClosed()
	mock.ExpectCommit()

	bookID, err := store.Add(context.Background(), parsedData)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to add a book: %v", failed, err)
	}

	if bookID != testBookID {
		t.Fatalf("\t\t%s\tShould get an inserted %d book ID: %d", failed, bookID, testBookID)
	}

	t.Logf("\t\t%s\tShould successfully add a book", succeed)
}

func TestPostgresStore_Update(t *testing.T) {
	t.Log("Given the need to test book updating.")
	db, mock := initMockDB(t)
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	parsedData := getTestParsedData()
	storedData := getTestStoredData()

	mockPublisherStore := publisher.NewMockStore(ctrl)
	mockPublisherStore.EXPECT().Upsert(gomock.Any(), gomock.Eq(parsedData.Publisher)).
		Return(testBookPublisherID, nil).Times(1)

	mockLanguageStore := lang.NewMockStore(ctrl)
	mockLanguageStore.EXPECT().Upsert(gomock.Any(), gomock.Eq(parsedData.Language)).
		Return(testBookLanguageID, nil).Times(1)

	mockAuthorStore := author.NewMockStore(ctrl)
	mockAuthorStore.EXPECT().UpsertAll(gomock.Any(), gomock.Eq(parsedData.Authors)).
		Return([]int64{testBookAuthorID}, nil).Times(1)
	mockAuthorStore.EXPECT().ReplaceBookAuthors(gomock.Any(), testBookID, []int64{testBookAuthorID}).
		Return(nil).Times(1)

	mockCategoryStore := category.NewMockStore(ctrl)
	mockCategoryStore.EXPECT().UpsertAll(gomock.Any(), gomock.Eq(parsedData.Categories)).
		Return([]int64{testBookCategoryID}, nil).Times(1)
	mockCategoryStore.EXPECT().ReplaceBookCategories(gomock.Any(), testBookID, []int64{testBookCategoryID}).
		Return(nil).Times(1)

	mockFileTypeStore := filetype.NewMockStore(ctrl)
	mockFileTypeStore.EXPECT().UpsertAll(gomock.Any(), gomock.Eq(parsedData.Formats)).
		Return([]int64{testBookFileTypeID}, nil).Times(1)
	mockFileTypeStore.EXPECT().ReplaceBookFileTypes(gomock.Any(), testBookID, []int64{testBookFileTypeID}).
		Return(nil).Times(1)

	mockTagStore := tag.NewMockStore(ctrl)
	mockTagStore.EXPECT().UpsertAll(gomock.Any(), gomock.Eq(parsedData.Tags)).
		Return([]int64{testBookTagID}, nil).Times(1)
	mockTagStore.EXPECT().ReplaceBookTags(gomock.Any(), testBookID, []int64{testBookTagID}).
		Return(nil).Times(1)

	store := NewPostgresStore(db, mockPublisherStore, mockLanguageStore, mockAuthorStore, mockCategoryStore,
		mockFileTypeStore, mockTagStore, log.Default())

	mock.ExpectBegin()

	updateStmt := mock.ExpectPrepare(updateBookQuery).WillBeClosed()
	updateStmt.ExpectExec().WithArgs(testBookTitle, testBookSubtitle, testBookDescription, testBookISBN10,
		testBookISBN13, testBookASIN, testBookPages, testBookLanguageID, testBookPublisherID, testBookPublisherURL,
		testBookEdition, testPublishDate, testBookFileName, testBookFileSize, testBookCoverFileName, testBookID).
		WillReturnResult(sqlmock.NewResult(testBookID, 1))
	mock.ExpectCommit()

	err := store.Update(context.Background(), &storedData, &parsedData)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to add a book: %v", failed, err)
	}

	t.Logf("\t\t%s\tShould successfully update a book", succeed)
}

func initMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to init the DB mock: %v", failed, err)
	}

	return db, mock
}
