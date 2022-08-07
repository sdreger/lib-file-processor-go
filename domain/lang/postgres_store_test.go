package lang

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"
)

var (
	languageToInsert = "Test Language"
)

func TestStore_Upsert(t *testing.T) {
	t.Log("Given the need to test language upsert")
	t.Run("Upsert a new record", testUpsertOneNew)
	t.Run("Upsert an existing record", testUpsertOneExisting)
	t.Run("Do not upsert language with empty language slice", testUpsertNoLanguage)
}

func testUpsertOneNew(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of a new language\n")

	var newLanguageID int64 = 1
	rowsSelected := sqlmock.NewRows([]string{"id"})
	rowsInserted := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id FROM ebook.languages WHERE name = \\$1").
		WillBeClosed()
	selectPrepare.ExpectQuery().WithArgs(languageToInsert).WillReturnRows(rowsSelected).RowsWillBeClosed()

	// One new insert
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.languages\\(name\\) VALUES \\(\\$1\\) RETURNING id").
		WillBeClosed()
	insertPrepare.ExpectQuery().WithArgs(languageToInsert).WillReturnRows(rowsInserted.AddRow(newLanguageID))
	mock.ExpectCommit()

	languageID, err := store.Upsert(context.Background(), languageToInsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted language ID: %v", failed, err)
	}

	if languageID != newLanguageID {
		t.Errorf("\t\t%s\tShould get a %d language ID: %d", failed, newLanguageID, languageID)
	}

	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add language", succeed)
}

func testUpsertOneExisting(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of an exisiting language\n")

	var existingLanguageID int64 = 1
	rowsSelected := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id FROM ebook.languages WHERE name = \\$1").
		WillBeClosed()
	selectQuery := selectPrepare.ExpectQuery().WithArgs(languageToInsert)
	rowToReturn := rowsSelected.AddRow(existingLanguageID)
	selectQuery.WillReturnRows(rowToReturn)

	// No new insert
	mock.ExpectCommit()

	languageID, err := store.Upsert(context.Background(), languageToInsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted language ID: %v", failed, err)
	}

	if languageID != existingLanguageID {
		t.Errorf("\t\t%s\tShould get a %d language ID: %d", failed, existingLanguageID, languageID)
	}

	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add language", succeed)
}

func testUpsertNoLanguage(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db)

	languageID, err := store.Upsert(context.Background(), "")
	if err == nil {
		t.Fatalf("\t\t%s\tShould not return an error when there is no language", failed)
	}

	if languageID != 0 {
		t.Fatalf("\t\t%s\tShould return an empty result when there is no language", failed)
	}

	t.Logf("\t\t%s\tShould return an empty result when there is no language", succeed)
}

func initMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to init the DB mock: %v", failed, err)
	}

	return db, mock
}

func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("\t\t%s\tShould be able to fulfill all mock expectations: %v", failed, err)
	}
}
