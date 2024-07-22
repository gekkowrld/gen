package src

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
)

type GitInput struct {
	Ignores []string
	IsInput bool
	Output  string
}

//go:embed gitignore/*.gitignore
var gitignores embed.FS

func GitIgnore(ig GitInput) string {
	var gitignore strings.Builder
	ignores := ig.Ignores

	for _, ignore := range ignores {
		cont, err := gitignores.ReadFile(fmt.Sprintf("gitignore/%s.gitignore", strings.ToLower(ignore)))
		if err != nil {
			log.Fatalf("%+v\n", err)
		}

		_, err = gitignore.Write(cont)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
	}

	// Add .gitignore file contents
	if ig.IsInput {
		fc, _ := os.ReadFile(ig.Output)
		gitignore.WriteRune('\n')
		gitignore.Write(fc)
		gitignore.WriteRune('\n')
	}

	return strings.TrimSpace(Unique(gitignore.String()))
}

func AllGitIgnore() {
	files, err := gitignores.ReadDir("gitignore")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	var __ag strings.Builder

	for _, agi := range files {
		file := strings.TrimSuffix(agi.Name(), ".gitignore")
		__ag.WriteString(file + "\n")
	}

	fmt.Println(__ag.String())
}

func uniq(input string) string {
	var uniqueLines strings.Builder
	uLines := make(map[string]bool)

	lines := strings.Split(input, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if _, ok := uLines[trimmedLine]; ok {
			continue
		}

		uLines[trimmedLine] = true
		uniqueLines.WriteString(line + "\n")
	}

	return uniqueLines.String()
}

func Unique(input string) string {
	var str strings.Builder
	// Get the blocks
	splits := SplitSec(input)
	for _, split := range splits {
		split = strings.TrimSpace(split)
		if len(strings.Split(split, "\n")) > 1 {
			str.WriteString(split + "\n\n")
		}
	}

	return str.String()
}

func SplitSec(input string) []string {
	var s_sec []string

	var __bstr strings.Builder
	lines := strings.Split(uniq(input), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			if __bstr.Len() > 1 {
				s_sec = append(s_sec, __bstr.String())
				__bstr.Reset()
			}
		}
		__bstr.WriteString(line + "\n")
	}
	if __bstr.Len() > 0 {
		s_sec = append(s_sec, __bstr.String())
	}

	return s_sec
}
