package validators

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ValidateYAMLSyntax checks if the provided file contains valid YAML syntax. It also detects if the file appears to be a template (contains {{ }}), in which case it advises the user to run 'kubelint build' first.

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

// ValidateYAMLSyntaxBytes is a variant of ValidateYAMLSyntax that operates on raw YAML content instead of file paths. This is useful for scanning rendered manifests in-memory without needing to read from disk.

func ValidateYAMLSyntaxBytes(content []byte) ValidationResult {

	var parsed interface{}

	err := yaml.Unmarshal(content, &parsed)
	if err != nil {
		return ValidationResult{
			Severity: Critical,
			Message:  "YAML syntax error → " + err.Error(),
		}
	}

	return ValidationResult{
		Severity: Info,
		Message:  "YAML syntax validation passed",
	}
}