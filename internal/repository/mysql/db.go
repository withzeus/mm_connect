package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() (*sql.DB, error) {
	dsn := "root@tcp(localhost:3306)/mm_connect?parseTime=true"

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to open mysql database: %v", err)
	}

	return db, db.Ping()
}
