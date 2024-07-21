package src

import "os"

// A wrapper around the standard os.WriteFile
func FileWrite(content, file string) error {
	return os.WriteFile(file, []byte(content), os.ModePerm)
}
