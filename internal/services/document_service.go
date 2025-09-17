package services

import (
	"APIScope/internal/database"
	"APIScope/internal/models"
	"APIScope/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DocumentService struct{}

func NewDocumentService() *DocumentService {
	return &DocumentService{}
}

func (s *DocumentService) CreateDocument(name, description string) (*models.Document, error) {
	doc := &models.Document{
		ID:          utils.GenerateDocumentID(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(time.Hour * 24 * 30),
		IsActive:    true,
		Versions:    []models.Version{},
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("document:%s", doc.ID)
	err = database.GetRedisClient().Set(database.GetContext(), key, docJSON, time.Until(doc.ExpiresAt)).Err()
	if err != nil {
		return nil, err
	}

	// Add to active documents set
	err = database.GetRedisClient().SAdd(database.GetContext(), "active_documents", doc.ID).Err()
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *DocumentService) GetDocumentByID(id string) (*models.Document, error) {
	key := fmt.Sprintf("document:%s", id)

	docJSON, err := database.GetRedisClient().Get(database.GetContext(), key).Result()
	if err != nil {
		return nil, errors.New("document not found or expired")
	}

	var doc models.Document
	err = json.Unmarshal([]byte(docJSON), &doc)
	if err != nil {
		return nil, err
	}

	// Check if document is still active and not expired
	if !doc.IsActive || time.Now().After(doc.ExpiresAt) {
		return nil, errors.New("document not found or expired")
	}

	// Load versions
	versions, err := s.getVersionsByDocumentID(id)
	if err == nil {
		doc.Versions = versions
	}

	return &doc, nil
}

func (s *DocumentService) DeleteDocument(id string) error {
	// Get document first to update it
	doc, err := s.GetDocumentByID(id)
	if err != nil {
		return err
	}

	doc.IsActive = false

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("document:%s", id)
	err = database.GetRedisClient().Set(database.GetContext(), key, docJSON, time.Until(doc.ExpiresAt)).Err()
	if err != nil {
		return err
	}

	// Remove from active documents set
	err = database.GetRedisClient().SRem(database.GetContext(), "active_documents", id).Err()
	return err
}

func (s *DocumentService) AddVersion(documentID string, filePath string, customVersion string) (*models.Version, error) {
	// Mark all existing versions as not latest
	versions, _ := s.getVersionsByDocumentID(documentID)
	for _, v := range versions {
		v.IsLatest = false
		s.saveVersion(&v)
	}

	if customVersion == "" {
		var existingVersions []string
		for _, v := range versions {
			existingVersions = append(existingVersions, v.Version)
		}
		customVersion = utils.GenerateVersionNumber(existingVersions)
	}

	newVersion := &models.Version{
		ID:         uuid.New().String(),
		DocumentID: documentID,
		Version:    customVersion,
		FilePath:   filePath,
		CreatedAt:  time.Now(),
		IsLatest:   true,
	}

	err := s.saveVersion(newVersion)
	if err != nil {
		return nil, err
	}

	return newVersion, nil
}

func (s *DocumentService) getVersionsByDocumentID(documentID string) ([]models.Version, error) {
	pattern := fmt.Sprintf("version:%s:*", documentID)
	keys, err := database.GetRedisClient().Keys(database.GetContext(), pattern).Result()
	if err != nil {
		return nil, err
	}

	var versions []models.Version
	for _, key := range keys {
		versionJSON, err := database.GetRedisClient().Get(database.GetContext(), key).Result()
		if err != nil {
			continue
		}

		var version models.Version
		err = json.Unmarshal([]byte(versionJSON), &version)
		if err != nil {
			continue
		}

		versions = append(versions, version)
	}

	return versions, nil
}

func (s *DocumentService) saveVersion(version *models.Version) error {
	versionJSON, err := json.Marshal(version)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("version:%s:%s", version.DocumentID, version.ID)
	return database.GetRedisClient().Set(database.GetContext(), key, versionJSON, 0).Err()
}