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

//The function to list files and directories recursively
func ListRecursive(path string, options OP.Options) {
	if path == "." {
		var NewPath string

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
				U.PrintFiles(files, options)
			} else {
				U.PrintFiles(FilterHidden(files), options)
			}
		}

		fmt.Println()

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
	if !strings.HasSuffix(path, ".") && !strings.HasSuffix(path, "..") {

		var NewPath string
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
				U.PrintFiles(files, options)
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

//The function will eliminate the directories and file that start in a period(.)
func FilterHidden(entries []FI.FileInfo) []FI.FileInfo {
	var filtered []FI.FileInfo
	for _, entry := range entries {
		if entry.Name[0] != '.' {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
