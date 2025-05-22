package database

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is your GORM handle
var DB *gorm.DB

// ConnectGORM opens a GORM Postgres connection, configures pooling,
// runs AutoMigrate on your models, and logs success/failure.
func ConnectGORM(pgDsn string, log *zap.Logger) error {
	// Open the connection
	db, err := gorm.Open(postgres.Open(pgDsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Configure the underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	DB = db
	log.Info("GORM connection established and migrated")
	return nil
}
