package repository

import (
	"database/sql"
)

type Repository interface {
	ExecuteQuery(q string, params ...interface{}) (sql.Result, error)
	RunQuery(q string, params ...interface{}) (*sql.Rows, error)
	DoTransaction(qs []string, params [][]interface{}) error
}
