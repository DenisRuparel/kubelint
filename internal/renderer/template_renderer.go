package renderer

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func RenderTemplate(filePath string, values map[string]interface{}) (string, error) {
	// Read template file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %v", err)
	}

	// Parse template with Sprig functions
	tmpl, err := template.New("k8s").
		Funcs(sprig.TxtFuncMap()).
		Option("missingkey=zero"). // 🔥 IMPORTANT FIX
		Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	// Execute template
	var output bytes.Buffer
	err = tmpl.Execute(&output, values)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %v", err)
	}

	return output.String(), nil
}