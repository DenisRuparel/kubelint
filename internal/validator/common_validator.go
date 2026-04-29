package validators

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ValidateCommon(filePath string) []ValidationResult {
	var results []ValidationResult

	data, err := os.ReadFile(filePath)
	if err != nil {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "Failed to read file",
		})
		return results
	}

	var parsed map[string]interface{}

	err = yaml.Unmarshal(data, &parsed)
	if err != nil {
		return results
	}

	// CRITICAL → apiVersion
	if _, exists := parsed["apiVersion"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "apiVersion is missing",
		})
	}

	// CRITICAL → kind
	if _, exists := parsed["kind"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "kind is missing",
		})
	}

	// CRITICAL → metadata
	metaRaw, exists := parsed["metadata"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "metadata section is missing",
		})
		return results
	}

	metadata, ok := metaRaw.(map[string]interface{})
	if !ok {
		return results
	}

	// CRITICAL → metadata.name
	if _, exists := metadata["name"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "metadata.name is missing",
		})
	}

	// WARNING → namespace
	if _, exists := metadata["namespace"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "metadata.namespace is missing",
		})
	}

	// WARNING → labels
	if _, exists := metadata["labels"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "metadata.labels are missing",
		})
	}

	// INFO → annotations
	if _, exists := metadata["annotations"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Info,
			Message:  "metadata.annotations not set (recommended)",
		})
	}

	return results
}