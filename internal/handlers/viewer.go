package handlers

import (
	"APIScope/internal/config"
	"APIScope/internal/models"
	"APIScope/internal/services"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type ViewerHandler struct {
	docService     *services.DocumentService
	storageService *services.StorageService
	config         *config.Config
}

func NewViewerHandler(docService *services.DocumentService, storageService *services.StorageService, cfg *config.Config) *ViewerHandler {
	return &ViewerHandler{
		docService:     docService,
		storageService: storageService,
		config:         cfg,
	}
}

func (h *ViewerHandler) ViewDocument(c *gin.Context) {
	documentID := c.Param("id")
	selectedVersion := c.Query("version")
	message := c.Query("message")
	messageType := c.DefaultQuery("type", "info")

	fmt.Printf("ViewDocument called for ID: %s, version: %s\n", documentID, selectedVersion)

	doc, err := h.docService.GetDocumentByID(documentID)
	if err != nil {
		fmt.Printf("Error getting document: %v\n", err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Document not found or expired",
			"title": "Document Not Found",
		})
		return
	}

	fmt.Printf("Document found: %s, Versions count: %d\n", doc.Name, len(doc.Versions))

	var targetVersion *models.Version
	var content string

	if selectedVersion != "" {
		// Find specific version
		for _, version := range doc.Versions {
			if version.Version == selectedVersion {
				targetVersion = &version
				break
			}
		}
		if targetVersion == nil {
			// Version not found, use latest
			for _, version := range doc.Versions {
				if version.IsLatest {
					targetVersion = &version
					selectedVersion = version.Version // Update to actual version being shown
					break
				}
			}
			message = "Version '" + selectedVersion + "' not found, showing latest version"
			messageType = "info"
		}
	} else {
		// Use latest version
		for _, version := range doc.Versions {
			if version.IsLatest {
				targetVersion = &version
				selectedVersion = version.Version // Set selectedVersion to latest
				break
			}
		}
	}

	if targetVersion == nil {
		fmt.Printf("No versions found\n")
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "No versions found for this document",
			"title": "No Versions Found",
		})
		return
	}

	fmt.Printf("Target version: %s, FilePath: %s\n", targetVersion.Version, targetVersion.FilePath)

	contentBytes, err := h.storageService.GetFile(targetVersion.FilePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Error reading document content: " + err.Error(),
			"title": "Error",
		})
		return
	}
	content = string(contentBytes)

	fmt.Printf("File content loaded successfully, length: %d\n", len(content))
	if len(content) > 0 {
		previewLen := 200
		if len(content) < previewLen {
			previewLen = len(content)
		}
		fmt.Printf("Content preview: %s\n", content[:previewLen])
	}

	// Sort versions by CreatedAt desc so newest first (latest should normally be first)
	sort.Slice(doc.Versions, func(i, j int) bool {
		return doc.Versions[i].CreatedAt.After(doc.Versions[j].CreatedAt)
	})

	// Prepare versions for template while marking which one is selected
	var versions []gin.H
	for _, version := range doc.Versions {
		versions = append(versions, gin.H{
			"ID":        version.ID,
			"Version":   version.Version,
			"CreatedAt": version.CreatedAt,
			"IsLatest":  version.IsLatest,
			"Selected":  version.Version == selectedVersion,
		})
	}

	templateData := gin.H{
		"Title":                   doc.Name,
		"Document":                doc,
		"Version":                 targetVersion,
		"Content":                 content,
		"DocumentID":              documentID,
		"SelectedVersion":         selectedVersion,
		"Versions":                versions,
		"Message":                 message,
		"MessageType":             messageType,
		"OpenAPIGeneratorEnabled": h.config.OpenAPIGeneratorEnabled,
		"OpenAPIGeneratorServer":  h.config.OpenAPIGeneratorServer,
		"AllowVersionDeletion":    h.config.AllowVersionDeletion,
		"AllowVersionDownload":    h.config.AllowVersionDownload,
		"AllowServerEditing":      h.config.AllowServerEditing,
		"AutoAdjustServerOrigin":  h.config.AutoAdjustServerOrigin,
		"StripServers":            h.config.StripServers,
	}

	c.HTML(http.StatusOK, "viewer.html", templateData)
}

func (h *ViewerHandler) DeleteDocument(c *gin.Context) {
	documentID := c.Param("id")

	err := h.docService.DeleteDocument(documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error deleting document",
		})
		return
	}

	h.storageService.DeleteDocument(documentID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Document deleted successfully",
	})
}
