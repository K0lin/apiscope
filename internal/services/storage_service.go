package services

import (
	"APIScope/internal/config"
	"fmt"
	"os"
	"path/filepath"
)

type StorageService struct {
	storagePath string
}

func NewStorageService(cfg *config.Config) *StorageService {
	return &StorageService{
		storagePath: cfg.StoragePath,
	}
}

func (s *StorageService) SaveFile(documentID, version string, content []byte) (string, error) {
	docDir := filepath.Join(s.storagePath, documentID)
	err := os.MkdirAll(docDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filename := fmt.Sprintf("%s.yaml", version)
	filePath := filepath.Join(docDir, filename)

	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filePath, nil
}

func (s *StorageService) GetFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return content, nil
}

func (s *StorageService) DeleteDocument(documentID string) error {
	docDir := filepath.Join(s.storagePath, documentID)
	return os.RemoveAll(docDir)
}
