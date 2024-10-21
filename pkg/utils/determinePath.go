package utils

import (
	"fmt"
	"os"
)

// checkPathType determines if the provided path is a file or a directory.
func checkPathType(path string) (string, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("path does not exist: %s", path)
	} else if err != nil {
		return "", err
	}

	if info.IsDir() {
		return "directory", nil
	} else {
		return "file", nil
	}
}
