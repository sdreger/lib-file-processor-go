package category

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

	testBookID = 1
)

var (
	categoriesToUpsert = []string{"Computers & Technology", "Programming"}
)

func TestStore_UpsertAll(t *testing.T) {
	t.Log("Given the need to test categories upsert")
	t.Run("Upsert two new records", testUpsertAllBothNew)
	t.Run("Upsert one new and one existing records", testUpsertAllOneNew)
	t.Run("Upsert two existing records", testUpsertAllBothExisting)
	t.Run("Do not upsert categories with empty category slice", testUpsertAllNoCategories)
}

func testUpsertAllBothNew(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of two new categories\n")

	var newCategoryID01 int64 = 1
	var newCategoryID02 int64 = 2
	expectedIDs := []int64{newCategoryID01, newCategoryID02}
	rowsSelected := sqlmock.NewRows([]string{"id", "name", "parent_id"})
	rowsInserted := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name, parent_id FROM ebook.categories WHERE name = ANY\\(\\$1\\)").
		WillBeClosed()
	selectPrepare.ExpectQuery().WithArgs(pq.Array(categoriesToUpsert)).WillReturnRows(rowsSelected).RowsWillBeClosed()

	// Two new inserts
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.categories\\(name, parent_id\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WillBeClosed()
	insertPrepare.ExpectQuery().WithArgs(categoriesToUpsert[0], nil).WillReturnRows(rowsInserted.AddRow(newCategoryID01))
	insertPrepare.ExpectQuery().WithArgs(categoriesToUpsert[1], newCategoryID01).WillReturnRows(rowsInserted.AddRow(newCategoryID02))
	mock.ExpectCommit()

	categoryIDs, err := store.UpsertAll(context.Background(), categoriesToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted category IDs: %v", failed, err)
	}

	assertCategoryIDs(t, expectedIDs, categoryIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add categories", succeed)
}

func testUpsertAllOneNew(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of one new categories\n")

	var existingCategoryID int64 = 1
	var newCategoryID int64 = 2
	expectedIDs := []int64{existingCategoryID, newCategoryID}
	rowsSelected := sqlmock.NewRows([]string{"id", "name", "parent_id"})
	rowsInserted := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name, parent_id FROM ebook.categories WHERE name = ANY\\(\\$1\\)").
		WillBeClosed()
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(categoriesToUpsert))
	rowToReturn := rowsSelected.AddRow(existingCategoryID, categoriesToUpsert[0], nil)
	selectQuery.WillReturnRows(rowToReturn)

	// One new insert
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.categories\\(name, parent_id\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WillBeClosed()
	insertPrepare.ExpectQuery().WithArgs(categoriesToUpsert[1], existingCategoryID).WillReturnRows(rowsInserted.AddRow(newCategoryID))
	mock.ExpectCommit()

	categoryIDs, err := store.UpsertAll(context.Background(), categoriesToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted category IDs: %v", failed, err)
	}

	assertCategoryIDs(t, expectedIDs, categoryIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add categories", succeed)
}

func testUpsertAllBothExisting(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of no new categories\n")

	var existingCategoryID01 int64 = 1
	var existingCategoryID02 int64 = 2
	expectedIDs := []int64{existingCategoryID01, existingCategoryID02}
	rowsSelected := sqlmock.NewRows([]string{"id", "name", "parent_id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name, parent_id FROM ebook.categories WHERE name = ANY\\(\\$1\\)").
		WillBeClosed()
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(categoriesToUpsert))
	rowToReturn := rowsSelected.
		AddRow(existingCategoryID01, categoriesToUpsert[0], nil).
		AddRow(existingCategoryID02, categoriesToUpsert[1], existingCategoryID01)
	selectQuery.WillReturnRows(rowToReturn)

	// No new inserts
	mock.ExpectPrepare("INSERT INTO ebook.categories\\(name, parent_id\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WillBeClosed()
	mock.ExpectCommit()

	categoryIDs, err := store.UpsertAll(context.Background(), categoriesToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted category IDs: %v", failed, err)
	}

	assertCategoryIDs(t, expectedIDs, categoryIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add categories", succeed)
}

func testUpsertAllNoCategories(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	categoryIDs, err := store.UpsertAll(context.Background(), []string{})
	if err != nil {
		t.Errorf("\t\t%s\tShould not return an error when there are no categories", failed)
	}

	if len(categoryIDs) > 0 {
		t.Errorf("\t\t%s\tShould return an empty result when there are no categories", failed)
	}

	t.Logf("\t\t%s\tShould return an empty result when there are no categories", succeed)
}

func TestStore_ReplaceBookCategories(t *testing.T) {
	t.Logf("\t\tWhen checking for book-category relations replacement\n")
	t.Run("Successfully replace book-category relations", testReplaceBookCategories)
	t.Run("Failed to replace book-category relations with empty category slice", testReplaceBookCategoriesErrorNoCategories)
	t.Run("Failed to replace book-category relations with bookID = 0", testReplaceBookCategoriesErrorNoBookID)
}

func testReplaceBookCategories(t *testing.T) {
	var newCategoryID01 int64 = 1
	var newCategoryID02 int64 = 2
	categoryIDs := []int64{newCategoryID01, newCategoryID02}

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	mock.ExpectBegin()
	deletePrepare := mock.ExpectPrepare("DELETE FROM ebook.book_category WHERE book_id = \\$1").WillBeClosed()
	deletePrepare.ExpectExec().WithArgs(testBookID).WillReturnResult(sqlmock.NewResult(0, 1))
	insertPrepare := mock.
		ExpectPrepare("INSERT INTO ebook.book_category\\(book_id, category_id\\) VALUES \\(\\$1, \\$2\\)").
		WillBeClosed()
	insertPrepare.ExpectExec().WithArgs(testBookID, newCategoryID01).WillReturnResult(sqlmock.NewResult(0, 1))
	insertPrepare.ExpectExec().WithArgs(testBookID, newCategoryID02).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := store.ReplaceBookCategories(context.Background(), testBookID, categoryIDs)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to add new book-category relatons: %v", failed, err)
	}
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to replace book-category relatons", succeed)
}

func testReplaceBookCategoriesErrorNoCategories(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	err := store.ReplaceBookCategories(context.Background(), testBookID, []int64{})
	if err == nil {
		t.Fatalf("\t\t%s\tAn error is expected when there are no categories", failed)
	}

	t.Logf("\t\t%s\tShould return an error when there are no categories", succeed)
}

func testReplaceBookCategoriesErrorNoBookID(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	err := store.ReplaceBookCategories(context.Background(), 0, []int64{1, 2})
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

func assertCategoryIDs(t *testing.T, expectedIDs, actualIDs []int64) {
	for i, _ := range expectedIDs {
		if expectedIDs[i] != actualIDs[i] {
			t.Errorf("\t\t%s\tShould get a %d category ID: %d", failed, expectedIDs[i], actualIDs[i])
		}
	}
}

func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("\t\t%s\tShould be able to fulfill all mock expectations: %v", failed, err)
	}
}
