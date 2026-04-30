package validators

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ValidateNamespacePolicy checks for best practices in Namespace manifests

func ValidateNamespacePolicy(filePath string) []ValidationResult {
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

	metaRaw, exists := parsed["metadata"]
	if !exists {
		return results
	}

	metadata, ok := metaRaw.(map[string]interface{})
	if !ok {
		return results
	}

	// namespace check
	namespaceRaw, exists := metadata["namespace"]
	if !exists {
		return results
	}

	namespace, ok := namespaceRaw.(string)
	if !ok {
		return results
	}

	// CRITICAL → default namespace
	if namespace == "default" {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "using default namespace is not recommended for production",
		})
	}

	// labels check
	labelsRaw, exists := metadata["labels"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "namespace labels missing",
		})
		return results
	}

	labels, ok := labelsRaw.(map[string]interface{})
	if !ok {
		return results
	}

	if _, exists := labels["environment"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "environment label missing",
		})
	}

	if _, exists := labels["owner"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "owner/team label missing",
		})
	}

	if _, exists := labels["managed-by"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "managed-by label missing",
		})
	}

	return results
}

// ValidateNamespacePolicyBytes is a helper for validating raw YAML content

func ValidateNamespacePolicyBytes(content []byte) []ValidationResult {
	var results []ValidationResult

	var parsed map[string]interface{}

	metaRaw, exists := parsed["metadata"]
	if !exists {
		return results
	}

	metadata, ok := metaRaw.(map[string]interface{})
	if !ok {
		return results
	}

	// namespace check
	namespaceRaw, exists := metadata["namespace"]
	if !exists {
		return results
	}

	namespace, ok := namespaceRaw.(string)
	if !ok {
		return results
	}

	// CRITICAL → default namespace
	if namespace == "default" {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "using default namespace is not recommended for production",
		})
	}

	// labels check
	labelsRaw, exists := metadata["labels"]
	if !exists {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  "namespace labels missing",
		})
		return results
	}

	labels, ok := labelsRaw.(map[string]interface{})
	if !ok {
		return results
	}

	if _, exists := labels["environment"]; !exists {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  "environment label missing",
		})
	}

	if _, exists := labels["owner"]; !exists {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  "owner/team label missing",
		})
	}

	if _, exists := labels["managed-by"]; !exists {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  "managed-by label missing",
		})
	}

	return results
}