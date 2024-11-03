package lsOptions

import (
	"fmt"
	"os"
	"strings"

	S "my-ls-1/internal/sort"
	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
)

/*
This function will take the path and the optioins issuedon the command line
and processes and returns a slice of fileinfos from the path entries
*/
func ReadDirectory(path string, options OP.Options) ([]FI.FileInfo, error) {
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

	S.SortFiles(files, options)

	return files, nil
}

/*
When the options of show all is set to true, the function will be called
to add the entries of the current directory and the parent directory
*/
func AddSpecialEntry(path, name string, files *[]FI.FileInfo) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	fileInfo := FI.CreateFileInfo(Dir(path), info)
	fileInfo.Name = name
	*files = append([]FI.FileInfo{fileInfo}, *files...)
}

/*
The Dir function is designed to return the directory portion of a given
file path. It processes the input path string and extracts the directory
component, taking care of various edge cases
*/
func Dir(path string) string {

	if path == "" {
		return "."
	}

	path = strings.TrimRight(path, string(os.PathSeparator))

	i := strings.LastIndex(path, string(os.PathSeparator))

	if i == -1 {
		return "."
	}

	if i == 0 {
		return string(os.PathSeparator)
	}

	return path[:i]
}
