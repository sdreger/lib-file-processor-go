package author

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"testing"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"

	testBookId = 1
)

var (
	authorsToUpsert = []string{"Bob", "John"}
)

func TestStore_UpsertAll(t *testing.T) {
	t.Log("Given the need to test upsert authors flow")
	t.Run("Upsert two new records", testUpsertAllBothNew)
	t.Run("Upsert one new and one existing records", testUpsertAllOneNew)
	t.Run("Upsert two existing records", testUpsertAllBothExisting)
	t.Run("Do not upsert authors with empty author slice", testUpsertAllNoAuthors)
}

func testUpsertAllBothNew(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of two new authors\n")

	var newAuthorID01 int64 = 1
	var newAuthorID02 int64 = 2
	expectedIDs := []int64{newAuthorID01, newAuthorID02}
	rowsSelected := sqlmock.NewRows([]string{"id", "name"})
	rowsInserted := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.authors WHERE name = ANY \\(\\$1\\)").
		WillBeClosed()
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(authorsToUpsert))
	selectQuery.WillReturnRows(rowsSelected)

	// Two new inserts
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.authors\\(name\\) VALUES \\(\\$1\\) RETURNING id").
		WillBeClosed()
	insertPrepare.ExpectQuery().WithArgs(authorsToUpsert[0]).WillReturnRows(rowsInserted.AddRow(newAuthorID01))
	insertPrepare.ExpectQuery().WithArgs(authorsToUpsert[1]).WillReturnRows(rowsInserted.AddRow(newAuthorID02))
	mock.ExpectCommit()

	authorIDs, err := store.UpsertAll(context.Background(), authorsToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted author IDs: %v", failed, err)
	}

	assertAuthorIDs(t, expectedIDs, authorIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add authors", succeed)
}

func testUpsertAllOneNew(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of one new authors\n")

	var existingAuthorID int64 = 1
	var newAuthorID int64 = 2
	expectedIDs := []int64{existingAuthorID, newAuthorID}
	rowsSelected := sqlmock.NewRows([]string{"id", "name"})
	rowsInserted := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.authors WHERE name = ANY \\(\\$1\\)").
		WillBeClosed()
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(authorsToUpsert))
	rowToReturn := rowsSelected.AddRow(existingAuthorID, authorsToUpsert[0])
	selectQuery.WillReturnRows(rowToReturn)

	// One new insert
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.authors\\(name\\) VALUES \\(\\$1\\) RETURNING id").
		WillBeClosed()
	insertPrepare.ExpectQuery().WithArgs(authorsToUpsert[1]).WillReturnRows(rowsInserted.AddRow(newAuthorID))
	mock.ExpectCommit()

	authorIDs, err := store.UpsertAll(context.Background(), authorsToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted author IDs: %v", failed, err)
	}

	assertAuthorIDs(t, expectedIDs, authorIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add authors", succeed)
}

func testUpsertAllBothExisting(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of no new authors\n")

	var existingAuthorID01 int64 = 1
	var existingAuthorID02 int64 = 2
	expectedIDs := []int64{existingAuthorID01, existingAuthorID02}
	rowsSelected := sqlmock.NewRows([]string{"id", "name"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.authors WHERE name = ANY \\(\\$1\\)").
		WillBeClosed()
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(authorsToUpsert))
	rowToReturn := rowsSelected.
		AddRow(existingAuthorID01, authorsToUpsert[0]).
		AddRow(existingAuthorID02, authorsToUpsert[1])
	selectQuery.WillReturnRows(rowToReturn)

	// No new inserts
	mock.ExpectPrepare("INSERT INTO ebook.authors\\(name\\) VALUES \\(\\$1\\) RETURNING id").WillBeClosed()
	mock.ExpectCommit()

	authorIDs, err := store.UpsertAll(context.Background(), authorsToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted author IDs: %v", failed, err)
	}

	assertAuthorIDs(t, expectedIDs, authorIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add authors", succeed)
}

func testUpsertAllNoAuthors(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	authorIDs, err := store.UpsertAll(context.Background(), []string{})
	if err != nil {
		t.Errorf("\t\t%s\tShould not return an error when there are no authors", failed)
	}

	if len(authorIDs) > 0 {
		t.Errorf("\t\t%s\tShould return an empty result when there are no authors", failed)
	}

	t.Logf("\t\t%s\tShould return an empty result when there are no authors", succeed)
}

func TestStore_ReplaceBookAuthors(t *testing.T) {
	t.Logf("\t\tWhen checking for book-author relations replacement\n")
	t.Run("Successfully replace book-author relations", testReplaceBookAuthors)
	t.Run("Failed to replace book-author relations with empty author slice", testReplaceBookAuthorsErrorNoAuthors)
	t.Run("Failed to replace book-author relations with bookID = 0", testReplaceBookAuthorsErrorNoBookID)
}

func testReplaceBookAuthors(t *testing.T) {
	var newAuthorID01 int64 = 1
	var newAuthorID02 int64 = 2
	authorIDs := []int64{newAuthorID01, newAuthorID02}

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	mock.ExpectBegin()
	deletePrepare := mock.ExpectPrepare("DELETE FROM ebook.book_author WHERE book_id = \\$1").WillBeClosed()
	deletePrepare.ExpectExec().WithArgs(testBookId).WillReturnResult(sqlmock.NewResult(0, 1))
	insertPrepare := mock.
		ExpectPrepare("INSERT INTO ebook.book_author\\(book_id, author_id\\) VALUES \\(\\$1, \\$2\\)").
		WillBeClosed()
	insertPrepare.ExpectExec().WithArgs(testBookId, newAuthorID01).WillReturnResult(sqlmock.NewResult(0, 1))
	insertPrepare.ExpectExec().WithArgs(testBookId, newAuthorID02).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := store.ReplaceBookAuthors(context.Background(), testBookId, authorIDs)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to add new book-author relatons: %v", failed, err)
	}
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to replace book-author relatons", succeed)
}

func testReplaceBookAuthorsErrorNoAuthors(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	err := store.ReplaceBookAuthors(context.Background(), testBookId, []int64{})
	if err == nil {
		t.Fatalf("\t\t%s\tAn error is expected when there are no authors", failed)
	}

	t.Logf("\t\t%s\tShould return an error when there are no authors", succeed)
}

func testReplaceBookAuthorsErrorNoBookID(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	err := store.ReplaceBookAuthors(context.Background(), 0, []int64{1, 2})
	if err == nil {
		t.Fatalf("\t\t%s\tAn error is expected when there is no book ID", failed)
	}

	t.Logf("\t\t%s\tShould return an error when there is no book ID", succeed)
}

func initMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to init the DB mock: %v", failed, err)
	}

	return db, mock
}

func assertAuthorIDs(t *testing.T, expectedIDs, actualIDs []int64) {
	for i, _ := range expectedIDs {
		if expectedIDs[i] != actualIDs[i] {
			t.Errorf("\t\t%s\tShould get a %d author ID: %d", failed, expectedIDs[i], actualIDs[i])
		}
	}
}

func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("\t\t%s\tShould be able to fulfill all mock expectations: %v", failed, err)
	}
}
