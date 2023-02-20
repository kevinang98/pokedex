package databases

import (
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

func NewMockMySQLConn() (*sqlx.DB, sqlmock.Sqlmock) {
	db, mockSQL, _ := sqlmock.New()
	mockDB := sqlx.NewDb(db, "sqlmock")

	return mockDB, mockSQL
}
