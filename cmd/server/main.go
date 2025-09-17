package main

import (
	"APIScope/internal/config"
	"APIScope/internal/database"
	"APIScope/internal/handlers"
	"APIScope/internal/services"
	"fmt"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	fmt.Println("Configuration loaded successfully!")

	err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Database error:", err)
	}
	fmt.Println("Database initialized successfully!")

	docService := services.NewDocumentService()
	storageService := services.NewStorageService(cfg)

	uploadHandler := handlers.NewUploadHandler(docService, storageService, cfg)
	viewerHandler := handlers.NewViewerHandler(docService, storageService)
	apiHandler := handlers.NewApiHandler(docService, storageService)

	router := gin.Default()

	// Set max multipart form size (64MB)
	router.MaxMultipartMemory = 64 << 20

	// Get absolute path for templates
	templatePath, _ := filepath.Abs("web/templates/*")
	router.LoadHTMLGlob(templatePath)
	fmt.Println("HTML templates loaded from:", templatePath)

	// Get absolute path for static files
	staticPath, _ := filepath.Abs("./web/static")
	router.Static("/static", staticPath)
	fmt.Println("Static files served from:", staticPath)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/upload")
	})

	router.GET("/upload", uploadHandler.ShowUploadPage)
	router.POST("/upload", uploadHandler.HandleUpload)

	router.GET("/view/:id", viewerHandler.ViewDocument)
	router.DELETE("/view/:id", viewerHandler.DeleteDocument)

	router.GET("/api/document/:id/content", apiHandler.GetDocumentContent)
	router.GET("/api/document/:id/versions", apiHandler.GetDocumentVersions)

	fmt.Printf("Server starting on port %s...\n", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}
