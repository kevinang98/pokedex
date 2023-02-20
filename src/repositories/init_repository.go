package repositories

import (
	"pokedex/src/databases"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func initRepositoryMock(t *testing.T) (Repository, sqlmock.Sqlmock) {
	mysqlDB, mock := databases.NewMockMySQLConn()

	mongoDB, _ := databases.NewMockMongoConn(t)

	repo := NewRepository(mysqlDB, mongoDB)

	return *repo, mock
}
