package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/DenisRuparel/kubelint/internal/loader"
	"github.com/DenisRuparel/kubelint/internal/renderer"
	validators "github.com/DenisRuparel/kubelint/internal/validator"
)

func Build(projectPath, valuesFile string) (string, error) {
	templateDir := filepath.Join(projectPath, "templates")

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return "", fmt.Errorf("templates folder not found at: %s", templateDir)
	}

	// Resolve values file path
	if valuesFile == "" {
		valuesFile = filepath.Join(templateDir, "values.yaml")
	} else {
		if !filepath.IsAbs(valuesFile) {
			valuesFile = filepath.Join(projectPath, valuesFile)
		}
	}

	if _, err := os.Stat(valuesFile); os.IsNotExist(err) {
		return "", fmt.Errorf("values.yaml not found at: %s", valuesFile)
	}

	// 1. CUE schema validation (Fixes Gap A)
	if err := validators.ValidateWithCUE(valuesFile); err != nil {
		return "", fmt.Errorf("CUE validation failed:\n%w", err)
	}

	// 2. Plaintext secrets detection in values.yaml (Fixes Gap C)
	secretIssues := validators.ValidatePlaintextSecrets(valuesFile)
	for _, issue := range secretIssues {
		if issue.Severity == validators.Critical {
			return "", fmt.Errorf("plaintext secret detected in %s:\n  • %s", valuesFile, issue.Message)
		}
	}

	// 3. Load values
	values, err := loader.LoadValues(valuesFile)
	if err != nil {
		return "", err
	}

	files, err := os.ReadDir(templateDir)
	if err != nil {
		return "", fmt.Errorf("failed to read templates directory: %w", err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	var namespaceDocs []string
	var otherDocs []string
	var renderErrors []string

	for _, f := range files {
		if f.IsDir() || f.Name() == "values.yaml" || filepath.Ext(f.Name()) != ".yaml" {
			continue
		}

		fp := filepath.Join(templateDir, f.Name())

		rendered, err := renderer.RenderTemplate(fp, values)
		if err != nil {
			renderErrors = append(renderErrors, fmt.Sprintf("[%s] %v", f.Name(), err))
			continue
		}

		// Attach file metadata for scanner
		rendered = fmt.Sprintf("# FILE: %s\n%s", fp, rendered)

		doc := "---\n" + rendered
		if len(rendered) > 0 && rendered[len(rendered)-1] != '\n' {
			doc += "\n"
		}

		// Detect Namespace resource to order first
		if strings.Contains(rendered, "kind: Namespace") {
			namespaceDocs = append(namespaceDocs, doc)
		} else {
			otherDocs = append(otherDocs, doc)
		}
	}

	// Check if any templates failed to render (Fixes Bug 3: discarded render errors)
	if len(renderErrors) > 0 {
		return "", fmt.Errorf("template rendering failed (%d error(s)):\n  • %s",
			len(renderErrors),
			strings.Join(renderErrors, "\n  • "),
		)
	}

	// Combine: Namespace FIRST
	var final string
	for _, d := range namespaceDocs {
		final += d
	}
	for _, d := range otherDocs {
		final += d
	}

	return final, nil
}
