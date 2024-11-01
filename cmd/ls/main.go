package main

import (
	"fmt"
	"os"

	L "my-ls-1/internal/list"
	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
	C "my-ls-1/pkg/utils/color"
)

func main() {

	//Initialize color function
	C.InitColorMap()

	//Parse command line flags and arguments
	options, args := OP.ParseFlags()

	if len(args) == 0 {
		args = []string{"."}
	}
	args, _ = AddFullPathAndSort(args)
	for i, arg := range args {
		if len(args) > 1 {
			if i > 0 {
				fmt.Println()
			}
			filIf, _ := os.Stat(arg)
			if FI.CreateFileInfo(arg, filIf); filIf.IsDir() {
				fmt.Printf("%s:\n", arg)
			}
		}

		fileInfo, err := os.Stat(arg)
		if err != nil {
			fmt.Printf("ls: cannot access '%s': %v\n", arg, err)
			continue
		}

		if fileInfo.IsDir() {
			if options.Recursive {
				L.ListRecursive(arg, options)
			} else {
				L.ListDir(arg, options)
			}
		} else {
			L.ListSingleFile(arg, options)
		}
	}
}

func AddFullPathAndSort(shortPaths []string) ([]string, error) {

	// Separate files and directories
	var files, dirs []string
	for _, path := range shortPaths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
	}

	// Combine files and directories
	return append(files, dirs...), nil
}
