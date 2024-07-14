package mysqldb

import (
	"database/sql"
	"fmt"
)

// Exec executes a single query with a transaction
func (d *Storage) Exec(qArgs QueryArgs) error {
	var err error

	tx := d.Db.MustBegin()

	exec, err := tx.NamedExec(qArgs.Query, qArgs.Args)
	if err != nil {
		d.Log.Warn(fmt.Sprintf("query execution failed. rolling back the changes"))
		_ = tx.Rollback()

		return err
	}

	// debug only: print results
	d.insertOnePrint(exec)

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// ExecMany executes multiple queries with a transaction
func (d *Storage) ExecMany(qArgsList []QueryArgs) error {
	var err error

	tx := d.Db.MustBegin()

	for k, _ := range qArgsList {
		exec, err := tx.NamedExec(qArgsList[k].Query, qArgsList[k].Args)
		if err != nil {
			d.Log.Warn(fmt.Sprintf("query execution failed. rolling back the changes"))
			_ = tx.Rollback()

			return err
		}

		// debug only: print results
		d.insertOnePrint(exec)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// insertOnePrint prints out the query result
func (d *Storage) insertOnePrint(exec sql.Result) {
	// prints log on debug mode only
	lastInsertedId, err := exec.LastInsertId()
	if err != nil {
		d.Log.Debug(fmt.Sprintf("last inserted ID: %v", lastInsertedId))
	}

	// prints log on debug mode only
	rowsAffected, err := exec.RowsAffected()
	if err != nil {
		d.Log.Debug(fmt.Sprintf("Number of row affected: %v", rowsAffected))
	}
}
