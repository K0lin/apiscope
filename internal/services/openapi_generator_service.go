package services

import (
	"APIScope/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenAPIGeneratorService struct {
	config *config.Config
}

type GeneratorLanguage struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type GenerateRequest struct {
	OpenAPIUrl string            `json:"openAPIUrl"`
	Options    map[string]string `json:"options,omitempty"`
}

type GenerateResponse struct {
	Code string `json:"code"`
	Link string `json:"link"`
}

func NewOpenAPIGeneratorService(cfg *config.Config) *OpenAPIGeneratorService {
	return &OpenAPIGeneratorService{
		config: cfg,
	}
}

func (s *OpenAPIGeneratorService) IsEnabled() bool {
	return s.config.OpenAPIGeneratorEnabled
}

func (s *OpenAPIGeneratorService) GetAvailableLanguages() ([]GeneratorLanguage, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("OpenAPI generator is disabled")
	}

	url := fmt.Sprintf("%s/api/gen/clients", s.config.OpenAPIGeneratorServer)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch generators: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("generator service returned status %d", resp.StatusCode)
	}

	// Read the raw response first
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try to decode as array of objects first
	var languages []GeneratorLanguage
	if err := json.Unmarshal(body, &languages); err != nil {
		// If that fails, try to decode as array of strings
		var languageNames []string
		if err := json.Unmarshal(body, &languageNames); err != nil {
			return nil, fmt.Errorf("failed to decode response as either object array or string array: %w", err)
		}

		// Convert strings to GeneratorLanguage objects
		languages = make([]GeneratorLanguage, len(languageNames))
		for i, name := range languageNames {
			languages[i] = GeneratorLanguage{
				Name:        name,
				DisplayName: name,
				Description: fmt.Sprintf("%s client generator", name),
			}
		}
	}

	return languages, nil
}

func (s *OpenAPIGeneratorService) GenerateSDK(generator, openAPIUrl string, options map[string]string) (*GenerateResponse, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("OpenAPI generator is disabled")
	}

	generateReq := GenerateRequest{
		OpenAPIUrl: openAPIUrl,
		Options:    options,
	}

	jsonData, err := json.Marshal(generateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/gen/clients/%s", s.config.OpenAPIGeneratorServer, generator)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to generate SDK: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("generator service returned status %d: %s", resp.StatusCode, string(body))
	}

	var generateResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&generateResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &generateResp, nil
}

func (s *OpenAPIGeneratorService) DownloadSDK(downloadUrl string) ([]byte, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("OpenAPI generator is disabled")
	}

	resp, err := http.Get(downloadUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to download SDK: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download service returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read download data: %w", err)
	}

	return data, nil
}
