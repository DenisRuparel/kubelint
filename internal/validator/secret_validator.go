package validators

import (
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// ValidateSecret reads a file and delegates to ValidateSecretBytes.
func ValidateSecret(filePath string) []ValidationResult {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []ValidationResult{{
			File:     filePath,
			Severity: Critical,
			Message:  "Failed to read file: " + err.Error(),
		}}
	}

	results := ValidateSecretBytes(data)
	for i := range results {
		results[i].File = filePath
	}
	return results
}

// ValidateSecretBytes validates a Kubernetes Secret manifest from raw YAML content.
// Bug fix: previously never called yaml.Unmarshal, so parsed was always nil and
// all checks silently passed.
func ValidateSecretBytes(content []byte) []ValidationResult {
	var results []ValidationResult
	var parsed map[string]interface{}

	if err := yaml.Unmarshal(content, &parsed); err != nil {
		return results
	}

	kind, ok := parsed["kind"].(string)
	if !ok || kind != "Secret" {
		return results
	}

	// type validation
	secretType, exists := parsed["type"]
	if !exists {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "Secret: type field is missing",
		})
	} else if secretType == "Opaque" {
		results = append(results, ValidationResult{
			Severity: Info,
			Message:  "Secret: type is Opaque",
		})
	}

	dataSection, hasData := parsed["data"]
	stringDataSection, hasStringData := parsed["stringData"]

	if !hasData && !hasStringData {
		results = append(results, ValidationResult{
			Severity: Critical,
			Message:  "Secret: must contain data or stringData",
		})
		return results
	}

	if hasStringData {
		results = append(results, ValidationResult{
			Severity: Info,
			Message:  "Secret: uses stringData (acceptable for development, avoid in production)",
		})

		if stringMap, ok := stringDataSection.(map[string]interface{}); ok {
			for key, value := range stringMap {
				val, _ := value.(string)
				if isPlaceholderValue(strings.ToLower(val)) {
					results = append(results, ValidationResult{
						Severity: Warning,
						Message:  "Secret: key '" + key + "' contains a placeholder or example value",
					})
				}
			}
		}
	}

	if hasData {
		if dataMap, ok := dataSection.(map[string]interface{}); ok && len(dataMap) == 0 {
			results = append(results, ValidationResult{
				Severity: Critical,
				Message:  "Secret: data section is empty",
			})
		}
	}

	return results
}

// secretKeyPattern matches common secret-like key names.
var secretKeyPattern = regexp.MustCompile(`(?i)(password|passwd|secret|token|apikey|api_key|credential|auth|jwt|private_key)`)

// weakValues is a set of known weak/demo secret values.
var weakValues = map[string]bool{
	"admin": true, "secret": true, "password": true, "12345": true,
	"changeme": true, "change-me": true, "test": true, "demo": true,
}

// isPlaceholderValue returns true if the value looks like a placeholder or demo value.
func isPlaceholderValue(val string) bool {
	keywords := []string{"your-api-key", "password-here", "change-me", "example", "placeholder", "todo", "fixme", "<", ">"}
	for _, kw := range keywords {
		if strings.Contains(val, kw) {
			return true
		}
	}
	return false
}

// ValidatePlaintextSecrets scans a values.yaml file for plaintext credentials
// in fields whose key names suggest they hold secrets. Returns Critical findings
// so that build and apply are blocked before secrets leak into manifests or git.
func ValidatePlaintextSecrets(valuesFile string) []ValidationResult {
	data, err := os.ReadFile(valuesFile)
	if err != nil {
		return []ValidationResult{{
			File:     valuesFile,
			Severity: Critical,
			Message:  "Failed to read values file: " + err.Error(),
		}}
	}

	var values map[string]interface{}
	if err := yaml.Unmarshal(data, &values); err != nil {
		return nil // YAML errors are reported elsewhere
	}

	var results []ValidationResult
	scanForSecrets(values, "", valuesFile, &results)
	return results
}

// scanForSecrets recursively walks a values map looking for secret-named keys
// with non-empty plaintext values.
func scanForSecrets(m map[string]interface{}, prefix, file string, results *[]ValidationResult) {
	for k, v := range m {
		keyPath := k
		if prefix != "" {
			keyPath = prefix + "." + k
		}

		switch val := v.(type) {
		case map[string]interface{}:
			scanForSecrets(val, keyPath, file, results)
		case string:
			if val == "" {
				continue
			}
			if secretKeyPattern.MatchString(k) {
				sev := Warning
				msg := "values.yaml: '" + keyPath + "' looks like a credential — use a secrets manager instead of plaintext values"

				// Escalate to Critical for known-weak values
				if weakValues[strings.ToLower(val)] || len(val) < 8 {
					sev = Critical
					msg = "values.yaml: '" + keyPath + "' contains a weak or demo credential — this must not reach production"
				}

				*results = append(*results, ValidationResult{
					File:     file,
					Severity: sev,
					Message:  msg,
				})
			}
		}
	}
}