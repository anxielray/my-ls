package lsOptions

import (
	"fmt"
	"os"
	"strings"

	S "my-ls-1/internal/sort"
	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
)

func ReadDir(path string, options OP.Options) ([]FI.FileInfo, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	entries, err := dir.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	files := make([]FI.FileInfo, 0, len(entries))

	if options.ShowHidden {
		parentPath := fmt.Sprintf("%s/..", path)
		AddSpecialEntry(parentPath, "..", &files)
		AddSpecialEntry(path, ".", &files)
	}

	for _, entry := range entries {
		if !options.ShowHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		fileInfo := FI.CreateFileInfo(path, info)
		files = append(files, fileInfo)
	}

	//we would only want to sort the entries from the second index if the option of hidden is true
	if options.ShowHidden {
		S.SortFiles(files[2:], options)
	} else {
		S.SortFiles(files, options)
	}

	return files, nil
}

func AddSpecialEntry(path, name string, files *[]FI.FileInfo) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	fileInfo := FI.CreateFileInfo(Dir(path), info)
	fileInfo.Name = name
	*files = append([]FI.FileInfo{fileInfo}, *files...)
}

// Implementation of the filePath.Dir(path) function
func Dir(path string) string {
	// Handle empty path
	if path == "" {
		return "."
	}

	// Remove trailing slashes
	path = strings.TrimRight(path, string(os.PathSeparator))

	// Find last separator
	i := strings.LastIndex(path, string(os.PathSeparator))

	if i == -1 {
		// No separator found
		return "."
	}

	if i == 0 {
		// Path starts with separator (root directory)
		return string(os.PathSeparator)
	}

	// Return everything before last separator
	return path[:i]
}
