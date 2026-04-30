package scanner

import (
	"fmt"
	validators "github.com/DenisRuparel/kubelint/internal/validator"
	"gopkg.in/yaml.v3"
	"strings"
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

func ScanRenderedYAML(yamlContent string) ScanResult {
	var result ScanResult

	docs := strings.Split(yamlContent, "---")

	for _, docStr := range docs {

		docStr = strings.TrimSpace(docStr)
		if docStr == "" {
			continue
		}

		// 🔥 Extract fileName per document
		var fileName string
		lines := strings.Split(docStr, "\n")
		if len(lines) > 0 && strings.HasPrefix(lines[0], "# FILE:") {
			fileName = strings.TrimSpace(strings.TrimPrefix(lines[0], "# FILE:"))
		}

		content := []byte(docStr)

		// 🔥 Syntax validation per document
		err := yaml.Unmarshal(content, &map[string]interface{}{})
		if err != nil {
			msg := err.Error()
			msg = strings.Replace(msg, "yaml:", "", 1)
			msg = strings.TrimSpace(msg)

			result.Issues = append(result.Issues, validators.ValidationResult{
				Severity: validators.Critical,
				Message: fmt.Sprintf("[%s] Invalid YAML syntax: %s (fix syntax before other validations)",
					fileName,
					msg,
				),
			})

			result.Summary.Critical++
			continue
		}

		// 🔥 Continue with your existing validators
		var parsed map[string]interface{}
		_ = yaml.Unmarshal(content, &parsed)

		var results []validators.ValidationResult

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

		// run validators
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

			if resourceID != "/" {
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
