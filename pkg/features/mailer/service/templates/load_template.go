package templates

import (
	"bytes"
	"fmt"
	"html/template"
)

func LoadTemplate(templatePath string, data any) (string, error) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template file %q: %w", templatePath, err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to execute template %q with provided data: %w", templatePath, err)
	}

	return body.String(), nil
}