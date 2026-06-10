package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustOpenDB(cfg Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	return db
}
