package database

import (
	"context"
	"database/sql"
	"log"
)

// Transactional executes a function within a database transaction.
func Transactional(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Panic(p) // re-throw after rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // if Commit returns error update err with commit err
		}
	}()

	err = fn(tx)
	return err
}
