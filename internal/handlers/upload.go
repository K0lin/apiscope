package handlers

import (
	"APIScope/internal/config"
	"APIScope/internal/models"
	"APIScope/internal/services"
	"APIScope/internal/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	docService     *services.DocumentService
	storageService *services.StorageService
	config         *config.Config
}

func NewUploadHandler(docService *services.DocumentService, storageService *services.StorageService, cfg *config.Config) *UploadHandler {
	return &UploadHandler{
		docService:     docService,
		storageService: storageService,
		config:         cfg,
	}
}

func (h *UploadHandler) ShowUploadPage(c *gin.Context) {
	message := c.Query("message")
	messageType := c.DefaultQuery("type", "info")

	c.HTML(http.StatusOK, "upload.html", gin.H{
		"title":       "Upload OpenAPI Document",
		"Message":     message,
		"MessageType": messageType,
	})
}

func (h *UploadHandler) HandleUpload(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")
	customVersion := c.PostForm("version")
	yamlContent := c.PostForm("yaml_content")
	documentID := c.PostForm("document_id") // Check if adding to existing document

	var content []byte
	var err error

	// Debug log
	fmt.Printf("Upload request - Name: %s, Description: %s, Version: %s, Document ID: %s, YAML Content length: %d\n",
		name, description, customVersion, documentID, len(yamlContent))

	// Check if text content was provided
	if yamlContent != "" && len(strings.TrimSpace(yamlContent)) > 0 {
		content = []byte(yamlContent)
		fmt.Printf("Using pasted content, length: %d\n", len(content))
	} else {
		// Try to get uploaded file
		file, err := c.FormFile("file")
		if err != nil {
			fmt.Printf("File upload error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Please provide either a file or YAML content. Error: " + err.Error(),
				"success": false,
			})
			return
		}

		fmt.Printf("File uploaded - Name: %s, Size: %d bytes\n", file.Filename, file.Size)

		if file.Size > h.config.MaxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "File too large. Maximum size: " + strconv.FormatInt(h.config.MaxFileSize/(1024*1024), 10) + "MB",
				"success": false,
			})
			return
		}

		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error reading file",
				"success": false,
			})
			return
		}
		defer fileContent.Close()

		content = make([]byte, file.Size)
		_, err = fileContent.Read(content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error reading file content",
				"success": false,
			})
			return
		}
	}

	err = utils.ValidateOpenAPIContent(content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid OpenAPI document: " + err.Error(),
			"success": false,
		})
		return
	}

	var doc *models.Document

	if documentID != "" {
		// Adding version to existing document
		fmt.Printf("Adding new version to existing document: %s\n", documentID)
		doc, err = h.docService.GetDocumentByID(documentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Document not found: " + err.Error(),
				"success": false,
			})
			return
		}
	} else {
		// Creating new document
		if name == "" {
			title, _, err := utils.GetDocumentInfo(content)
			if err == nil {
				name = title
			} else {
				name = "Untitled API"
			}
		}

		doc, err = h.docService.CreateDocument(name, description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error creating document: " + err.Error(),
				"success": false,
			})
			return
		}
	}

	// Generate version identifier for file storage
	versionID := customVersion
	if versionID == "" {
		versionID = fmt.Sprintf("v%d", len(doc.Versions)+1)
	}

	filePath, err := h.storageService.SaveFile(doc.ID, versionID, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error saving file: " + err.Error(),
			"success": false,
		})
		return
	}

	version, err := h.docService.AddVersion(doc.ID, filePath, customVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating version: " + err.Error(),
			"success": false,
		})
		return
	}

	// Return JSON with document ID and success info
	if c.GetHeader("Accept") == "application/json" || c.Query("ajax") == "1" {
		c.JSON(http.StatusCreated, gin.H{
			"success":     true,
			"document_id": doc.ID,
			"version_id":  version.ID,
			"message":     "Document uploaded successfully",
			"view_url":    "/view/" + doc.ID,
		})
	} else {
		// Redirect to viewer page for regular form submissions
		c.Redirect(http.StatusFound, "/view/"+doc.ID)
	}
}
