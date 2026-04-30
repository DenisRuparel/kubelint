package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/DenisRuparel/kubelint/internal/loader"
	"github.com/DenisRuparel/kubelint/internal/renderer"
)

func Build(projectPath, valuesFile string) (string, error) {
	templateDir := filepath.Join(projectPath, "templates")

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return "", fmt.Errorf("templates folder not found at: %s", templateDir)
	}

	// Resolve values file path
	if valuesFile == "" {
		valuesFile = filepath.Join(projectPath, "templates", "values.yaml")
	} else {
		if !filepath.IsAbs(valuesFile) {
			valuesFile = filepath.Join(projectPath, valuesFile)
		}
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

	var namespaceDocs []string
	var otherDocs []string
	var errors []string

	for _, f := range files {
		if f.IsDir() || f.Name() == "values.yaml" || filepath.Ext(f.Name()) != ".yaml" {
			continue
		}

		fp := filepath.Join(templateDir, f.Name())

		rendered, err := renderer.RenderTemplate(fp, values)
		if err != nil {
			errors = append(errors,
				fmt.Sprintf("[%s] %v", f.Name(), err),
			)
			continue
		}

		doc := "---\n" + rendered
		if len(rendered) > 0 && rendered[len(rendered)-1] != '\n' {
			doc += "\n"
		}

		// Detect Namespace resource
		if strings.Contains(rendered, "kind: Namespace") {
			namespaceDocs = append(namespaceDocs, doc)
		} else {
			otherDocs = append(otherDocs, doc)
		}
	}

	// Combine: Namespace FIRST
	var final string

	for _, d := range namespaceDocs {
		final += d
	}
	for _, d := range otherDocs {
		final += d
	}

	if len(errors) > 0 {
		fmt.Println("\n⚠️ Build Issues:")
		fmt.Println("---------------------------------")

		for _, e := range errors {
			fmt.Printf("• %s\n", e)
		}

		fmt.Println("---------------------------------")
		fmt.Printf("Total Errors: %d\n\n", len(errors))
	}

	return final, nil
}
