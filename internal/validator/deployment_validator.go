package validators

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ValidateDeployment(filePath string) []ValidationResult {
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
	if !ok || kind != "Deployment" {
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

	// replicas
	if _, exists := spec["replicas"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "replicas not defined",
		})
	}

	// selector
	if _, exists := spec["selector"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "selector is missing",
		})
	}

	// template
	templateRaw, exists := spec["template"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "template section missing",
		})
		return results
	}

	template, ok := templateRaw.(map[string]interface{})
	if !ok {
		return results
	}

	templateSpecRaw, exists := template["spec"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "template.spec missing",
		})
		return results
	}

	templateSpec, ok := templateSpecRaw.(map[string]interface{})
	if !ok {
		return results
	}

	// serviceAccountName
	if _, exists := templateSpec["serviceAccountName"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Info,
			Message:  "serviceAccountName not set (recommended)",
		})
	}

	containersRaw, exists := templateSpec["containers"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "containers section missing",
		})
		return results
	}

	containers, ok := containersRaw.([]interface{})
	if !ok || len(containers) == 0 {
		return results
	}

	container, ok := containers[0].(map[string]interface{})
	if !ok {
		return results
	}

	// resources
	if _, exists := container["resources"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "resources section missing",
		})
	}

	// livenessProbe
	if _, exists := container["livenessProbe"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "livenessProbe missing",
		})
	}

	// readinessProbe
	if _, exists := container["readinessProbe"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Warning,
			Message:  "readinessProbe missing",
		})
	}

	// imagePullPolicy
	if _, exists := container["imagePullPolicy"]; !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Info,
			Message:  "imagePullPolicy not set (recommended)",
		})
	}

	// securityContext
	securityRaw, exists := container["securityContext"]
	if !exists {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "securityContext missing",
		})
		return results
	}

	securityContext, ok := securityRaw.(map[string]interface{})
	if !ok {
		return results
	}

	checkSecurityField := func(field string) {
		if _, exists := securityContext[field]; !exists {
			results = append(results, ValidationResult{
				File:     filePath,
				Severity: Warning,
				Message:  field + " not configured in securityContext",
			})
		}
	}

	checkSecurityField("privileged")
	checkSecurityField("runAsNonRoot")
	checkSecurityField("readOnlyRootFilesystem")
	checkSecurityField("allowPrivilegeEscalation")

	return results
}