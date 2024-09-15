package db

import (
	"encrypt-decrypt-api/app/config"
	"encrypt-decrypt-api/app/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB stellt die Verbindung zur Datenbank her und führt die Migration durch.
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Datenbank-DSN erstellen
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	// Retry-Mechanismus für die Datenbankverbindung
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Datenbankverbindung fehlgeschlagen. Neuer Versuch in 5 Sekunden... (%d/5)", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("konnte keine Verbindung zur Datenbank herstellen: %v", err)
	}

	// Automatische Migration
	err = db.AutoMigrate(&models.File{})
	if err != nil {
		return nil, fmt.Errorf("automatische Migration fehlgeschlagen: %v", err)
	}

	return db, nil
}
