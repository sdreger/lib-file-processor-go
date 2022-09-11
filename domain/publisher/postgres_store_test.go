package publisher

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"
)

var (
	publisherToInsert = "Test Publisher"
)

func TestStore_Upsert(t *testing.T) {
	t.Log("Given the need to test publisher upsert")
	t.Run("Upsert a new record", testUpsertOneNew)
	t.Run("Upsert an existing record", testUpsertOneExisting)
	t.Run("Do not upsert publisher with empty publisher slice", testUpsertNoPublisher)
}

func testUpsertOneNew(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of a new publisher\n")

	var newPublisherID int64 = 1
	rowsSelected := sqlmock.NewRows([]string{"id"})
	rowsInserted := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id FROM ebook.publishers WHERE name = \\$1").
		WillBeClosed()
	selectPrepare.ExpectQuery().WithArgs(publisherToInsert).WillReturnRows(rowsSelected).RowsWillBeClosed()

	// One new insert
	insertPrepare := mock.ExpectPrepare("INSERT INTO ebook.publishers\\(name\\) VALUES \\(\\$1\\) RETURNING id").
		WillBeClosed()
	insertPrepare.ExpectQuery().WithArgs(publisherToInsert).WillReturnRows(rowsInserted.AddRow(newPublisherID))
	mock.ExpectCommit()

	publisherID, err := store.Upsert(context.Background(), publisherToInsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted publisher ID: %v", failed, err)
	}

	if publisherID != newPublisherID {
		t.Errorf("\t\t%s\tShould get a %d publisher ID: %d", failed, newPublisherID, publisherID)
	}

	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add publisher", succeed)
}

func testUpsertOneExisting(t *testing.T) {
	t.Logf("\t\tWhen checking for insertion of an exisiting publisher\n")

	var existingPublisherID int64 = 1
	rowsSelected := sqlmock.NewRows([]string{"id"})

	db, mock := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	mock.ExpectBegin()
	selectPrepare := mock.ExpectPrepare("SELECT id FROM ebook.publishers WHERE name = \\$1").
		WillBeClosed()
	selectQuery := selectPrepare.ExpectQuery().WithArgs(publisherToInsert)
	rowToReturn := rowsSelected.AddRow(existingPublisherID)
	selectQuery.WillReturnRows(rowToReturn)

	// No new insert
	mock.ExpectCommit()

	publisherID, err := store.Upsert(context.Background(), publisherToInsert)
	if err != nil {
		t.Errorf("\t\t%s\tShould be able to get upserted publisher ID: %v", failed, err)
	}

	if publisherID != existingPublisherID {
		t.Errorf("\t\t%s\tShould get a %d publisher ID: %d", failed, existingPublisherID, publisherID)
	}

	assertMockExpectations(t, mock)

	t.Logf("\t\t%s\tShould be able to add publisher", succeed)
}

func testUpsertNoPublisher(t *testing.T) {
	db, _ := initMockDB(t)
	defer db.Close()
	store := NewPostgresStore(db, log.Default())

	publisherID, err := store.Upsert(context.Background(), "")
	if err == nil {
		t.Fatalf("\t\t%s\tShould not return an error when there is no publisher", failed)
	}

	if publisherID != 0 {
		t.Fatalf("\t\t%s\tShould return an empty result when there is no publisher", failed)
	}

	t.Logf("\t\t%s\tShould return an empty result when there is no publisher", succeed)
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
