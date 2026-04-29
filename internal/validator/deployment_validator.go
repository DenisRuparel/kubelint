package validators

import (
	"os"
	"strings"
)

func ValidateDeployment(filePath string) []ValidationResult {
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

	// Check if this is actually a Deployment
	if !strings.Contains(content, "kind: Deployment") {
		return results
	}

	// Example Rule 1: resources.limits check
	if !strings.Contains(content, "limits:") {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "resources.limits not defined",
		})
	}

	// Example Rule 2: replicas check
	if !strings.Contains(content, "replicas:") {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "replicas not defined",
		})
	}

	// Example Rule 3: imagePullPolicy check
	if !strings.Contains(content, "imagePullPolicy:") {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Info,
			Message:  "imagePullPolicy not set (recommended)",
		})
	}

	return results
}