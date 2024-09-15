package server

import (
	"encrypt-decrypt-api/app/config"
	"encrypt-decrypt-api/app/router"
	"log"
)

// StartAPIServer startet den API-Server auf dem angegebenen Port.
func StartAPIServer(cfg *config.Config, port string) {
	// Router aus dem Router-Paket erhalten
	r := router.NewRouter(cfg.AllowedOrigins)

	// API-Server starten
	if err := r.Run(port); err != nil {
		log.Fatalf("API-Server konnte nicht gestartet werden: %v", err)
	}
}
