package src

import "os"

// A thin wrapper around the standard os.WriteFile
func FileWrite(content, file string) error {
	return os.WriteFile(file, []byte(content), 0644)
}
