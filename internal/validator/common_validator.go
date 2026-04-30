package validators

import (
	"gopkg.in/yaml.v3"
	"os"
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

	var node yaml.Node
	err := yaml.Unmarshal(content, &node)
	if err != nil {
		return []ValidationResult{
			{
				Severity: Critical,
				Message:  "Invalid YAML structure",
			},
		}
	}

	if len(node.Content) == 0 {
		return results
	}

	doc := node.Content[0]

	var apiVersion, kind, name string
	var hasMetadata bool

	for i := 0; i < len(doc.Content); i += 2 {
		key := doc.Content[i].Value
		value := doc.Content[i+1]

		switch key {
		case "apiVersion":
			apiVersion = value.Value
		case "kind":
			kind = value.Value
		case "metadata":
			hasMetadata = true

			for j := 0; j < len(value.Content); j += 2 {
				if value.Content[j].Value == "name" {
					name = value.Content[j+1].Value
				}
			}
		}
	}

	// ✅ Correct checks
	if apiVersion == "" {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "apiVersion is missing or empty",
		})
	}

	if kind == "" {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "kind is missing or empty",
		})
	}

	if !hasMetadata {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "metadata section is missing",
		})
	}

	if name == "" {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "metadata.name is missing or empty",
		})
	}

	return results
}