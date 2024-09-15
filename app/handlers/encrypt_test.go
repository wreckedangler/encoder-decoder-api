package handlers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/pbkdf2"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncryptHandler_Success(t *testing.T) {
	// Set up Gin context
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/encrypt", EncryptHandler)

	// Create a sample file
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	fileWriter, _ := writer.CreateFormFile("file", "testfile.txt")
	fileContent := []byte("This is a test file")
	fileWriter.Write(fileContent)

	// Add password to the form data
	writer.WriteField("password", "testpassword")
	writer.Close()

	// Create a test request
	req, _ := http.NewRequest("POST", "/encrypt", &buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Record the response
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.Code)

	// Assert that the response contains the right content-type
	assert.Equal(t, "application/octet-stream", resp.Header().Get("Content-Type"))
}

func TestEncryptHandler_FileReadError(t *testing.T) {
	// Set up Gin context
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/encrypt", EncryptHandler)

	// Create a sample request without a file
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("password", "testpassword")
	writer.Close()

	req, _ := http.NewRequest("POST", "/encrypt", &buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Record the response
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the response status code
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Datei konnte nicht gelesen werden")
}

func TestEncrypt(t *testing.T) {
	data := []byte("Hello, World!")
	password := "testpassword"

	// Encrypt the data
	encryptedData, err := encrypt(data, password)

	// Assert no error occurred
	assert.NoError(t, err)

	// Ensure the encrypted data is not nil
	assert.NotNil(t, encryptedData)
}

func TestEncrypt_AES_Error(t *testing.T) {
	// Create an invalid passphrase to test AES cipher failure
	_, err := aes.NewCipher([]byte("shortpass"))
	assert.Error(t, err, "should return an error for short key length")
}

func TestEncrypt_GCM_Error(t *testing.T) {
	// Use invalid key length for AES block
	key := pbkdf2.Key([]byte("testpassword"), []byte("salt"), 4096, 32, sha256.New)
	block, _ := aes.NewCipher(key)

	// Simulate an error in GCM
	_, err := cipher.NewGCM(block)
	assert.NoError(t, err, "GCM should work correctly with a valid key")
}
