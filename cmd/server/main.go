package main

import (
	"APIScope/internal/config"
	"APIScope/internal/database"
	"APIScope/internal/handlers"
	"APIScope/internal/services"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

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
	openAPIGeneratorService := services.NewOpenAPIGeneratorService(cfg)

	uploadHandler := handlers.NewUploadHandler(docService, storageService, cfg)
	viewerHandler := handlers.NewViewerHandler(docService, storageService, cfg)
	apiHandler := handlers.NewApiHandler(docService, storageService, openAPIGeneratorService)

	router := gin.Default()

	// Enhanced CORS middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		reqOrigin := c.GetHeader("Origin")
		method := c.Request.Method

		// Decide allowed origin
		allowOrigin := ""
		if len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "*" {
			allowOrigin = "*"
		} else {
			for _, o := range cfg.AllowedOrigins {
				if strings.EqualFold(o, reqOrigin) {
					allowOrigin = reqOrigin
					break
				}
			}
		}
		// If credentials requested and wildcard configured, echo specific origin instead of *
		if cfg.CORSAllowCredentials && allowOrigin == "*" && reqOrigin != "" {
			allowOrigin = reqOrigin
		}

		if allowOrigin != "" || reqOrigin == "" { // allow same-origin/no-origin cases to pass with wildcard config
			if allowOrigin == "" && reqOrigin == "" && len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "*" {
				allowOrigin = "*"
			}
			if allowOrigin != "" {
				c.Header("Access-Control-Allow-Origin", allowOrigin)
			}
			if cfg.CORSAllowCredentials && allowOrigin != "" && allowOrigin != "*" {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			// Allow headers (merge requested ones)
			allowHeadersSet := map[string]struct{}{}
			for _, h := range cfg.CORSAllowedHeaders {
				allowHeadersSet[strings.TrimSpace(h)] = struct{}{}
			}
			if req := c.GetHeader("Access-Control-Request-Headers"); req != "" {
				for _, h := range strings.Split(req, ",") {
					allowHeadersSet[strings.TrimSpace(h)] = struct{}{}
				}
			}
			var allowHeadersList []string
			for h := range allowHeadersSet {
				if h != "" {
					allowHeadersList = append(allowHeadersList, h)
				}
			}
			c.Header("Access-Control-Allow-Headers", strings.Join(allowHeadersList, ", "))
			c.Header("Access-Control-Allow-Methods", strings.Join(cfg.CORSAllowedMethods, ", "))
			if len(cfg.CORSExposeHeaders) > 0 {
				c.Header("Access-Control-Expose-Headers", strings.Join(cfg.CORSExposeHeaders, ", "))
			}
			if cfg.CORSMaxAgeSeconds > 0 {
				c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", cfg.CORSMaxAgeSeconds))
			}
		}

		if method == http.MethodOptions {
			c.AbortWithStatus(204)
			if cfg.CORSDebug {
				log.Printf("CORS preflight %s origin=%s allowed=%s duration=%s", c.FullPath(), reqOrigin, allowOrigin, time.Since(start))
			}
			return
		}
		c.Next()
		if cfg.CORSDebug {
			log.Printf("CORS request %s %s origin=%s status=%d allowed=%s duration=%s", method, c.FullPath(), reqOrigin, c.Writer.Status(), allowOrigin, time.Since(start))
		}
	})

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

	// Basic health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Version deletion endpoint (conditional)
	if cfg.AllowVersionDeletion {
		router.DELETE("/api/document/:id/version/:version", apiHandler.DeleteDocumentVersion)
	}
	if cfg.AllowVersionDownload {
		router.GET("/api/document/:id/version/:version/download", apiHandler.DownloadDocumentVersion)
	}

	fmt.Printf("Server starting on port %s...\n", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}
