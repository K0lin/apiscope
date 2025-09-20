package utils

import (
	"encoding/json"
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

type OpenAPIDocument struct {
	OpenAPI string `yaml:"openapi" json:"openapi"`
	Swagger string `yaml:"swagger" json:"swagger"`
	Info    struct {
		Title   string `yaml:"title" json:"title"`
		Version string `yaml:"version" json:"version"`
	} `yaml:"info" json:"info"`
	Paths      map[string]any `yaml:"paths" json:"paths"`
	Components any            `yaml:"components" json:"components"`
}

func ValidateOpenAPIContent(content []byte) error {
	trimmed := strings.TrimSpace(string(content))
	if trimmed == "" {
		return errors.New("document is empty")
	}

	var doc OpenAPIDocument

	// Try YAML first then JSON (works either way due to superset nature)
	if err := yaml.Unmarshal(content, &doc); err != nil {
		if err2 := json.Unmarshal(content, &doc); err2 != nil {
			return errors.New("invalid YAML or JSON format")
		}
	}

	if doc.OpenAPI == "" && doc.Swagger == "" {
		return errors.New("missing 'openapi' or 'swagger' version field")
	}

	if doc.OpenAPI != "" && !strings.HasPrefix(doc.OpenAPI, "3.") {
		return errors.New("only OpenAPI 3.x is supported")
	}
	if doc.Swagger != "" && !strings.HasPrefix(doc.Swagger, "2.") {
		return errors.New("only Swagger 2.x is supported")
	}

	if strings.TrimSpace(doc.Info.Title) == "" || strings.TrimSpace(doc.Info.Version) == "" {
		return errors.New("missing required fields: info.title or info.version")
	}

	// Minimal structural check: either at least one path or components defined
	if len(doc.Paths) == 0 && doc.Components == nil {
		return errors.New("document missing both 'paths' and 'components' sections")
	}

	return nil
}

func GetDocumentInfo(content []byte) (string, string, error) {
	var doc OpenAPIDocument

	err := yaml.Unmarshal(content, &doc)
	if err != nil {
		err = json.Unmarshal(content, &doc)
		if err != nil {
			return "", "", err
		}
	}

	return doc.Info.Title, doc.Info.Version, nil
}
