package validators

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ValidateSecret checks for common best practices in Secret manifests

func ValidateSecret(filePath string) []ValidationResult {
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
	if !ok || kind != "Secret" {
		return results
	}

	// type validation
	secretType, exists := parsed["type"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "Secret type is missing",
		})
	} else {
		if secretType == "Opaque" {
			results = append(results, ValidationResult{
				File:     filePath,
				Severity: Info,
				Message:  "Secret type is Opaque",
			})
		}
	}

	// data or stringData required
	dataSection, hasData := parsed["data"]
	stringDataSection, hasStringData := parsed["stringData"]

	if !hasData && !hasStringData {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "Secret must contain data or stringData",
		})
		return results
	}

	// Validate stringData placeholders
	if hasStringData {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Info,
			Message:  "Using stringData for Secret (acceptable for development)",
		})

		stringMap, ok := stringDataSection.(map[string]interface{})
		if ok {
			for key, value := range stringMap {
				val := strings.ToLower(value.(string))

				if strings.Contains(val, "your-api-key") ||
					strings.Contains(val, "password-here") ||
					strings.Contains(val, "change-me") ||
					strings.Contains(val, "example") {
					results = append(results, ValidationResult{
						File:     filePath,
						Severity: Warning,
						Message:  key + " contains placeholder/example value",
					})
				}
			}
		}
	}

	// Empty data validation
	if hasData {
		dataMap, ok := dataSection.(map[string]interface{})
		if ok && len(dataMap) == 0 {
			results = append(results, ValidationResult{
				File:     filePath,
				Severity: Critical,
				Message:  "Secret data section is empty",
			})
		}
	}

	return results
}

// ValidateSecretBytes is a helper for validating Secret content from byte slices (used in rendered YAML scanning)

func ValidateSecretBytes(content []byte) []ValidationResult {
	var results []ValidationResult
	
	var parsed map[string]interface{}

	kind, ok := parsed["kind"].(string)
	if !ok || kind != "Secret" {
		return results
	}

	// type validation
	secretType, exists := parsed["type"]
	if !exists {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "Secret type is missing",
		})
	} else {
		if secretType == "Opaque" {
			results = append(results, ValidationResult{
				Severity: Info,
				Message:  "Secret type is Opaque",
			})
		}
	}

	// data or stringData required
	dataSection, hasData := parsed["data"]
	stringDataSection, hasStringData := parsed["stringData"]

	if !hasData && !hasStringData {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "Secret must contain data or stringData",
		})
		return results
	}

	// Validate stringData placeholders
	if hasStringData {
		results = append(results, ValidationResult{
			Severity: Info,
			Message:  "Using stringData for Secret (acceptable for development)",
		})

		stringMap, ok := stringDataSection.(map[string]interface{})
		if ok {
			for key, value := range stringMap {
				val := strings.ToLower(value.(string))

				if strings.Contains(val, "your-api-key") ||
					strings.Contains(val, "password-here") ||
					strings.Contains(val, "change-me") ||
					strings.Contains(val, "example") {
					results = append(results, ValidationResult{
						Severity: Warning,
						Message:  key + " contains placeholder/example value",
					})
				}
			}
		}
	}

	// Empty data validation
	if hasData {
		dataMap, ok := dataSection.(map[string]interface{})
		if ok && len(dataMap) == 0 {
			results = append(results, ValidationResult{
				Severity: Critical,
				Message:  "Secret data section is empty",
			})
		}
	}

	return results
}