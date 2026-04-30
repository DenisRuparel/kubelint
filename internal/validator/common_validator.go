package validators

import (
	"os"
	"fmt"
	"gopkg.in/yaml.v3"
)

// ValidateCommon checks for common fields in Kubernetes manifests

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

// ValidateCommonBytes is a helper for validating YAML content directly

func ValidateCommonBytes(content []byte) []ValidationResult {
	var results []ValidationResult

	var parsed map[string]interface{}

	// CRITICAL → apiVersion
	apiVersion, exists := parsed["apiVersion"]
	if !exists || apiVersion == nil || fmt.Sprintf("%v", apiVersion) == "" {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "apiVersion is missing or empty",
		})
	}

	// CRITICAL → kind
	kind, exists := parsed["kind"]
	if !exists || kind == nil || fmt.Sprintf("%v", kind) == "" {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "kind is missing or empty",
		})
	}

	// CRITICAL → metadata
	metaRaw, exists := parsed["metadata"]
	if !exists || metaRaw == nil || fmt.Sprintf("%v", metaRaw) == "" {
		results = append(results, ValidationResult{
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
	name, exists := metadata["name"]
	if !exists || name == nil || fmt.Sprintf("%v", name) == "" {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "metadata.name is missing or empty",
		})
	}

	// WARNING → namespace
	if _, exists := metadata["namespace"]; !exists {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  "metadata.namespace is missing",
		})
	}

	// WARNING → labels
	if _, exists := metadata["labels"]; !exists {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  "metadata.labels are missing",
		})
	}

	// INFO → annotations
	if _, exists := metadata["annotations"]; !exists {
		results = append(results, ValidationResult{
			Severity: Info,
			Message:  "metadata.annotations not set (recommended)",
		})
	}

	return results
}
