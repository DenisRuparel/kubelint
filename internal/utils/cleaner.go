package utils

import "strings"

func CleanYAML(yaml string) string {
	var result []string

	lines := strings.Split(yaml, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "# FILE:") {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}