package lsOptions

import (
	"fmt"
	"os"
	"path/filepath"
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

	S.SortFiles(files, options)
	return files, nil
}

func AddSpecialEntry(path, name string, files *[]FI.FileInfo) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	fileInfo := FI.CreateFileInfo(filepath.Dir(path), info)
	fileInfo.Name = name
	*files = append([]FI.FileInfo{fileInfo}, *files...)
}
