package filetype

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"log"
	"testing"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"

	testBookID = 1
)

var (
	fileTypesToUpsert = []string{"pdf", "epub"}
)

func TestStore_UpsertAll(t *testing.T) {
	t.Log("Given the need to test file types upsert")
	t.Run("Upsert two new records", testUpsertAllBothNew)
	t.Run("Upsert one new and one existing records", testUpsertAllOneNew)
	t.Run("Upsert two existing records", testUpsertAllBothExisting)
	t.Run("Do not upsert fileTypes with empty fileType slice", testUpsertAllNoFileTypes)
}

func testUpsertAllBothNew(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of two new file types\n")

	var newFileTypeID01 int64 = 1
	var newFileTypeID02 int64 = 2
	expectedIDs := []int64{newFileTypeID01, newFileTypeID02}
	rowsSelected := sqlmock.NewRows([]string{"id", "name"})
	rowsInserted := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.file_types WHERE name = ANY\\(\\$1\\)").
		WillBeClosed()
	selectPrepare.ExpectQuery().WithArgs(pq.Array(fileTypesToUpsert)).WillReturnRows(rowsSelected).RowsWillBeClosed()

	// Two new inserts
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.file_types\\(name\\) VALUES \\(\\$1\\) RETURNING id").
		WillBeClosed()
	insertPrepare.ExpectQuery().WithArgs(fileTypesToUpsert[0]).WillReturnRows(rowsInserted.AddRow(newFileTypeID01))
	insertPrepare.ExpectQuery().WithArgs(fileTypesToUpsert[1]).WillReturnRows(rowsInserted.AddRow(newFileTypeID02))
	mock.ExpectCommit()

	fileTypeIDs, err := store.UpsertAll(context.Background(), fileTypesToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted fileType IDs: %v", failed, err)
	}

	assertFileTypeIDs(t, expectedIDs, fileTypeIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add fileTypes", succeed)
}

func testUpsertAllOneNew(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of one new file types\n")

	var existingFileTypeID int64 = 1
	var newFileTypeID int64 = 2
	expectedIDs := []int64{existingFileTypeID, newFileTypeID}
	rowsSelected := sqlmock.NewRows([]string{"id", "name"})
	rowsInserted := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.file_types WHERE name = ANY\\(\\$1\\)").
		WillBeClosed()
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(fileTypesToUpsert))
	rowToReturn := rowsSelected.AddRow(existingFileTypeID, fileTypesToUpsert[0])
	selectQuery.WillReturnRows(rowToReturn)

	// One new insert
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.file_types\\(name\\) VALUES \\(\\$1\\) RETURNING id").
		WillBeClosed()
	insertPrepare.ExpectQuery().WithArgs(fileTypesToUpsert[1]).WillReturnRows(rowsInserted.AddRow(newFileTypeID))
	mock.ExpectCommit()

	fileTypeIDs, err := store.UpsertAll(context.Background(), fileTypesToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted fileType IDs: %v", failed, err)
	}

	assertFileTypeIDs(t, expectedIDs, fileTypeIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add fileTypes", succeed)
}

func testUpsertAllBothExisting(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of no new file types\n")

	var existingFileTypeID01 int64 = 1
	var existingFileTypeID02 int64 = 2
	expectedIDs := []int64{existingFileTypeID01, existingFileTypeID02}
	rowsSelected := sqlmock.NewRows([]string{"id", "name"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.file_types WHERE name = ANY\\(\\$1\\)").
		WillBeClosed()
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(fileTypesToUpsert))
	rowToReturn := rowsSelected.
		AddRow(existingFileTypeID01, fileTypesToUpsert[0]).
		AddRow(existingFileTypeID02, fileTypesToUpsert[1])
	selectQuery.WillReturnRows(rowToReturn)

	// No new inserts
	mock.ExpectPrepare("INSERT INTO ebook.file_types\\(name\\) VALUES \\(\\$1\\) RETURNING id").WillBeClosed()
	mock.ExpectCommit()

	fileTypeIDs, err := store.UpsertAll(context.Background(), fileTypesToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted fileType IDs: %v", failed, err)
	}

	assertFileTypeIDs(t, expectedIDs, fileTypeIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add fileTypes", succeed)
}

func testUpsertAllNoFileTypes(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	fileTypeIDs, err := store.UpsertAll(context.Background(), []string{})
	if err != nil {
		t.Errorf("\t\t%s\tShould not return an error when there are no file types", failed)
	}

	if len(fileTypeIDs) > 0 {
		t.Errorf("\t\t%s\tShould return an empty result when there are no file types", failed)
	}

	t.Logf("\t\t%s\tShould return an empty result when there are no file types", succeed)
}

func TestStore_ReplaceBookFileTypes(t *testing.T) {
	t.Logf("\t\tWhen checking for book-fileType relations replacement\n")
	t.Run("Successfully replace book-fileType relations", testReplaceBookFileTypes)
	t.Run("Failed to replace book-fileType relations with empty fileType slice", testReplaceBookFileTypesErrorNoFileTypes)
	t.Run("Failed to replace book-fileType relations with bookID = 0", testReplaceBookFileTypesErrorNoBookID)
}

func testReplaceBookFileTypes(t *testing.T) {
	var newFileTypeID01 int64 = 1
	var newFileTypeID02 int64 = 2
	fileTypeIDs := []int64{newFileTypeID01, newFileTypeID02}

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	mock.ExpectBegin()
	deletePrepare := mock.ExpectPrepare("DELETE FROM ebook.book_file_type WHERE book_id = \\$1").WillBeClosed()
	deletePrepare.ExpectExec().WithArgs(testBookID).WillReturnResult(sqlmock.NewResult(0, 1))
	insertPrepare := mock.
		ExpectPrepare("INSERT INTO ebook.book_file_type\\(book_id, file_type_id\\) VALUES \\(\\$1, \\$2\\)").
		WillBeClosed()
	insertPrepare.ExpectExec().WithArgs(testBookID, newFileTypeID01).WillReturnResult(sqlmock.NewResult(0, 1))
	insertPrepare.ExpectExec().WithArgs(testBookID, newFileTypeID02).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := store.ReplaceBookFileTypes(context.Background(), testBookID, fileTypeIDs)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to add new book-fileType relatons: %v", failed, err)
	}
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to replace book-fileType relatons", succeed)
}

func testReplaceBookFileTypesErrorNoFileTypes(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	err := store.ReplaceBookFileTypes(context.Background(), testBookID, []int64{})
	if err == nil {
		t.Fatalf("\t\t%s\tAn error is expected when there are no fileTypes", failed)
	}

	t.Logf("\t\t%s\tShould return an error when there are no fileTypes", succeed)
}

func testReplaceBookFileTypesErrorNoBookID(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	err := store.ReplaceBookFileTypes(context.Background(), 0, []int64{1, 2})
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

func assertFileTypeIDs(t *testing.T, expectedIDs, actualIDs []int64) {
	for i, _ := range expectedIDs {
		if expectedIDs[i] != actualIDs[i] {
			t.Errorf("\t\t%s\tShould get a %d fileType ID: %d", failed, expectedIDs[i], actualIDs[i])
		}
	}
}

func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("\t\t%s\tShould be able to fulfill all mock expectations: %v", failed, err)
	}
}
