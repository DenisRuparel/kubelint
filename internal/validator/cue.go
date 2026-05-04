package validators

import (
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
	"gopkg.in/yaml.v3"
)

// ValidateWithCUE validates values.yaml against CUE schemas
func ValidateWithCUE(valuesFile string) error {
	ctx := cuecontext.New()

	// 🔹 Load schema
	instances := load.Instances([]string{"./schemas/..."}, nil)
	if len(instances) == 0 {
		return fmt.Errorf("failed to load CUE schemas")
	}

	schema := ctx.BuildInstance(instances[0])
	if schema.Err() != nil {
		return schema.Err()
	}

	// 🔹 Read YAML
	data, err := os.ReadFile(valuesFile)
	if err != nil {
		return fmt.Errorf("failed to read values file: %w", err)
	}

	var yamlData interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return fmt.Errorf("invalid YAML format: %w", err)
	}

	dataMap := yamlData.(map[string]interface{})

	// 🔹 Encode into CUE
	value := ctx.Encode(map[string]interface{}{
		"deployment": dataMap["deployment"],
		"service":    dataMap["service"],
		"configMap":  dataMap["configMap"],
		"secret":     dataMap["secret"],
		"tlsSecret":  dataMap["tlsSecret"],
		"ingress":    dataMap["ingress"],
		"namespace":  dataMap["namespace"],
	})

	// 🔹 Unify schema + values
	result := schema.Unify(value)

	// 🔥 VALIDATION (THIS WAS MISSING)
	if err := result.Validate(
		cue.Concrete(true),
		cue.Final(),
	); err != nil {
		return enhanceCUEError(err)
	}

	return nil
}

// 🔥 Custom error formatter (OUTSIDE function)
func enhanceCUEError(err error) error {
	details := errors.Details(err, nil)

	// 🔥 TLS-specific error
	if strings.Contains(details, "tls.crt") || strings.Contains(details, "tls.key") {
		return fmt.Errorf(
			"❌ TLS certificate is not configured properly.\n"+
				"👉 Replace <base64-cert> and <base64-key> with real base64 values.\n\n%s",
			details,
		)
	}

	return fmt.Errorf("CUE validation failed:\n\n%s", details)
}