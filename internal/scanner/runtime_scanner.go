package scanner

import (
	"fmt"
	validators "github.com/DenisRuparel/kubelint/internal/validator"
	"gopkg.in/yaml.v3"
	"io"
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

	// 🔥 validate full YAML ONCE
	syntax := validators.ValidateYAMLSyntaxBytes([]byte(yamlContent))

	if syntax.Severity == validators.Critical {
		result.Issues = append(result.Issues, syntax)
		result.Summary.Critical++
		return result
	}

	decoder := yaml.NewDecoder(strings.NewReader(yamlContent))

	for {
		var node yaml.Node

		err := decoder.Decode(&node)
		if err != nil {
			if err == io.EOF {
				break
			}
			result.Issues = append(result.Issues, validators.ValidationResult{
				Severity: validators.Critical,
				Message:  fmt.Sprintf("[%s] Invalid YAML syntax: %v (fix syntax before other validations)", err),
			})

			result.Summary.Critical++
			return result
		}

		if len(node.Content) == 0 {
			continue
		}

		// 👇 THIS IS THE FIX — preserve original YAML
		if node.Kind != yaml.DocumentNode || len(node.Content) == 0 {
			continue
		}

		doc := node.Content[0]

		content, err := yaml.Marshal(doc)
		if err != nil {
			continue
		}

		// 🔥 Extract file name from injected comment
		var fileName string

		lines := strings.Split(string(content), "\n")
		if len(lines) > 0 && strings.HasPrefix(lines[0], "# FILE:") {
			fileName = strings.TrimSpace(strings.TrimPrefix(lines[0], "# FILE:"))
		}

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

		// YAML syntax
		docSyntax := validators.ValidateYAMLSyntaxBytes(content)
		results = append(results, docSyntax)

		if docSyntax.Severity != validators.Critical {
			results = append(results, validators.ValidateCommonBytes(content)...)
			results = append(results, validators.ValidateStructureBytes(content))
			results = append(results, validators.ValidateDeploymentBytes(content)...)
			results = append(results, validators.ValidateSecurityBytes(content)...)
			results = append(results, validators.ValidateServiceBytes(content)...)
			results = append(results, validators.ValidateConfigMapBytes(content)...)
			results = append(results, validators.ValidateSecretBytes(content)...)
			results = append(results, validators.ValidateIngressBytes(content)...)
			results = append(results, validators.ValidateNamespacePolicyBytes(content)...)
		}

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
