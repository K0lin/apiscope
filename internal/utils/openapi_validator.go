package utils

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"strings"
)

type OpenAPIDocument struct {
	OpenAPI string `yaml:"openapi" json:"openapi"`
	Swagger string `yaml:"swagger" json:"swagger"`
	Info    struct {
		Title   string `yaml:"title" json:"title"`
		Version string `yaml:"version" json:"version"`
	} `yaml:"info" json:"info"`
}

func ValidateOpenAPIContent(content []byte) error {
	var doc OpenAPIDocument

	err := yaml.Unmarshal(content, &doc)
	if err != nil {
		err = json.Unmarshal(content, &doc)
		if err != nil {
			return errors.New("invalid YAML or JSON format")
		}
	}

	if doc.OpenAPI == "" && doc.Swagger == "" {
		return errors.New("not a valid OpenAPI/Swagger document")
	}

	if doc.OpenAPI != "" && !strings.HasPrefix(doc.OpenAPI, "3.") {
		return errors.New("only OpenAPI 3.x is supported")
	}

	if doc.Swagger != "" && !strings.HasPrefix(doc.Swagger, "2.") {
		return errors.New("only Swagger 2.x is supported")
	}

	if doc.Info.Title == "" || doc.Info.Version == "" {
		return errors.New("missing required fields: info.title or info.version")
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
