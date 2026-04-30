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

	var root yaml.Node
	if err := yaml.Unmarshal(content, &root); err != nil {
		return []ValidationResult{
			{
				Severity: Critical,
				Message:  "Invalid YAML",
			},
		}
	}

	if len(root.Content) == 0 {
		return results
	}

	doc := root.Content[0] // actual document

	var apiVersion, kind, name string
	var hasMetadata bool

	for i := 0; i < len(doc.Content); i += 2 {
		key := doc.Content[i].Value
		val := doc.Content[i+1]

		switch key {

		case "apiVersion":
			apiVersion = val.Value

		case "kind":
			kind = val.Value

		case "metadata":
			hasMetadata = true

			for j := 0; j < len(val.Content); j += 2 {
				if val.Content[j].Value == "name" {
					name = val.Content[j+1].Value
				}
			}
		}
	}

	// 🔥 THESE WILL NOW WORK CORRECTLY
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