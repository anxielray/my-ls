package utils

import (
	"fmt"
	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
	"os"
	"path/filepath"
	"strings"
)

// const blockSize int = 512 // Default block size in bytes

// This function calculates the total size in 1 KB blocks for the specified directory's entries only.
func calculateTotalBlocks(dir string, options OP.Options) (int64, error) {
	var totalBlocks int64
	var files []FI.FileInfo

	// Read the direct entries in the specified directory
	// Add current (.) and parent (..) directory entries
	if options.ShowHidden {
		AddSpecialEntry(dir, ".", &files)
		AddSpecialEntry(fmt.Sprintf("%s/%s", dir, ".."), "..", &files)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		// Skip hidden entries if options.ShowHidden is false
		// Retrieve FileInfo for each entry and add it to files
		if !options.ShowHidden && !IsHidden(entry) {

			info, err := entry.Info()
			if err != nil {
				return 0, err
			}
			fileInfo := FI.CreateFileInfo(dir, info)
			files = append(files, fileInfo)
		}

	}

	// Calculate total blocks based on the `files` slice, including "." and ".."
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
	fileInfo := FI.CreateFileInfo(filepath.Dir(path), info)
	fileInfo.Name = name
	*files = append(*files, fileInfo)
}

/// IsHidden checks if a given DirEntry is a hidden directory.
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
