package validators

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ValidateConfigMap(filePath string) []ValidationResult {
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

	kind, ok := parsed["kind"].(string)
	if !ok || kind != "ConfigMap" {
		return results
	}

	// Check data section
	dataSection, hasData := parsed["data"]
	binaryDataSection, hasBinaryData := parsed["binaryData"]

	if !hasData && !hasBinaryData {
		results = append(results, ValidationResult{
			File:     filePath,
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
				File:     filePath,
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
				File:     filePath,
				Severity: Warning,
				Message:  "binaryData section exists but is empty",
			})
		}
	}

	return results
}