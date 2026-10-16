package scanner

import (
	"fmt"
	"os"
	"path/filepath"
)

// CollectYAMLFiles gathers all .yaml and .yml manifest files from the specified paths.
func CollectYAMLFiles(paths []string) ([]string, error) {
	var yamlFiles []string

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("cannot access path %q: %w", path, err)
		}

		// If single file
		if !info.IsDir() {
			if isYAMLFile(path) {
				yamlFiles = append(yamlFiles, path)
			}
			continue
		}

		// If directory → recursive scan
		err = filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip hidden directories (like .git, .github)
			if fileInfo.IsDir() && len(fileInfo.Name()) > 0 && fileInfo.Name()[0] == '.' {
				return filepath.SkipDir
			}

			if !fileInfo.IsDir() && isYAMLFile(filePath) {
				yamlFiles = append(yamlFiles, filePath)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	if len(yamlFiles) == 0 {
		return nil, fmt.Errorf("no YAML files found in provided paths")
	}

	return yamlFiles, nil
}

func isYAMLFile(path string) bool {
	ext := filepath.Ext(path)
	base := filepath.Base(path)

	if ext != ".yaml" && ext != ".yml" {
		return false
	}

	// Ignore config and values files from direct manifest scanning
	if base == ".kubelint.yaml" ||
		base == ".kubelint.yml" ||
		base == "values.yaml" ||
		base == "values.yml" {
		return false
	}

	return true
}