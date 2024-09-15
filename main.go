package main

import (
	"encrypt-decrypt-api/app/config"
	"encrypt-decrypt-api/app/db"
	"encrypt-decrypt-api/app/models"
	"encrypt-decrypt-api/app/server"
	"log"
	"sync"
)

func main() {
	// Konfiguration einlesen
	cfg := config.NewConfig()

	// Initialisiere die Datenbankverbindung
	dbConn, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Fehler bei der Datenbankinitialisierung: %v", err)
	}
	// Datenbankverbindung global verfügbar machen
	models.DB = dbConn

	// Starte beide Server parallel mit Goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Starte statischen Webserver auf Port 8081
	go func() {
		defer wg.Done()
		server.StartStaticFileServer(":8081")
	}()

	// Starte den API-Server auf Port 8080
	go func() {
		defer wg.Done()
		server.StartAPIServer(cfg, ":8080")
	}()

	// Warte, bis beide Goroutines fertig sind
	wg.Wait()
}
