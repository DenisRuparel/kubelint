package validators

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ValidateYAMLSyntax(filePath string) ValidationResult {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "Failed to read file",
		}
	}

	content := string(data)

	// Template detection
	if strings.Contains(content, "{{") && strings.Contains(content, "}}") {
		return ValidationResult{
			File:     filePath,
			Severity: Info,
			Message:  "Template file detected → run 'kubelint build' first",
		}
	}

	var parsed interface{}

	err = yaml.Unmarshal(data, &parsed)
	if err != nil {
		return ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "YAML syntax error → " + err.Error(),
		}
	}

	return ValidationResult{
		File:     filePath,
		Severity: Info,
		Message:  "YAML syntax validation passed",
	}
}
