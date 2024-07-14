package mysqldb

import (
	"fmt"

	"github.com/developersismedika/sqlx"
	_ "github.com/go-sql-driver/mysql"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/config"
	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
)

const (
	// MysqlDriver defines the name of the mysql driver
	MysqlDriver = "mysql"
	// Percent defines the symbol that can be used to simplify the query statement building
	//  for example, when using fmt.Sprintf(), using this variable will avoid issue with double % inside the print
	Percent = "%"
	// NoRow defines the error string when no record found in the query result
	NoRow = "sql: no rows in result set"
)

type Storage struct {
	Db  *sqlx.DB
	Log *logger.Logger
}

// DbQuery defines the query parameter
//
//	it is used to wrap query param and
type DbQuery struct {
	Q string `json:"query"`
}

// QueryArgs defines the query arguments
//
//	it will be used for building each query for the Named Query Execution
type QueryArgs struct {
	Query string      `db:"query"`
	Args  interface{} `db:"args"`
}

// DbConnect opens MySQL database connection
func DbConnect(log *logger.Logger, conf config.DbMySQL) (*sqlx.DB, error) {
	// connects
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.User, conf.Pass, conf.Host, conf.Port, conf.DbName)
	db, err := sqlx.Open(MysqlDriver, dsn)
	if err != nil {
		log.Error("failed to connect to the database")
		return nil, err
	}

	// validates connection with ping
	if err = db.Ping(); err != nil {
		log.Error("failed to ping the connected database")
		return nil, err
	}

	return db, nil
}
