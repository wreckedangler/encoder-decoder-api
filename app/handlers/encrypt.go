package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encrypt-decrypt-api/app/models"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

func EncryptHandler(c *gin.Context) {
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

	// Datei verschlüsseln
	encryptedData, err := encrypt(fileData, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Verschlüsselung fehlgeschlagen"})
		return
	}

	// Passwort hashen
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Passwort-Hashing fehlgeschlagen"})
		return
	}

	// Metadaten in der Datenbank speichern
	newFile := models.File{
		Filename:   header.Filename + ".enc",
		Filesize:   header.Size,
		UploadDate: time.Now(),
		Password:   string(hashedPassword),
	}
	result := models.DB.Create(&newFile)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Datenbankspeicherung fehlgeschlagen"})
		return
	}

	// Verschlüsselte Datei zurückgeben
	c.Header("Content-Disposition", "attachment; filename="+header.Filename+".enc")
	c.Data(http.StatusOK, "application/octet-stream", encryptedData)
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	// Salt generieren
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

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

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Daten verschlüsseln
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Salt und Ciphertext kombinieren
	encryptedData := append(salt, ciphertext...)

	return encryptedData, nil
}
