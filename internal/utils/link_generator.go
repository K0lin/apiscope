package utils

import (
    "fmt"
    "github.com/google/uuid"
    "strings"
)

func GenerateDocumentID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func GenerateVersionNumber(existingVersions []string) string {
	maxVersion := 0

	for _, version := range existingVersions {
		if strings.HasPrefix(version, "v") {
			var num int
			if n, err := fmt.Sscanf(version, "v%d", &num); n == 1 && err == nil {
				if num > maxVersion {
					maxVersion = num
				}
			}
		}
	}

	return fmt.Sprintf("v%d", maxVersion+1)
}
