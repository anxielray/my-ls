package internal

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	T "my-ls-1/cmd/terminal/lsOptions"
	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
	U "my-ls-1/pkg/utils"
)

func ListSingleFile(path string, options OP.Options) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		fmt.Printf("ls: cannot access '%s': %v\n", path, err)
		return
	}

	file := FI.FileInfo{
		Name:    fileInfo.Name(),
		Size:    fileInfo.Size(),
		Mode:    fileInfo.Mode(),
		ModTime: fileInfo.ModTime(),
		IsDir:   fileInfo.IsDir(),
		IsLink:  fileInfo.Mode()&os.ModeSymlink != 0,
	}

	if file.IsLink {
		linkTarget, err := os.Readlink(path)
		if err == nil {
			file.LinkTarget = linkTarget
		}
	}

	if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		file.Nlink = stat.Nlink
		file.Uid = stat.Uid
		file.Gid = stat.Gid
		file.Rdev = stat.Rdev
	}

	if options.LongFormat {
		U.PrintLongFormat([]FI.FileInfo{file}, options)
	} else {
		fmt.Println(U.FormatFileName(file, options))
	}
}

func ListDir(path string, options OP.Options) {
	files, err := T.ReadDirectory(path, options)
	if err != nil {
		fmt.Printf("ls: cannot access '%s': %v\n", path, err)
		return
	}
	U.PrintFiles(files, options)
}

// ListRecursive function to list files and directories recursively in reverse order
func ListRecursive(path string, options OP.Options) {
	var NewPath string
	if !strings.HasSuffix(path, ".") && !strings.HasSuffix(path, "..") {

		fmt.Printf("%s:\n", path)
		files, _ := T.ReadDirectory(path, options)

		if options.LongFormat {
			if options.ShowHidden {
				U.PrintLongFormat(files, options)
			} else {
				U.PrintLongFormat(FilterHidden(files), options)
			}
		} else {
			if options.ShowHidden {
			} else {
				U.PrintFiles(FilterHidden(files), options)
			}
		}

		fmt.Println()

		// open  a loop to update the path for every entry
		for _, file := range files {
			if file.IsDir {

				if strings.HasSuffix(path, "/") {
					NewPath = fmt.Sprintf("%s%s", path, file.Name)
				} else {
					NewPath = fmt.Sprintf("%s/%s", path, file.Name)
				}

				ListRecursive(NewPath, options)
			}
		}
	}
}

func FilterHidden(entries []FI.FileInfo) []FI.FileInfo {
	var filtered []FI.FileInfo
	for _, entry := range entries {
		if entry.Name[0] != '.' {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func AddSpecialEntryReverse(path, name string, files *[]FI.FileInfo) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	fileInfo := FI.CreateFileInfo(Dir(path), info)
	fileInfo.Name = name
	*files = append(*files, fileInfo)
}

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
