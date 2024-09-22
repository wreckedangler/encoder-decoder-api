package handlers

import (
	"encrypt-decrypt-api/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListFilesHandler(c *gin.Context) {
	var files []models.File
	result := models.DB.Find(&files)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen der Daten"})
		return
	}
	c.JSON(http.StatusOK, files)
}
