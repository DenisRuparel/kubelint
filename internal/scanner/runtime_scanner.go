package scanner

import (
	"strings"

	"github.com/DenisRuparel/kubelint/internal/validator"
)

type ScanSummary struct {
	Critical int
	Warning  int
	Info     int
}

func ScanRenderedYAML(yamlContent string) ScanSummary {
	var summary ScanSummary

	// Split multi-doc YAML
	docs := strings.Split(yamlContent, "---")

	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		content := []byte(doc)

		var results []validators.ValidationResult

		// YAML syntax
		syntax := validators.ValidateYAMLSyntaxBytes(content)
		results = append(results, syntax)

		// Only continue if syntax OK
		if syntax.Severity != validators.Critical {

			// Structure
			results = append(results, validators.ValidateStructureBytes(content))

			// Common
			results = append(results, validators.ValidateCommonBytes(content)...)

			// Deployment
			results = append(results, validators.ValidateDeploymentBytes(content)...)

			// Security
			results = append(results, validators.ValidateSecurityBytes(content)...)

			// Service
			results = append(results, validators.ValidateServiceBytes(content)...)

			// ConfigMap
			results = append(results, validators.ValidateConfigMapBytes(content)...)

			// Secret
			results = append(results, validators.ValidateSecretBytes(content)...)

			// Ingress
			results = append(results, validators.ValidateIngressBytes(content)...)

			// Namespace
			results = append(results, validators.ValidateNamespacePolicyBytes(content)...)
		}

		// Count results
		for _, r := range results {
			switch r.Severity {
			case validators.Critical:
				summary.Critical++
			case validators.Warning:
				summary.Warning++
			case validators.Info:
				summary.Info++
			}
		}
	}

	return summary
}