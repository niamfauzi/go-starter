package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
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

// NewGorm membuka koneksi GORM untuk query aplikasi sehari-hari.
func NewGorm(dsn string, appLogger *zap.Logger) (*gorm.DB, error) {
	gormLogLevel := gormlogger.Warn
	if appLogger != nil {
		gormLogLevel = gormlogger.Info
	}

	gormDB, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.New(
			log.New(zap.NewStdLog(appLogger).Writer(), "", log.LstdFlags),
			gormlogger.Config{
				SlowThreshold:             500 * time.Millisecond,
				LogLevel:                  gormLogLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql gagal: %w", err)
	}

	return gormDB, nil
}
