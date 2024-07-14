package queue

import (
	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/mysqldb"
)

type Store mysqldb.Storage

// Insert inserts data
func (store *Store) Insert() error {
	// TODO: implements code here
	store.Log.Debug("stores data to the queue table")

	return nil
}
