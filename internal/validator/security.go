package validators

import (
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

type SecurityDeployment struct {
	Kind string `yaml:"kind"`
	Spec struct {
		Template struct {
			Spec struct {
				Containers []struct {
					Name string `yaml:"name"`

					SecurityContext struct {
						Privileged               *bool `yaml:"privileged"`
						RunAsNonRoot             *bool `yaml:"runAsNonRoot"`
						ReadOnlyRootFilesystem   *bool `yaml:"readOnlyRootFilesystem"`
						AllowPrivilegeEscalation *bool `yaml:"allowPrivilegeEscalation"`
					} `yaml:"securityContext"`

					Resources struct {
						Limits map[string]string `yaml:"limits"`
					} `yaml:"resources"`
				} `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

// ValidateSecurity checks for security best practices in Deployment manifests

func ValidateSecurity(filePath string) []ValidationResult {
	var results []ValidationResult

	content, err := os.ReadFile(filePath)
	if err != nil {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  "Failed to read file for security validation",
		})
		return results
	}

	var deployment SecurityDeployment

	err = yaml.Unmarshal(content, &deployment)
	if err != nil {
		results = append(results, ValidationResult{
			File:     filePath,
			Severity: Critical,
			Message:  fmt.Sprintf("Failed to parse security validation: %v", err),
		})
		return results
	}

	if deployment.Kind != "Deployment" {
		return results
	}

	for _, c := range deployment.Spec.Template.Spec.Containers {

		if c.SecurityContext.Privileged == nil {
			results = append(results, ValidationResult{
				File:     filePath,
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should explicitly define privileged=false", c.Name),
			})
		}

		if c.SecurityContext.RunAsNonRoot == nil ||
			!*c.SecurityContext.RunAsNonRoot {
			results = append(results, ValidationResult{
				File:     filePath,
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should set runAsNonRoot=true", c.Name),
			})
		}

		if c.SecurityContext.ReadOnlyRootFilesystem == nil ||
			!*c.SecurityContext.ReadOnlyRootFilesystem {
			results = append(results, ValidationResult{
				File:     filePath,
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should set readOnlyRootFilesystem=true", c.Name),
			})
		}

		if c.SecurityContext.AllowPrivilegeEscalation == nil ||
			*c.SecurityContext.AllowPrivilegeEscalation {
			results = append(results, ValidationResult{
				File:     filePath,
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should set allowPrivilegeEscalation=false", c.Name),
			})
		}

		if _, exists := c.Resources.Limits["cpu"]; !exists {
			results = append(results, ValidationResult{
				File:     filePath,
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should define CPU limits", c.Name),
			})
		}

		if _, exists := c.Resources.Limits["memory"]; !exists {
			results = append(results, ValidationResult{
				File:     filePath,
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should define Memory limits", c.Name),
			})
		}
	}

	return results
}

// ValidateSecurityBytes is a helper for validating security in YAML content without file I/O

func ValidateSecurityBytes(content []byte) []ValidationResult {
	var results []ValidationResult

	var deployment SecurityDeployment

	if deployment.Kind != "Deployment" {
		return results
	}

	for _, c := range deployment.Spec.Template.Spec.Containers {

		if c.SecurityContext.Privileged == nil {
			results = append(results, ValidationResult{
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should explicitly define privileged=false", c.Name),
			})
		}

		if c.SecurityContext.RunAsNonRoot == nil ||
			!*c.SecurityContext.RunAsNonRoot {
			results = append(results, ValidationResult{
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should set runAsNonRoot=true", c.Name),
			})
		}

		if c.SecurityContext.ReadOnlyRootFilesystem == nil ||
			!*c.SecurityContext.ReadOnlyRootFilesystem {
			results = append(results, ValidationResult{
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should set readOnlyRootFilesystem=true", c.Name),
			})
		}

		if c.SecurityContext.AllowPrivilegeEscalation == nil ||
			*c.SecurityContext.AllowPrivilegeEscalation {
			results = append(results, ValidationResult{
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should set allowPrivilegeEscalation=false", c.Name),
			})
		}

		if _, exists := c.Resources.Limits["cpu"]; !exists {
			results = append(results, ValidationResult{
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should define CPU limits", c.Name),
			})
		}

		if _, exists := c.Resources.Limits["memory"]; !exists {
			results = append(results, ValidationResult{
				Severity: Warning,
				Message:  fmt.Sprintf("Container '%s' should define Memory limits", c.Name),
			})
		}
	}

	return results
}