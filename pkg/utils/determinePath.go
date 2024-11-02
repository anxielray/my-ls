package utils

import (
	"fmt"
	"os"

	FI "my-ls-1/pkg/fileinfo"
)

// checkPathType determines if the provided path is a file or a directory.
func checkPathType(path string) (string, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("path does not exist: %s", path)
	} else if err != nil {
		return "", err
	}

	if value, _ := IsSymlink(path); value {
		return "symlink", nil
	}

	if info.IsDir() {
		return "directory", nil
	} else {
		return "file", nil
	}
}

//checks if thep path is a symlink
func IsSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}

//Gets all the symlinks in a given path
func GetSymlinksInDir(dirPath string) ([]FI.FileInfo, error) {
	var symlinks []FI.FileInfo

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		inf, err := entry.Info()
		if err != nil {
			continue
		}
		fileInfo := FI.CreateFileInfo(dirPath, inf)
		if fileInfo.IsLink {
			symlinks = append(symlinks, fileInfo)
		}
	}

	return symlinks, nil
}
