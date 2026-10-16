package validators

import (
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"gopkg.in/yaml.v3"
)

// ValidateWithCUE validates values.yaml against CUE schemas
func ValidateWithCUE(valuesFile string) error {
	ctx := cuecontext.New()

	schema, err := loadEmbeddedSchema(ctx)
	if err != nil {
		return fmt.Errorf("failed to load embedded schemas: %w", err)
	}

	data, err := os.ReadFile(valuesFile)
	if err != nil {
		return fmt.Errorf("failed to read values file: %w", err)
	}

	var yamlData interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return fmt.Errorf("invalid YAML format: %w", err)
	}

	dataMap, ok := yamlData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("values file must contain a YAML mapping")
	}

	// 🔹 Only encode resource sections that are actually present in values.yaml
	// This ensures projects with only a Deployment or Service aren't falsely forced
	// to define Ingress, TLS, etc.
	knownSections := []string{"deployment", "service", "configMap", "secret", "tlsSecret", "ingress", "namespace"}
	encodedMap := make(map[string]interface{})
	hasTLSSecret := false

	for _, section := range knownSections {
		if val, exists := dataMap[section]; exists && val != nil {
			encodedMap[section] = val
			if section == "tlsSecret" {
				hasTLSSecret = true
			}
		}
	}

	value := ctx.Encode(encodedMap)

	// 🔹 Unify schema + values
	result := schema.Unify(value)

	// Validate concrete and final
	if err := result.Validate(
		cue.Concrete(true),
		cue.Final(),
	); err != nil {
		return enhanceCUEError(err, hasTLSSecret)
	}

	return nil
}

// enhanceCUEError formats CUE validation errors into user-friendly messages
func enhanceCUEError(err error, hadTLS bool) error {
	details := errors.Details(err, nil)

	// Only show TLS remediation if the user actually configured tlsSecret
	if hadTLS && (strings.Contains(details, "tls.crt") || strings.Contains(details, "tls.key")) {
		return fmt.Errorf(
			"❌ TLS certificate is not configured properly in tlsSecret.\n"+
				"👉 Replace <base64-cert> and <base64-key> with real base64 values.\n\n%s",
			details,
		)
	}

	return fmt.Errorf("CUE validation failed:\n\n%s", details)
}

func loadEmbeddedSchema(ctx *cue.Context) (cue.Value, error) {
	files, err := schemaFS.ReadDir("schemas")
	if err != nil {
		return cue.Value{}, err
	}

	var unified cue.Value

	for _, file := range files {
		data, err := schemaFS.ReadFile("schemas/" + file.Name())
		if err != nil {
			return cue.Value{}, err
		}

		val := ctx.CompileBytes(data)

		if unified.Exists() {
			unified = unified.Unify(val)
		} else {
			unified = val
		}
	}

	return unified, nil
}