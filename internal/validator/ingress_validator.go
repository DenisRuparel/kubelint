package validators

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ValidateIngress(filePath string) []ValidationResult {
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
	if !ok || kind != "Ingress" {
		return results
	}

	// spec check
	specRaw, exists := parsed["spec"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "spec section missing",
		})
		return results
	}

	spec, ok := specRaw.(map[string]interface{})
	if !ok {
		return results
	}

	// ingressClassName
	if _, exists := spec["ingressClassName"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "ingressClassName missing",
		})
	}

	// tls check
	if _, exists := spec["tls"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Info,
			Message:  "TLS not configured (recommended for production)",
		})
	}

	// rules check
	rulesRaw, exists := spec["rules"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "rules section missing",
		})
		return results
	}

	rules, ok := rulesRaw.([]interface{})
	if !ok || len(rules) == 0 {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "rules configuration is invalid",
		})
		return results
	}

	rule, ok := rules[0].(map[string]interface{})
	if !ok {
		return results
	}

	// host check
	if _, exists := rule["host"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "host missing in ingress rule",
		})
	}

	httpRaw, exists := rule["http"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "http section missing in ingress rule",
		})
		return results
	}

	httpMap, ok := httpRaw.(map[string]interface{})
	if !ok {
		return results
	}

	pathsRaw, exists := httpMap["paths"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "paths section missing",
		})
		return results
	}

	paths, ok := pathsRaw.([]interface{})
	if !ok || len(paths) == 0 {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "paths configuration is invalid",
		})
		return results
	}

	path, ok := paths[0].(map[string]interface{})
	if !ok {
		return results
	}

	// pathType check
	if _, exists := path["pathType"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "pathType missing",
		})
	}

	backendRaw, exists := path["backend"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "backend section missing",
		})
		return results
	}

	backend, ok := backendRaw.(map[string]interface{})
	if !ok {
		return results
	}

	serviceRaw, exists := backend["service"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "backend service missing",
		})
		return results
	}

	service, ok := serviceRaw.(map[string]interface{})
	if !ok {
		return results
	}

	if _, exists := service["name"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "backend service name missing",
		})
	}

	portRaw, exists := service["port"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "backend service port missing",
		})
		return results
	}

	port, ok := portRaw.(map[string]interface{})
	if !ok {
		return results
	}

	if _, exists := port["number"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "backend service port number missing",
		})
	}

	return results
}