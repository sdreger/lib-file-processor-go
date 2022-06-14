package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"
)

func TestWithTransaction(t *testing.T) {
	t.Log("Given the need to test transaction wrapper")
	t.Run("Successful execution of a new TX", testWithNewTransaction)
	t.Run("Rollback after panic during a new TX execution", testWithNewTransactionPanic)
	t.Run("Rollback after error during a new TX execution", testWithNewTransactionError)
	t.Run("Successful reusing of existent TX", testWithExistingTransactionSuccess)
	t.Run("Rollback after panic during reusing existing TX", testWithExistingTransactionPanic)
	t.Run("Rollback after error during reusing existing TX", testWithExistingTransactionError)
}
func testWithNewTransaction(t *testing.T) {
	t.Logf("\t\tWhen checking for a new transaction creation\n")

	var queryResult int64
	var expectedQueryResult int64 = 1

	db, mock := initMockDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .+ FROM test").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedQueryResult))
	mock.ExpectCommit()

	txErr := WithTransaction(context.Background(), db, func(txCtx context.Context, tx *sql.Tx) error {
		row := tx.QueryRow("SELECT id FROM test")
		err := row.Scan(&queryResult)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		t.Errorf("\t\t%s\tShould be able to execute transaction block: %v", failed, txErr)
	}

	if queryResult != expectedQueryResult {
		t.Errorf("\t\t%s\tShould get: %d query result: %d", failed, queryResult, expectedQueryResult)
	}

	assertMockExpectations(t, mock)
	t.Logf("\t\t%s\tShould be able to successfully commit transaction", succeed)
}

func testWithNewTransactionPanic(t *testing.T) {
	t.Logf("\t\tWhen checking for a new transaction creation with panic\n")

	db, mock := initMockDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	txErr := WithTransaction(context.Background(), db, func(txCtx context.Context, tx *sql.Tx) error {
		panic("panic inside new transaction")
	})

	if txErr == nil {
		t.Errorf("\t\t%s\tShould return an error after panic inside transaction block: %v", failed, txErr)
	}

	assertMockExpectations(t, mock)
	t.Logf("\t\t%s\tShould be able to rollback transaction after panic", succeed)
}

func testWithNewTransactionError(t *testing.T) {
	t.Logf("\t\tWhen checking for a new transaction creation with error\n")

	db, mock := initMockDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	txErr := WithTransaction(context.Background(), db, func(txCtx context.Context, tx *sql.Tx) error {
		return fmt.Errorf("error inside new transaction")
	})

	if txErr == nil {
		t.Errorf("\t\t%s\tShould return an error after error inside transaction block: %v", failed, txErr)
	}

	assertMockExpectations(t, mock)
	t.Logf("\t\t%s\tShould be able to rollback transaction after error", succeed)
}

func testWithExistingTransactionSuccess(t *testing.T) {
	t.Logf("\t\tWhen checking for existing transaction reusing\n")

	var queryResult int64
	var expectedQueryResult int64 = 1

	db, mock := initMockDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .+ FROM test").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedQueryResult))
	mock.ExpectCommit()

	txErr := WithTransaction(context.Background(), db, func(txCtx context.Context, tx *sql.Tx) error {
		innerTxErr := WithTransaction(txCtx, db, func(txCtx context.Context, tx *sql.Tx) error {
			row := tx.QueryRow("SELECT id FROM test")
			err := row.Scan(&queryResult)
			if err != nil {
				return err
			}

			return nil
		})

		return innerTxErr
	})

	if txErr != nil {
		t.Errorf("\t\t%s\tShould be able to execute transaction block: %v", failed, txErr)
	}

	if queryResult != expectedQueryResult {
		t.Errorf("\t\t%s\tShould get: %d query result: %d", failed, queryResult, expectedQueryResult)
	}

	assertMockExpectations(t, mock)
	t.Logf("\t\t%s\tShould be able to successfully commit transaction", succeed)
}

func testWithExistingTransactionPanic(t *testing.T) {
	t.Logf("\t\tWhen checking for existing transaction reusing with panic\n")

	db, mock := initMockDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	txErr := WithTransaction(context.Background(), db, func(txCtx context.Context, tx *sql.Tx) error {
		innerTxErr := WithTransaction(txCtx, db, func(txCtx context.Context, tx *sql.Tx) error {
			panic("panic inside new transaction")
		})

		return innerTxErr
	})

	if txErr == nil {
		t.Errorf("\t\t%s\tShould return an error after panic inside transaction block: %v", failed, txErr)
	}

	assertMockExpectations(t, mock)
	t.Logf("\t\t%s\tShould be able to rollback transaction after panic", succeed)
}

func testWithExistingTransactionError(t *testing.T) {
	t.Logf("\t\tWhen checking for existing transaction reusing with error\n")

	db, mock := initMockDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	txErr := WithTransaction(context.Background(), db, func(txCtx context.Context, tx *sql.Tx) error {
		innerTxErr := WithTransaction(txCtx, db, func(txCtx context.Context, tx *sql.Tx) error {
			return fmt.Errorf("error inside new transaction")
		})

		return innerTxErr
	})

	if txErr == nil {
		t.Errorf("\t\t%s\tShould return an error after error inside transaction block: %v", failed, txErr)
	}

	assertMockExpectations(t, mock)
	t.Logf("\t\t%s\tShould be able to rollback transaction after error", succeed)
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
