package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encrypt-decrypt-api/app/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

func DecryptHandler(c *gin.Context) {
	// Passwort und Datei aus der Anfrage lesen
	password := c.PostForm("password")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datei konnte nicht gelesen werden"})
		return
	}
	defer file.Close()

	// Dateiinhalt lesen
	fileData := make([]byte, header.Size)
	_, err = file.Read(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Datei konnte nicht gelesen werden"})
		return
	}

	// Metadaten aus der Datenbank abrufen
	var dbFile models.File
	result := models.DB.Where("filename = ?", header.Filename).First(&dbFile)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Datei nicht gefunden"})
		return
	}

	// Passwort überprüfen
	err = bcrypt.CompareHashAndPassword([]byte(dbFile.Password), []byte(password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Falsches Passwort"})
		return
	}

	// Datei entschlüsseln
	decryptedData, err := decrypt(fileData, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Entschlüsselung fehlgeschlagen"})
		return
	}

	// Entschlüsselte Datei zurückgeben
	originalFilename := header.Filename
	if len(originalFilename) > 4 && originalFilename[len(originalFilename)-4:] == ".enc" {
		originalFilename = originalFilename[:len(originalFilename)-4]
	}
	c.Header("Content-Disposition", "attachment; filename="+originalFilename)
	c.Data(http.StatusOK, "application/octet-stream", decryptedData)
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	// Salt und Ciphertext trennen
	if len(data) < 16 {
		return nil, errors.New("Daten sind zu kurz")
	}
	salt := data[:16]
	ciphertext := data[16:]

	// Schlüssel ableiten
	key := pbkdf2.Key([]byte(passphrase), salt, 4096, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("Ciphertext ist zu kurz")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Daten entschlüsseln
	decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}
