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

	// Split multi-doc YAML
	decoder := yaml.NewDecoder(strings.NewReader(yamlContent))

	for {
		var parsed map[string]interface{}

		err := decoder.Decode(&parsed)
		if err != nil {
			break
		}

		if parsed == nil {
			continue
		}

		content, _ := yaml.Marshal(parsed)

		var results []validators.ValidationResult

		kind, _ := parsed["kind"].(string)

		metadata, _ := parsed["metadata"].(map[string]interface{})
		name, _ := metadata["name"].(string)

		resourceID := fmt.Sprintf("%s/%s", kind, name)

		// YAML syntax
		syntax := validators.ValidateYAMLSyntaxBytes(content)
		results = append(results, syntax)

		if syntax.Severity != validators.Critical {

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
