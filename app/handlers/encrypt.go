package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"encrypt-decrypt-api/app/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

func EncryptHandler(c *gin.Context) {
	// Read password and file from the request
	password := c.PostForm("password")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read the file"})
		return
	}
	defer file.Close()

	// Read the entire file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read the file data"})
		return
	}

	// Encrypt the file data
	encryptedData, err := encrypt(fileData, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failed"})
		return
	}

	// Compute hash of the encrypted file
	hash := sha256.Sum256(encryptedData)
	fileHash := hex.EncodeToString(hash[:])

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	// Get the original filename
	originalFilename := header.Filename

	// Remove the file extension(s)
	filenameWithoutExt := removeAllExtensions(originalFilename)

	// Create a new filename with .enc extension
	encryptedFilename := filenameWithoutExt + ".enc"

	// Save metadata to the database
	newFile := models.File{
		OriginalFilename: originalFilename,
		Filename:         encryptedFilename,
		Filesize:         int64(len(fileData)),
		UploadDate:       time.Now(),
		Password:         string(hashedPassword),
		FileHash:         fileHash,
	}
	result := models.DB.Create(&newFile)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database save failed"})
		return
	}

	// Return the encrypted file
	c.Header("Content-Disposition", "attachment; filename="+newFile.Filename)
	c.Data(http.StatusOK, "application/octet-stream", encryptedData)
}

// In your handlers package, add this function
func removeAllExtensions(filename string) string {
	filenameWithoutExt := filename
	for {
		ext := filepath.Ext(filenameWithoutExt)
		if ext == "" {
			break
		}
		filenameWithoutExt = strings.TrimSuffix(filenameWithoutExt, ext)
	}
	return filenameWithoutExt
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// Derive a key using PBKDF2
	key := pbkdf2.Key([]byte(passphrase), salt, 4096, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Use AES-GCM for encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Prepend the salt to the ciphertext
	encryptedData := append(salt, ciphertext...)

	return encryptedData, nil
}
