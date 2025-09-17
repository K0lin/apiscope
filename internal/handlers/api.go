package handlers

import (
	"APIScope/internal/models"
	"APIScope/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	docService     *services.DocumentService
	storageService *services.StorageService
}

func NewApiHandler(docService *services.DocumentService, storageService *services.StorageService) *ApiHandler {
	return &ApiHandler{
		docService:     docService,
		storageService: storageService,
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
