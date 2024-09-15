package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBHost         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBPort         string
	AllowedOrigins string
}

func NewConfig() *Config {
	// .env-Datei laden
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Fehler beim Laden der .env-Datei")
	}

	// Konfiguration erstellen
	cfg := &Config{
		DBHost:         os.Getenv("DB_HOST"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		DBPort:         os.Getenv("DB_PORT"),
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
	}

	// Überprüfen, ob alle Variablen gesetzt sind
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" || cfg.DBPort == "" {
		log.Fatal("Fehler: Datenbank-Konfigurationsvariablen sind nicht alle gesetzt.")
	}

	return cfg
}
