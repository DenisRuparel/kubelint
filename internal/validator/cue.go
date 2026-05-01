package validators

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
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

	// 🔥 IMPORTANT FIX: wrap under root
	dataMap := yamlData.(map[string]interface{})

	value := ctx.Encode(map[string]interface{}{
		"deployment": dataMap["deployment"],
		"service":    dataMap["service"],
		"configMap":  dataMap["configMap"],
		"secret":     dataMap["secret"],
		"tlsSecret":  dataMap["tlsSecret"],
		"ingress":    dataMap["ingress"],
		"namespace":  dataMap["namespace"],
	})

	// 🔥 Unify
	result := schema.Unify(value)

	// 🔥 Extract only concrete data
	if err := result.Validate(
		cue.Concrete(true),
		cue.Final(),
	); err != nil {
		return formatCUEError(err)
	}

	return nil
}

func formatCUEError(err error) error {
	details := errors.Details(err, nil)
	return fmt.Errorf("CUE validation failed:\n\n%s", details)
}