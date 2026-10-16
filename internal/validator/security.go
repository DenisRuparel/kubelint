package validators

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SecurityDeployment is the typed model used for security-focused validation.
type SecurityDeployment struct {
	Kind string `yaml:"kind"`
	Spec struct {
		Template struct {
			Spec struct {
				Containers []SecurityContainer `yaml:"containers"`
				InitContainers []SecurityContainer `yaml:"initContainers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

// SecurityContainer holds only the fields relevant to security checks.
type SecurityContainer struct {
	Name      string `yaml:"name"`
	Resources struct {
		Limits map[string]string `yaml:"limits"`
	} `yaml:"resources"`
	SecurityContext struct {
		Privileged               *bool `yaml:"privileged"`
		RunAsNonRoot             *bool `yaml:"runAsNonRoot"`
		ReadOnlyRootFilesystem   *bool `yaml:"readOnlyRootFilesystem"`
		AllowPrivilegeEscalation *bool `yaml:"allowPrivilegeEscalation"`
	} `yaml:"securityContext"`
}

// ValidateSecurity reads a file and delegates to ValidateSecurityBytes.
func ValidateSecurity(filePath string) []ValidationResult {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []ValidationResult{{
			File:     filePath,
			Severity: Critical,
			Message:  "Failed to read file: " + err.Error(),
		}}
	}

	results := ValidateSecurityBytes(data)
	for i := range results {
		results[i].File = filePath
	}
	return results
}

// ValidateSecurityBytes validates security settings in raw YAML content.
// Bug fix: previously never called yaml.Unmarshal, so deployment was always
// a zero-value struct and all checks were silently skipped.
func ValidateSecurityBytes(content []byte) []ValidationResult {
	var results []ValidationResult
	var deployment SecurityDeployment

	if err := yaml.Unmarshal(content, &deployment); err != nil {
		// Not our document to validate; not a parse error we own.
		return results
	}

	if deployment.Kind != "Deployment" {
		return results
	}

	allContainers := append(deployment.Spec.Template.Spec.Containers, deployment.Spec.Template.Spec.InitContainers...)
	for _, c := range allContainers {
		results = append(results, checkContainerSecurity(c)...)
	}

	return results
}

// checkContainerSecurity runs all security checks for a single container.
func checkContainerSecurity(c SecurityContainer) []ValidationResult {
	var results []ValidationResult
	name := c.Name

	if c.SecurityContext.Privileged == nil {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  fmt.Sprintf("Container '%s': privileged not set (should be false)", name),
		})
	}

	if c.SecurityContext.RunAsNonRoot == nil || !*c.SecurityContext.RunAsNonRoot {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  fmt.Sprintf("Container '%s': runAsNonRoot should be true", name),
		})
	}

	if c.SecurityContext.ReadOnlyRootFilesystem == nil || !*c.SecurityContext.ReadOnlyRootFilesystem {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  fmt.Sprintf("Container '%s': readOnlyRootFilesystem should be true", name),
		})
	}

	if c.SecurityContext.AllowPrivilegeEscalation == nil || *c.SecurityContext.AllowPrivilegeEscalation {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  fmt.Sprintf("Container '%s': allowPrivilegeEscalation should be false", name),
		})
	}

	if _, ok := c.Resources.Limits["cpu"]; !ok {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  fmt.Sprintf("Container '%s': CPU limit not set", name),
		})
	}

	if _, ok := c.Resources.Limits["memory"]; !ok {
		results = append(results, ValidationResult{
			Severity: Warning,
			Message:  fmt.Sprintf("Container '%s': memory limit not set", name),
		})
	}

	return results
}