package handlers

import (
	"APIScope/internal/config"
	"APIScope/internal/models"
	"APIScope/internal/services"
	"APIScope/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	docService              *services.DocumentService
	storageService          *services.StorageService
	openAPIGeneratorService *services.OpenAPIGeneratorService
	cfg                     *config.Config
}

func NewApiHandler(docService *services.DocumentService, storageService *services.StorageService, openAPIGeneratorService *services.OpenAPIGeneratorService, cfg *config.Config) *ApiHandler {
	return &ApiHandler{
		docService:              docService,
		storageService:          storageService,
		openAPIGeneratorService: openAPIGeneratorService,
		cfg:                     cfg,
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

// DownloadDocumentVersion serves the raw stored file for a specific version (attachment)
func (h *ApiHandler) DownloadDocumentVersion(c *gin.Context) {
	documentID := c.Param("id")
	versionStr := c.Param("version")
	if versionStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version required"})
		return
	}

	doc, err := h.docService.GetDocumentByID(documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	var target *models.Version
	for i := range doc.Versions {
		if doc.Versions[i].Version == versionStr {
			target = &doc.Versions[i]
			break
		}
	}
	if target == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}
	content, err := h.storageService.GetFile(target.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read file"})
		return
	}
	// Force download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s-%s.yaml\"", documentID, versionStr))
	c.Data(http.StatusOK, "application/x-yaml", content)
}

// DeleteDocumentVersion deletes a single version (by human version string) if allowed.
func (h *ApiHandler) DeleteDocumentVersion(c *gin.Context) {
	// We don't have config reference here; route should be registered only if allowed.
	documentID := c.Param("id")
	version := c.Param("version")
	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version is required"})
		return
	}

	// Get version list to identify file to delete
	doc, err := h.docService.GetDocumentByID(documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	var filePath string
	for _, v := range doc.Versions {
		if v.Version == version {
			filePath = v.FilePath
			break
		}
	}
	if filePath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}

	if err := h.docService.DeleteVersion(documentID, version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Remove file from storage silently
	if filePath != "" {
		_ = os.Remove(filePath)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "version deleted"})
}

// SetShareLink assigns a one-time share slug to a document (if feature enabled).
// POST /api/document/:id/share  body: {"slug":"optional-custom"}
func (h *ApiHandler) SetShareLink(c *gin.Context) {
	if h.cfg == nil || !h.cfg.AllowCustomShareLink {
		c.JSON(http.StatusNotFound, gin.H{"error": "feature disabled"})
		return
	}
	documentID := c.Param("id")
	doc, err := h.docService.GetDocumentByID(documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	if doc.ShareSlug != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "share link already set", "share_slug": doc.ShareSlug})
		return
	}
	// Parse body (optional slug)
	var payload struct {
		Slug string `json:"slug"`
	}
	if c.Request.Body != nil {
		_ = json.NewDecoder(c.Request.Body).Decode(&payload) // ignore decode errors -> treat as empty
	}
	slug := payload.Slug
	if slug != "" {
		sanitized, ok := utils.SanitizeSlug(slug)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slug format (3-40 chars, a-z, 0-9, dashes)"})
			return
		}
		slug = sanitized
	} else {
		// generate random
		slug = utils.GenerateShareSlug()
	}
	if err := h.docService.SetShareSlug(doc, slug); err != nil {
		if err.Error() == "slug already taken" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fullURL := fmt.Sprintf("%s://%s/share/%s", func() string {
		if c.Request.TLS != nil {
			return "https"
		}
		return "http"
	}(), c.Request.Host, slug)
	c.JSON(http.StatusCreated, gin.H{"share_slug": slug, "url": fullURL})
}
