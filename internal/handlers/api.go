package handlers

import (
	"APIScope/internal/models"
	"APIScope/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	docService              *services.DocumentService
	storageService          *services.StorageService
	openAPIGeneratorService *services.OpenAPIGeneratorService
}

func NewApiHandler(docService *services.DocumentService, storageService *services.StorageService, openAPIGeneratorService *services.OpenAPIGeneratorService) *ApiHandler {
	return &ApiHandler{
		docService:              docService,
		storageService:          storageService,
		openAPIGeneratorService: openAPIGeneratorService,
	}
}

func (h *ApiHandler) GetDocumentContent(c *gin.Context) {
	documentID := c.Param("id")
	requestedVersion := c.Query("version") // Get version from query parameter

	// Get the document
	doc, err := h.docService.GetDocumentByID(documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Document not found or expired",
		})
		return
	}

	var targetVersion *models.Version

	if requestedVersion != "" {
		// Look for specific version
		for _, version := range doc.Versions {
			if version.Version == requestedVersion {
				targetVersion = &version
				break
			}
		}

		if targetVersion == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Version not found: " + requestedVersion,
			})
			return
		}
	} else {
		// Get latest version if no specific version requested
		for _, version := range doc.Versions {
			if version.IsLatest {
				targetVersion = &version
				break
			}
		}
	}

	if targetVersion == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No versions found",
		})
		return
	}

	content, err := h.storageService.GetFile(targetVersion.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error reading file",
		})
		return
	}

	c.Header("Content-Type", "application/yaml")
	c.String(http.StatusOK, string(content))
}

func (h *ApiHandler) GetDocumentVersions(c *gin.Context) {
	documentID := c.Param("id")

	doc, err := h.docService.GetDocumentByID(documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Document not found or expired",
		})
		return
	}

	var versions []gin.H
	for _, version := range doc.Versions {
		versions = append(versions, gin.H{
			"id":         version.ID,
			"version":    version.Version,
			"created_at": version.CreatedAt,
			"is_latest":  version.IsLatest,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"document_id": doc.ID,
		"versions":    versions,
	})
}

func (h *ApiHandler) GetAvailableLanguages(c *gin.Context) {
	if !h.openAPIGeneratorService.IsEnabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "OpenAPI Generator service is disabled",
		})
		return
	}

	languages, err := h.openAPIGeneratorService.GetAvailableLanguages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch available languages: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"languages": languages,
	})
}

func (h *ApiHandler) GenerateSDK(c *gin.Context) {
	if !h.openAPIGeneratorService.IsEnabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "OpenAPI Generator service is disabled",
		})
		return
	}

	documentID := c.Param("id")
	generator := c.Param("generator")
	requestedVersion := c.Query("version")

	// Get the document content URL
	baseURL := fmt.Sprintf("%s://%s", func() string {
		if c.Request.TLS != nil {
			return "https"
		}
		return "http"
	}(), c.Request.Host)

	contentURL := fmt.Sprintf("%s/api/document/%s/content", baseURL, documentID)
	if requestedVersion != "" {
		contentURL += "?version=" + requestedVersion
	}

	// Generate SDK
	result, err := h.openAPIGeneratorService.GenerateSDK(generator, contentURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate SDK: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ApiHandler) DownloadSDK(c *gin.Context) {
	if !h.openAPIGeneratorService.IsEnabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "OpenAPI Generator service is disabled",
		})
		return
	}

	downloadURL := c.Query("url")
	if downloadURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Download URL is required",
		})
		return
	}

	data, err := h.openAPIGeneratorService.DownloadSDK(downloadURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to download SDK: " + err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=sdk.zip")
	c.Data(http.StatusOK, "application/zip", data)
}
