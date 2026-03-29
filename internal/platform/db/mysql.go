package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// NewSQL membuka koneksi database menggunakan database/sql.
// Ini dipakai oleh migration runner karena migration kita berbasis SQL file manual.
func NewSQL(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
