package tag

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
	tagsToUpsert = []string{"pdf", "epub"}
)

func TestStore_UpsertAll(t *testing.T) {
	t.Log("Given the need to test tags upsert")
	t.Run("Upsert two new records", testUpsertAllBothNew)
	t.Run("Upsert one new and one existing records", testUpsertAllOneNew)
	t.Run("Upsert two existing records", testUpsertAllBothExisting)
	t.Run("Do not upsert tags with empty tag slice", testUpsertAllNoTags)
}

func testUpsertAllBothNew(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of two new tags\n")

	var newTagID01 int64 = 1
	var newTagID02 int64 = 2
	expectedIDs := []int64{newTagID01, newTagID02}
	rowsSelected := sqlmock.NewRows([]string{"id", "name"})
	rowsInserted := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.tags WHERE name = ANY\\(\\$1\\)").
		WillBeClosed()
	selectPrepare.ExpectQuery().WithArgs(pq.Array(tagsToUpsert)).WillReturnRows(rowsSelected).RowsWillBeClosed()

	// Two new inserts
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.tags\\(name\\) VALUES \\(\\$1\\) RETURNING id").
		WillBeClosed()
	insertPrepare.ExpectQuery().WithArgs(tagsToUpsert[0]).WillReturnRows(rowsInserted.AddRow(newTagID01))
	insertPrepare.ExpectQuery().WithArgs(tagsToUpsert[1]).WillReturnRows(rowsInserted.AddRow(newTagID02))
	mock.ExpectCommit()

	tagIDs, err := store.UpsertAll(context.Background(), tagsToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted tag IDs: %v", failed, err)
	}

	assertTagIDs(t, expectedIDs, tagIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add tags", succeed)
}

func testUpsertAllOneNew(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of one new tags\n")

	var existingTagID int64 = 1
	var newTagID int64 = 2
	expectedIDs := []int64{existingTagID, newTagID}
	rowsSelected := sqlmock.NewRows([]string{"id", "name"})
	rowsInserted := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.tags WHERE name = ANY\\(\\$1\\)").
		WillBeClosed()
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(tagsToUpsert))
	rowToReturn := rowsSelected.AddRow(existingTagID, tagsToUpsert[0])
	selectQuery.WillReturnRows(rowToReturn)

	// One new insert
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.tags\\(name\\) VALUES \\(\\$1\\) RETURNING id").
		WillBeClosed()
	insertPrepare.ExpectQuery().WithArgs(tagsToUpsert[1]).WillReturnRows(rowsInserted.AddRow(newTagID))
	mock.ExpectCommit()

	tagIDs, err := store.UpsertAll(context.Background(), tagsToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted tag IDs: %v", failed, err)
	}

	assertTagIDs(t, expectedIDs, tagIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add tags", succeed)
}

func testUpsertAllBothExisting(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of no new tags\n")

	var existingTagID01 int64 = 1
	var existingTagID02 int64 = 2
	expectedIDs := []int64{existingTagID01, existingTagID02}
	rowsSelected := sqlmock.NewRows([]string{"id", "name"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.tags WHERE name = ANY\\(\\$1\\)").
		WillBeClosed()
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(tagsToUpsert))
	rowToReturn := rowsSelected.
		AddRow(existingTagID01, tagsToUpsert[0]).
		AddRow(existingTagID02, tagsToUpsert[1])
	selectQuery.WillReturnRows(rowToReturn)

	// No new inserts
	mock.ExpectPrepare("INSERT INTO ebook.tags\\(name\\) VALUES \\(\\$1\\) RETURNING id").WillBeClosed()
	mock.ExpectCommit()

	tagIDs, err := store.UpsertAll(context.Background(), tagsToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted tag IDs: %v", failed, err)
	}

	assertTagIDs(t, expectedIDs, tagIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add tags", succeed)
}

func testUpsertAllNoTags(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	tagIDs, err := store.UpsertAll(context.Background(), []string{})
	if err != nil {
		t.Errorf("\t\t%s\tShould not return an error when there are no tags", failed)
	}

	if len(tagIDs) > 0 {
		t.Errorf("\t\t%s\tShould return an empty result when there are no tags", failed)
	}

	t.Logf("\t\t%s\tShould return an empty result when there are no tags", succeed)
}

func TestStore_ReplaceBookTags(t *testing.T) {
	t.Logf("\t\tWhen checking for book-tag relations replacement\n")
	t.Run("Successfully replace book-tag relations", testReplaceBookTags)
	t.Run("Failed to replace book-tag relations with empty tag slice", testReplaceBookTagsErrorNoTags)
	t.Run("Failed to replace book-tag relations with bookID = 0", testReplaceBookTagsErrorNoBookID)
}

func testReplaceBookTags(t *testing.T) {
	var newTagID01 int64 = 1
	var newTagID02 int64 = 2
	tagIDs := []int64{newTagID01, newTagID02}

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	mock.ExpectBegin()
	deletePrepare := mock.ExpectPrepare("DELETE FROM ebook.book_tag WHERE book_id = \\$1").WillBeClosed()
	deletePrepare.ExpectExec().WithArgs(testBookID).WillReturnResult(sqlmock.NewResult(0, 1))
	insertPrepare := mock.
		ExpectPrepare("INSERT INTO ebook.book_tag\\(book_id, tag_id\\) VALUES \\(\\$1, \\$2\\)").
		WillBeClosed()
	insertPrepare.ExpectExec().WithArgs(testBookID, newTagID01).WillReturnResult(sqlmock.NewResult(0, 1))
	insertPrepare.ExpectExec().WithArgs(testBookID, newTagID02).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := store.ReplaceBookTags(context.Background(), testBookID, tagIDs)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to add new book-tag relatons: %v", failed, err)
	}
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to replace book-tag relatons", succeed)
}

func testReplaceBookTagsErrorNoTags(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	err := store.ReplaceBookTags(context.Background(), testBookID, []int64{})
	if err == nil {
		t.Fatalf("\t\t%s\tAn error is expected when there are no tags", failed)
	}

	t.Logf("\t\t%s\tShould return an error when there are no tags", succeed)
}

func testReplaceBookTagsErrorNoBookID(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	err := store.ReplaceBookTags(context.Background(), 0, []int64{1, 2})
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

func assertTagIDs(t *testing.T, expectedIDs, actualIDs []int64) {
	for i, _ := range expectedIDs {
		if expectedIDs[i] != actualIDs[i] {
			t.Errorf("\t\t%s\tShould get a %d tag ID: %d", failed, expectedIDs[i], actualIDs[i])
		}
	}
}

func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("\t\t%s\tShould be able to fulfill all mock expectations: %v", failed, err)
	}
}
