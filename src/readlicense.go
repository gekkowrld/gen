package src

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

// The values that the input expects
type LicenseInput struct {
	Project    string
	Year       string
	Author     string
	License    string
	Output     string
	IsTemplate bool
	Template   string
}

// The metadata of a license text
type LicenseMeta struct {
	Title       string
	SpdxId      string
	Nickname    string
	Description string
	Note        string
	Permissions []string
	Conditions  []string
	Limitations []string
}

//go:embed licenses/*.txt
var licenses embed.FS

// Take in the template to be used.
// this will be used to generate the final license
func License(input LicenseInput) (string, error) {
	var exf ext

	if input.IsTemplate {
		exf.istmpl = true
		exf.tmpl = input.Template
	} else {
		exf.input = input.License
	}

	cont, err := extractText(exf)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("license").Parse(cont)
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

type ext struct {
	input  string
	ismeta bool
	istmpl bool
	tmpl   string
}

func AllLicense() map[string]string {
	lis := make(map[string]string)
	files, err := licenses.ReadDir("licenses")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	for _, agi := range files {
		agin := strings.TrimSpace(agi.Name())
		file := strings.TrimSuffix(agin, ".txt")
		lis[file] = agin
	}

	return lis
}

func extractText(fe ext) (string, error) {
	var content []byte
	var err error
	LicenseMaps := AllLicense()

	if !fe.istmpl {

		file, ok := LicenseMaps[strings.ToLower(fe.input)]
		if !ok {
			return "", fmt.Errorf("no license found for %s", fe.input)
		}

		filePath := fmt.Sprintf("licenses/%s", file)
		content, err = licenses.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read license file: %w", err)
		}
	} else {
		content, err = os.ReadFile(fe.tmpl)
		if err != nil {
			return "", err
		}
	}

	var text string
	blks := strings.SplitN(string(content), "|||", 2)

	// If the blks is less one, then return it
	if len(blks) == 1 {
		return strings.TrimSpace(blks[0]), nil
	}
	// If 0, then return an error
	if len(blks) <= 0 {
		return "", fmt.Errorf("no data in the provided file")
	}

	texts := strings.SplitN(string(content), "|||", 2)
	if fe.ismeta {
		if len(texts) >= 1 {
			text = texts[0]
		}
	} else {
		if len(texts) >= 2 {
			text = texts[1]
		}
	}
	return strings.TrimSpace(text), nil
}

func Metadata(name string) (LicenseMeta, error) {
	var lm LicenseMeta
	cont, err := extractText(ext{input: name, ismeta: true})
	if err != nil {
		return lm, err
	}

	var listI struct {
		perm bool
		cond bool
		lim  bool
	}

	lines := strings.Split(cont, "\n")
	var il bool
	var temp_v []string

	for _, line := range lines {
		var parts []string
		if strings.TrimSpace(line) == "" {
			continue
		}
		if il && !strings.ContainsRune(line, ':') {
			parts = strings.SplitN(line, "-", 2)
			if len(parts) >= 2 {
				temp_v = append(temp_v, strings.TrimSpace(parts[1]))
			}
			continue
		}

		if il && strings.ContainsRune(line, ':') {
			if listI.lim {
				lm.Limitations = temp_v
				listI.lim = false
			}

			if listI.cond {
				lm.Conditions = temp_v
				listI.cond = false
			}

			if listI.perm {
				lm.Permissions = temp_v
				listI.perm = false
			}

			temp_v = []string{}
			il = false
		}

		if !il {
			// Spit the line using ':'
			parts = strings.SplitN(line, ":", 2)
		}

		if !il && len(parts) > 0 {
			var val string
			if len(parts) >= 2 {
				val = strings.TrimSpace(parts[1])
			}
			switch strings.ToLower(strings.TrimSpace(parts[0])) {
			case "title":
				lm.Title = val
			case "spdx-id":
				lm.SpdxId = val
			case "nickname":
				lm.Nickname = val
			case "description":
				lm.Description = val
			case "note":
				lm.Note = val // This is not that common
			case "permissions":
				il = true
				listI.perm = true
			case "conditions":
				il = true
				listI.cond = true
			case "limitations":
				il = true
				listI.lim = true
			}
		}
	}

	// A clunky way to handle  End Of Input
	switch {
	case listI.lim:
		lm.Limitations = temp_v
	case listI.cond:
		lm.Conditions = temp_v
	case listI.perm:
		lm.Permissions = temp_v
	}

	return lm, nil
}
