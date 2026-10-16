package validators

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ValidateService checks if the provided Kubernetes manifest is a Service and validates its structure and required fields.
func ValidateService(filePath string) []ValidationResult {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []ValidationResult{
			{
				File:     filePath,
				Severity: Critical,
				Message:  "Failed to read file: " + err.Error(),
			},
		}
	}

	results := ValidateServiceBytes(data)
	for i := range results {
		results[i].File = filePath
	}
	return results
}

// ValidateServiceBytes performs validation on raw YAML content.
func ValidateServiceBytes(content []byte) []ValidationResult {
	var results []ValidationResult
	var parsed map[string]interface{}

	if err := yaml.Unmarshal(content, &parsed); err != nil {
		return results
	}

	kind, ok := parsed["kind"].(string)
	if !ok || kind != "Service" {
		return results
	}

	spec, ok := parsed["spec"].(map[string]interface{})
	if !ok {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "spec section missing",
		})
		return results
	}

	// selector check
	if _, exists := spec["selector"]; !exists {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  "selector is missing",
		})
	}

	// type check
	if _, exists := spec["type"]; !exists {
		results = append(results, ValidationResult{
			Severity: Info,
			Message:  "service type not specified (default: ClusterIP)",
		})
	}

	// ports check
	portsRaw, exists := spec["ports"]
	if !exists {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "ports section missing",
		})
		return results
	}

	ports, ok := portsRaw.([]interface{})
	if !ok || len(ports) == 0 {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "ports configuration is invalid",
		})
		return results
	}

	port, ok := ports[0].(map[string]interface{})
	if !ok {
		return results
	}

	if _, exists := port["port"]; !exists {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "service port is missing",
		})
	}

	if _, exists := port["targetPort"]; !exists {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  "targetPort not defined",
		})
	}

	if _, exists := port["protocol"]; !exists {
		results = append(results, ValidationResult{
			Severity: Info,
			Message:  "protocol not specified (default: TCP)",
		})
	}

	return results
}