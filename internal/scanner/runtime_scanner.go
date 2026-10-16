package scanner

import (
	"fmt"
	"strings"

	validators "github.com/DenisRuparel/kubelint/internal/validator"
	"gopkg.in/yaml.v3"
)

type ScanSummary struct {
	Critical int
	Warning  int
	Info     int
}

type ScanResult struct {
	Summary ScanSummary
	Issues  []validators.ValidationResult
}

// splitYAMLDocuments splits a multi-document YAML stream on lines that are exactly "---".
// This prevents naive substring splits on "-----BEGIN CERTIFICATE-----" or inside strings.
func splitYAMLDocuments(yamlContent string) []string {
	var docs []string
	var current strings.Builder

	lines := strings.Split(yamlContent, "\n")
	for _, line := range lines {
		trimmed := strings.TrimRight(line, "\r")
		if strings.TrimSpace(trimmed) == "---" {
			docStr := strings.TrimSpace(current.String())
			if docStr != "" {
				docs = append(docs, current.String())
			}
			current.Reset()
			continue
		}
		current.WriteString(line)
		current.WriteByte('\n')
	}

	docStr := strings.TrimSpace(current.String())
	if docStr != "" {
		docs = append(docs, current.String())
	}

	return docs
}

func ScanRenderedYAML(yamlContent string) ScanResult {
	var result ScanResult

	docs := splitYAMLDocuments(yamlContent)

	for _, docStr := range docs {
		trimmedDoc := strings.TrimSpace(docStr)
		if trimmedDoc == "" {
			continue
		}

		// Extract fileName per document if present
		var fileName string
		lines := strings.Split(trimmedDoc, "\n")
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, "# FILE:") {
				fileName = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "# FILE:"))
				break
			}
		}

		content := []byte(trimmedDoc)

		// Parse once (fixes bug #8: duplicate yaml.Unmarshal)
		var parsed map[string]interface{}
		err := yaml.Unmarshal(content, &parsed)
		if err != nil {
			msg := err.Error()
			msg = strings.Replace(msg, "yaml:", "", 1)
			msg = strings.TrimSpace(msg)

			displayID := fileName
			if displayID == "" {
				displayID = "YAML Stream"
			}

			result.Issues = append(result.Issues, validators.ValidationResult{
				Severity: validators.Critical,
				Message: fmt.Sprintf("[%s] Invalid YAML syntax: %s (fix syntax before other validations)",
					displayID,
					msg,
				),
			})

			result.Summary.Critical++
			continue
		}

		if len(parsed) == 0 {
			continue
		}

		kind, _ := parsed["kind"].(string)
		metadata, _ := parsed["metadata"].(map[string]interface{})
		name, _ := metadata["name"].(string)

		var resourceID string
		if fileName != "" {
			resourceID = fileName
		} else {
			if kind == "" {
				kind = "Unknown"
			}
			if name == "" {
				name = "unknown"
			}
			resourceID = fmt.Sprintf("%s/%s", kind, name)
		}

		// Run in-memory validators
		var results []validators.ValidationResult
		results = append(results, validators.ValidateCommonBytes(content)...)
		results = append(results, validators.ValidateStructureBytes(content))
		results = append(results, validators.ValidateDeploymentBytes(content)...)
		results = append(results, validators.ValidateSecurityBytes(content)...)
		results = append(results, validators.ValidateServiceBytes(content)...)
		results = append(results, validators.ValidateConfigMapBytes(content)...)
		results = append(results, validators.ValidateSecretBytes(content)...)
		results = append(results, validators.ValidateIngressBytes(content)...)
		results = append(results, validators.ValidateNamespacePolicyBytes(content)...)

		for _, r := range results {
			if r.Severity == "" {
				r.Severity = validators.Critical
			}

			if resourceID != "/" && resourceID != "" {
				r.Message = fmt.Sprintf("[%s] %s", resourceID, r.Message)
			}

			result.Issues = append(result.Issues, r)

			switch r.Severity {
			case validators.Critical:
				result.Summary.Critical++
			case validators.Warning:
				result.Summary.Warning++
			case validators.Info:
				result.Summary.Info++
			}
		}
	}
	return result
}
