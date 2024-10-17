package main

import (
	"fmt"
	"os"

	L "my-ls-1/internal/list"
	OP "my-ls-1/pkg/options"
	C "my-ls-1/pkg/utils/color"
)

func main() {
	C.InitColorMap()
	options, args := OP.ParseFlags()

	if len(args) == 0 {
		args = []string{"."}
	}

	for i, arg := range args {
		if len(args) > 1 {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("%s:\n", arg)
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
