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

// Return the string format of the gitignore content that has been computed.
// It takes in the GitInput struct.
// If the user passes multiple `Ignores` then and an output file (optional), the following is done.
//   - Append all the file contents into a strings.Builder for next manipulations.
//     This step potentially has multiple duplicate files, e.g both `Go` and `C` have executables `*.exe`
//   - Check if the file passed as output contains anything, if it does, then:
//     add the content to as if it was got from a file.
//   - Remove the duplicate lines and then return it.
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

// Read from the 'directory' containing the .gitignore and then return the list.
// The list is in form of array string.
func AllGitIgnore() []string {
	var gits []string
	files, err := gitignores.ReadDir("gitignore")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	for _, agi := range files {
		file := strings.TrimSuffix(agi.Name(), ".gitignore")
		gits = append(gits, file)
	}

	return gits
}

// use the simple fact that go `map` doesn't accept duplicate keys.
// using this, then just 'append' the keys and remove all that are the same
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

// Given a string, get the unique lines split in 'blocks'
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

// Split a string using '#' as the section 'header'
// An example:
//
//	'''
//	  # Go files
//	  *.exe
//	  # C files
//	  *.exe
//	  *.out
//	'''
//
// This will be the split into two sections, the new lines between them is not useful as they are ignored.
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
