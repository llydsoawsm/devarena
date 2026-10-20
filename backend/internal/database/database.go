package database

import (
	"fmt"
	"log"

	"devcode/backend/internal/config"
	"devcode/backend/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	if cfg.DBDriver == "postgres" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err != nil {
			log.Printf("⚠️ PostgreSQL connection failed: %v. Falling back to SQLite database (%s)\n", err, cfg.SQLitePath)
			db, err = gorm.Open(sqlite.Open(cfg.SQLitePath), gormConfig)
		} else {
			log.Println("✅ Successfully connected to PostgreSQL database")
		}
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.SQLitePath), gormConfig)
		if err == nil {
			log.Printf("✅ Connected to SQLite database (%s)\n", cfg.SQLitePath)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Auto-migrate tables
	if err := db.AutoMigrate(
		&models.User{},
		&models.Problem{},
		&models.TestCase{},
		&models.Match{},
		&models.MatchParticipant{},
		&models.Submission{},
		&models.MatchEvent{},
	); err != nil {
		return nil, fmt.Errorf("auto-migration failed: %w", err)
	}

	DB = db
	log.Println("✅ Database schema migrated successfully")
	return db, nil
}
