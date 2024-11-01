package utils

import (
	"fmt"
	"os"
	"strings"

	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
)

// This function calculates the total size in 1 KB blocks for the specified directory's entries only.
func calculateTotalBlocks(dir string, options OP.Options) (int64, error) {
	var totalBlocks int64
	var files []FI.FileInfo
	var hiddens []os.DirEntry

	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		// capture all the hidden entries
		if IsHidden(entry) {
			hiddens = append(hiddens, entry)
		}
	}

	// Read the direct entries in the specified directory
	if options.ShowHidden {
		AddSpecialEntry(dir, ".", &files)
		AddSpecialEntry(fmt.Sprintf("%s/%s", dir, ".."), "..", &files)
		for _, ety := range hiddens {
			AddSpecialEntry(fmt.Sprintf("%s/%s", dir, ety.Name()), ety.Name(), &files)
		}
	}

	for _, entry := range entries {
		if !IsHidden(entry) {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			fileInfo := FI.CreateFileInfo(dir, info)
			files = append(files, fileInfo)
		}
	}

	for _, file := range files {
		totalBlocks += file.Blocks
	}

	return totalBlocks, nil
}

func AddSpecialEntry(path, name string, files *[]FI.FileInfo) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	fileInfo := FI.CreateFileInfo(GetDir(path), info)
	fileInfo.Name = name
	*files = append(*files, fileInfo)
}

// / IsHidden checks if a given DirEntry is a hidden directory.
func IsHidden(entry os.DirEntry) bool {
	// Check if it's a directory
	if !entry.IsDir() {
		return false
	}

	// Get the base name of the directory
	dirName := entry.Name()

	// Check if the directory name starts with a dot
	return strings.HasPrefix(dirName, ".")
}

// GetDir returns the directory part of a file path
func GetDir(path string) string {
	// Trim any trailing slashes
	path = strings.TrimRight(path, "/")
	// Find the last occurrence of the slash
	lastSlashIndex := strings.LastIndex(path, "/")
	if lastSlashIndex == -1 {
		return "." // No directory part, return current directory
	}
	// Return the substring up to the last slash
	return path[:lastSlashIndex]
}
