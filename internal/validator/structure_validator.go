package validators

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ValidateStructure(filePath string) ValidationResult {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "Failed to read file",
		}
	}

	var root yaml.Node

	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  err.Error(),
		}
	}

	// YAML document root check
	if len(root.Content) == 0 {
		return ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "Empty YAML document",
		}
	}

	doc := root.Content[0]

	// Required top-level fields
	requiredFields := []string{"apiVersion", "kind", "metadata"}

	for _, field := range requiredFields {
		if !hasField(doc, field) {
			return ValidationResult{
				File:     filePath,
				Severity: Critical,
				Message:  fmt.Sprintf("Missing required field: %s", field),
			}
		}
	}

	// Deployment-specific structure validation
	// kindNode := getField(doc, "kind")
	// if kindNode != nil && kindNode.Value == "Deployment" {
	// 	specNode := getField(doc, "spec")
	// 	if specNode == nil {
	// 		return ValidationResult{
	// 			File:     filePath,
	// 			Severity: Critical,
	// 			Message:  "Missing spec section",
	// 		}
	// 	}

	// 	templateNode := getField(specNode, "template")
	// 	if templateNode == nil {
	// 		return ValidationResult{
	// 			File:     filePath,
	// 			Severity: Critical,
	// 			Message:  fmt.Sprintf(
	// 				"Line %d: Missing spec.template section",
	// 				specNode.Line,
	// 			),
	// 		}
	// 	}

	// 	templateSpecNode := getField(templateNode, "spec")
	// 	if templateSpecNode == nil {
	// 		return ValidationResult{
	// 			File:     filePath,
	// 			Severity: Critical,
	// 			Message: fmt.Sprintf(
	// 				"Line %d: Indentation issue → spec.template.spec missing",
	// 				templateNode.Line,
	// 			),
	// 		}
	// 	}
	// }

	return ValidationResult{
		File:     filePath,
		Severity: Info,
		Message:  "Kubernetes structure validation passed",
	}
}

func hasField(node *yaml.Node, field string) bool {
	return getField(node, field) != nil
}

func getField(node *yaml.Node, field string) *yaml.Node {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content)-1; i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		if keyNode.Value == field {
			return valueNode
		}
	}

	return nil
}