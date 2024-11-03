package options

import (
	"fmt"
	"os"
	"strings"
)

type Options struct {
	LongFormat bool //-l
	Recursive  bool // -R
	ShowHidden bool // -a
	Reverse    bool // -r
	SortByTime bool // -t
	SortBySize bool // -S
	OnePerLine bool // -1
	NoColor    bool
}

/*
The function will collect the command line arguments and sort them into flags
and files. Declares each flag as well. If the option is -- the functioin will
treat it as current directory. The function returns options boolean values
as well as the array containing the files/directories to be displayed by the ls command.
*/
func ParseFlags() (Options, []string) {
	var options Options

	args := os.Args[1:]

	var dirs []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") && len(arg) > 1 && arg != "--" {
			for _, flag := range arg[1:] {
				switch flag {
				case 'l':
					options.LongFormat = true
				case 'R':
					options.Recursive = true
				case 'a':
					options.ShowHidden = true
				case 'r':
					options.Reverse = true
				case 't':
					options.SortByTime = true
				case 'S':
					options.SortBySize = true
				case '1':
					options.OnePerLine = true
				case 'G':
					options.NoColor = true
				default:
					fmt.Printf("ls: invalid option -- '%c'\n", flag)
					os.Exit(1)
				}
			}
		} else if arg == "--" {
			dirs = append(dirs, args[i+1:]...)
			break
		} else {
			if strings.HasSuffix(arg, "/") {
				IsNormalFile(arg)
			}
			dirs = append(dirs, arg)
		}
	}

	return options, dirs
}

func LinkMessage(path string) {
	fmt.Printf("ls: cannot access '%s': Not a directory\n", path)
}

func IsNormalFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		LinkMessage(path)
		os.Exit(0)
	}
	return !fileInfo.IsDir()
}
