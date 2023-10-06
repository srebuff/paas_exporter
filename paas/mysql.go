package paas

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func ConnectMysql(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return err // proper error handling instead of panic in your app
	}
	return nil
}
