package validators

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ValidateConfigMap checks for common best practices in ConfigMap manifests
func ValidateConfigMap(filePath string) []ValidationResult {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []ValidationResult{
			{
				File:     filePath,
				Severity: Critical,
				Message:  "Failed to read file: " + err.Error(),
			},
		}
	}

	results := ValidateConfigMapBytes(data)
	for i := range results {
		results[i].File = filePath
	}
	return results
}

// ValidateConfigMapBytes operates on raw YAML content instead of file paths.
func ValidateConfigMapBytes(content []byte) []ValidationResult {
	var results []ValidationResult
	var parsed map[string]interface{}

	if err := yaml.Unmarshal(content, &parsed); err != nil {
		return results
	}

	kind, ok := parsed["kind"].(string)
	if !ok || kind != "ConfigMap" {
		return results
	}

	// Check data section
	dataSection, hasData := parsed["data"]
	binaryDataSection, hasBinaryData := parsed["binaryData"]

	if !hasData && !hasBinaryData {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "ConfigMap must contain data or binaryData",
		})
		return results
	}

	// Empty data check
	if hasData {
		dataMap, ok := dataSection.(map[string]interface{})
		if ok && len(dataMap) == 0 {
			results = append(results, ValidationResult{
				Severity: Warning,
				Message:  "data section exists but is empty",
			})
		}
	}

	// Empty binaryData check
	if hasBinaryData {
		binaryMap, ok := binaryDataSection.(map[string]interface{})
		if ok && len(binaryMap) == 0 {
			results = append(results, ValidationResult{
				Severity: Warning,
				Message:  "binaryData section exists but is empty",
			})
		}
	}

	return results
}