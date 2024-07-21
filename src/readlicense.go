package src

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
)

// The values that the input expects
type LicenseInput struct {
	Project string
	Year    string
	Author  string
	License string
}

//go:embed licenses/*.txt
var licenses embed.FS

var LicenseMaps = map[string]string{
	"gpl3": "gpl-3.0.txt",
	"mit":  "mit.txt",
}

func License(input LicenseInput) (string, error) {
	file, ok := LicenseMaps[strings.ToLower(input.License)]
	if !ok {
		return "", fmt.Errorf("no license found for %s", input.License)
	}

	filePath := fmt.Sprintf("licenses/%s", file)

	cont, err := licenses.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read license file: %w", err)
	}

	tmpl, err := template.New("license").Parse(string(cont))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, input)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return output.String(), nil
}
