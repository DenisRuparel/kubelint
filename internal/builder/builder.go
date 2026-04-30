package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/DenisRuparel/kubelint/internal/loader"
	"github.com/DenisRuparel/kubelint/internal/renderer"
)

func Build(projectPath, valuesFile string) (string, error) {
	templateDir := filepath.Join(projectPath, "templates")

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return "", fmt.Errorf("templates folder not found at: %s", templateDir)
	}

	if valuesFile == "" {
		valuesFile = filepath.Join(templateDir, "values.yaml")
	}

	// allow project-relative values path
	if !filepath.IsAbs(valuesFile) {
		valuesFile = filepath.Join(projectPath, valuesFile)
	}

	if _, err := os.Stat(valuesFile); os.IsNotExist(err) {
		return "", fmt.Errorf("values.yaml not found at: %s", valuesFile)
	}

	values, err := loader.LoadValues(valuesFile)
	if err != nil {
		return "", err
	}

	files, err := os.ReadDir(templateDir)
	if err != nil {
		return "", fmt.Errorf("failed to read templates directory")
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	var final string

	for _, f := range files {
		if f.IsDir() || f.Name() == "values.yaml" || filepath.Ext(f.Name()) != ".yaml" {
			continue
		}

		fp := filepath.Join(templateDir, f.Name())

		rendered, err := renderer.RenderTemplate(fp, values)
		if err != nil {
			return "", fmt.Errorf("render %s: %v", f.Name(), err)
		}

		final += "---\n"
		final += rendered
		if len(rendered) > 0 && rendered[len(rendered)-1] != '\n' {
			final += "\n"
		}
	}

	return final, nil
}