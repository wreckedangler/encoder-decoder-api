package router

import (
	"encrypt-decrypt-api/app/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(allowedOrigins string) *gin.Engine {
	router := gin.Default()

	// CORS-Konfiguration hinzufügen
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigins},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// Neue Routen definieren
	router.POST("/encrypt", handlers.EncryptHandler)
	router.POST("/decrypt", handlers.DecryptHandler)
	router.GET("/files", handlers.ListFilesHandler)

	// Route für Root-URL
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Willkommen zur Verschlüsselungs-API!")
	})

	return router
}
