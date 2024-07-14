package storage

import (
	"github.com/developersismedika/sqlx"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/config"
	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
	mySqlx "github.com/ardihikaru/example-golang-google-chat-sender/pkg/mysqldb"
)

// DbConnect opens MySQL database connection
func DbConnect(log *logger.Logger, cfg config.DbMySQL) (*sqlx.DB, error) {
	return mySqlx.DbConnect(log, cfg)
}
