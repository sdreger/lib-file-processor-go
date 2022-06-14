package transaction

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	TxKey = "TX"
)

// WithTransaction wrapper function is used to execute a set of DB-related actions
// within a transaction block. A new transaction will be created before function execution,
// and will be committed after its execution. In case if error is returned or panic is called
// inside the function, the transaction will be rolled back, and an error will be returned.
// If one wrapper function is wrapped inside another, then existing transaction will be used.
func WithTransaction(ctx context.Context, db *sql.DB, f func(txCtx context.Context, tx *sql.Tx) error) (err error) {
	var isNewTransaction bool
	txCtx := ctx
	ctxValue := ctx.Value(TxKey)
	var tx *sql.Tx

	if ctxValue != nil {
		//log.Printf("[DEBUG] - Using ongoing transaction")
		tx = ctxValue.(*sql.Tx)
	} else {
		//log.Printf("[DEBUG] - Creating a new transaction")
		isNewTransaction = true
		newTx, txErr := db.BeginTx(ctx, nil)
		if txErr != nil {
			return txErr
		}
		tx = newTx
		txCtx = context.WithValue(ctx, TxKey, tx)
	}

	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("panic inside the transaction block: %v", p)
			if isNewTransaction {
				//log.Printf("[WARN] - Rolling back transaction")
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					err = fmt.Errorf("can not rollback transaction: %v. %w", rollbackErr, err)
				}
			}
		} else if err != nil && isNewTransaction {
			//log.Printf("[WARN] - Rolling back transaction")
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("can not rollback transaction: %v. %w", rollbackErr, err)
			}
		} else if isNewTransaction {
			//log.Printf("[DEBUG] - Committing transaction")
			if commitErr := tx.Commit(); err != nil && isNewTransaction {
				err = fmt.Errorf("can not commit transaction: %v. %w", commitErr, err)
			}
		}
	}()

	err = f(txCtx, tx)

	return err
}
