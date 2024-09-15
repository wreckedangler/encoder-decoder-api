package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// APIClient enthält die Konfigurationsinformationen für den API-Client
type APIClient struct {
	BaseURL string
	Client  *http.Client
}

// NewAPIClient erstellt eine neue Instanz des API-Clients
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 30 * time.Second}, // Setze einen Timeout
	}
}

// EncryptFile verschlüsselt eine Datei über die API
func (api *APIClient) EncryptFile(filePath string, password string) ([]byte, error) {
	url := fmt.Sprintf("%s/encrypt", api.BaseURL)
	return api.uploadFile(url, filePath, password)
}

// DecryptFile entschlüsselt eine Datei über die API
func (api *APIClient) DecryptFile(filePath string, password string) ([]byte, error) {
	url := fmt.Sprintf("%s/decrypt", api.BaseURL)
	return api.uploadFile(url, filePath, password)
}

// uploadFile sendet eine Datei und ein Passwort an die angegebene URL
func (api *APIClient) uploadFile(url string, filePath string, password string) ([]byte, error) {
	// Datei- und Passwortvalidierung
	if filePath == "" {
		return nil, fmt.Errorf("Dateipfad darf nicht leer sein")
	}
	if password == "" {
		return nil, fmt.Errorf("Passwort darf nicht leer sein")
	}

	// Datei öffnen
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("konnte Datei nicht öffnen: %v", err)
	}
	defer file.Close()

	// Einen neuen multipart Writer erstellen
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Passwortfeld hinzufügen
	err = writer.WriteField("password", password)
	if err != nil {
		return nil, fmt.Errorf("konnte Passwortfeld nicht erstellen: %v", err)
	}

	// Dateifeld hinzufügen
	fileField, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("konnte Dateifeld nicht erstellen: %v", err)
	}

	// Dateiinhalt kopieren
	_, err = io.Copy(fileField, file)
	if err != nil {
		return nil, fmt.Errorf("konnte Datei nicht lesen: %v", err)
	}

	// Writer schließen, um die Grenze zu finalisieren
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("konnte Writer nicht schließen: %v", err)
	}

	// HTTP-Anfrage erstellen
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("konnte HTTP-Anfrage nicht erstellen: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Anfrage senden
	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP-Anfrage fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()

	// Statuscode prüfen
	if resp.StatusCode != http.StatusOK {
		// Fehlernachricht lesen
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Fehler beim Lesen der Fehlernachricht: %v", err)
		}
		return nil, fmt.Errorf("Fehler vom Server: %s", string(bodyBytes))
	}

	// Antwortdaten lesen
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("konnte Antwort nicht lesen: %v", err)
	}

	return responseData, nil
}
