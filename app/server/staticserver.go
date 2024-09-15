package server

import (
	"log"
	"net/http"
)

// StartStaticFileServer startet den Webserver für statische Dateien auf dem angegebenen Port.
func StartStaticFileServer(port string) {
	// Statische Dateien bereitstellen (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("app/static"))
	http.Handle("/", fs)

	// Server starten
	log.Printf("Starte Webserver für statische Dateien auf %s...\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
