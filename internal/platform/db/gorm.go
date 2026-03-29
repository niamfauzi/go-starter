package db

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

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
