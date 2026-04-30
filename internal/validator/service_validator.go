package validators

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ValidateService checks if the provided Kubernetes manifest is a Service and validates its structure and required fields. It returns a slice of ValidationResult indicating any issues found.

func ValidateService(filePath string) []ValidationResult {
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

	var parsed map[string]interface{}

	err = yaml.Unmarshal(data, &parsed)
	if err != nil {
		return results
	}

	kind, ok := parsed["kind"].(string)
	if !ok || kind != "Service" {
		return results
	}

	spec, ok := parsed["spec"].(map[string]interface{})
	if !ok {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "spec section missing",
		})
		return results
	}

	// selector check
	if _, exists := spec["selector"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "selector is missing",
		})
	}

	// type check
	if _, exists := spec["type"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Info,
			Message:  "service type not specified (default: ClusterIP)",
		})
	}

	// ports check
	portsRaw, exists := spec["ports"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "ports section missing",
		})
		return results
	}

	ports, ok := portsRaw.([]interface{})
	if !ok || len(ports) == 0 {
		results = append(results, ValidationResult{
			File:     filePath,
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
			File:     filePath,
			Severity: Critical,
			Message:  "service port is missing",
		})
	}

	if _, exists := port["targetPort"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "targetPort not defined",
		})
	}

	if _, exists := port["protocol"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Info,
			Message:  "protocol not specified (default: TCP)",
		})
	}

	return results
}

// ValidateServiceBytes performs the same validation as ValidateService but operates on raw YAML content instead of a file path. It returns a slice of ValidationResult indicating any issues found.

func ValidateServiceBytes(content []byte) []ValidationResult {
	var results []ValidationResult

	var parsed map[string]interface{}

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