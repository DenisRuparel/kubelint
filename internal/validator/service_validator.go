package validators

import (
	"os"
	"strings"
)

func ValidateService(filePath string) []ValidationResult {
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

	content := string(data)

	// Check if this is actually a Service
	if !strings.Contains(content, "kind: Service") {
		return results
	}

	// Rule 1: selector check
	if !strings.Contains(content, "selector:") {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "selector not defined",
		})
	}

	// Rule 2: ports check
	if !strings.Contains(content, "ports:") {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "ports section missing",
		})
	}

	// Rule 3: service type recommendation
	if !strings.Contains(content, "type:") {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Info,
			Message:  "service type not specified (recommended)",
		})
	}

	return results
}