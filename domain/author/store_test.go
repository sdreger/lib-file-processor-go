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
)

var (
	authorsToUpsert = []string{"Bob", "John"}
)

func TestUpsertAll(t *testing.T) {
	t.Log("Given the need to test upsert authors flow")
	t.Run("Upsert two new records", testUpsertAllBothNew)
	t.Run("Upsert one new and one existing records", testUpsertAllOneNew)
	t.Run("Upsert two existing records", testUpsertAllBothExisting)
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
	store := NewStore(db)

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.authors WHERE name = ANY (.+)")
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(authorsToUpsert))
	selectQuery.WillReturnRows(rowsSelected)

	// Two new inserts
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.authors(.+) VALUES (.+) RETURNING id")
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
	store := NewStore(db)

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.authors WHERE name = ANY (.+)")
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(authorsToUpsert))
	rowToReturn := rowsSelected.AddRow(existingAuthorID, authorsToUpsert[0])
	selectQuery.WillReturnRows(rowToReturn)

	// One new insert
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.authors(.+) VALUES (.+) RETURNING id")
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
	store := NewStore(db)

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id, name FROM ebook.authors WHERE name = ANY (.+)")
	selectQuery := selectPrepare.ExpectQuery().WithArgs(pq.Array(authorsToUpsert))
	rowToReturn := rowsSelected.
		AddRow(existingAuthorID01, authorsToUpsert[0]).
		AddRow(existingAuthorID02, authorsToUpsert[1])
	selectQuery.WillReturnRows(rowToReturn)

	// No new inserts
	mock.ExpectPrepare("INSERT INTO ebook.authors(.+) VALUES (.+) RETURNING id")
	mock.ExpectCommit()

	authorIDs, err := store.UpsertAll(context.Background(), authorsToUpsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted author IDs: %v", failed, err)
	}

	assertAuthorIDs(t, expectedIDs, authorIDs)
	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add authors", succeed)
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
