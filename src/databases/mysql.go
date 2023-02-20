package databases

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"

	"github.com/jmoiron/sqlx"
)

func NewMySQLConn() (*sqlx.DB, error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", user, pass, host, port, name)

	db, err := sqlx.Open("mysql", source)
	if err != nil {
		return db, err
	}

	if err = db.Ping(); err != nil {
		return db, err
	}

	if os.Getenv("GOOSE") == "ENABLE" {
		if err := goose.SetDialect("mysql"); err != nil {
			return db, err
		}

		if err = goose.Run("up", db.DB, "./src/databases/ddl"); err != nil {
			return db, err
		}
	}

	return db, nil
}
