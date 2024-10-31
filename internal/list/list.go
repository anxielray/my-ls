package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

func ListRecursive(path string, options OP.Options) {
	fmt.Printf("%s:\n", path)
	filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("ls: cannot access '%s': %v\n", p, err)
			return nil
		}
		if d.IsDir() {
			if options.ShowHidden {
				if p != path {
					fmt.Printf("\n%s:\n", p)
				}
				files, err := T.ReadDirectory(p, options)
				if err != nil {
					fmt.Printf("ls: cannot access '%s': %v\n", p, err)
					return nil
				}
				U.PrintFiles(files, options)
			} else {
				if !strings.HasPrefix(p, ".") {
					if p != path {
						fmt.Printf("\n./%s:\n", p)
					}
					files, err := T.ReadDirectory(p, options)
					if err != nil {
						fmt.Printf("ls: cannot access '%s': %v\n", p, err)
						return nil
					}
					files = FilterHidden(files)
					U.PrintFiles(files, options)
				} else if p == "." {

					files, err := T.ReadDirectory(p, options)
					if err != nil {
						fmt.Printf("ls: cannot access '%s': %v\n", p, err)
						return nil
					}

					files = FilterHidden(files)
					U.PrintFiles(files, options)
				}
			}
		}
		return nil
	})
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
