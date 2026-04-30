package loader

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadValues(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read values file: %v", err)
	}

	var values map[string]interface{}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values.yaml: %v", err)
	}

	return values, nil
}