package utils

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	T "my-ls-1/cmd/terminal"
	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
)

var Path string

// This function will print entries in the long format. (ls -l)
func PrintLongFormat(files []FI.FileInfo, options OP.Options) {

	if len(os.Args) > 2 {

		//check if a file has been passed
		if len(files) > 1 || options.Recursive {

			// check for the path to display size
			for _, arg := range os.Args[1:] {
				if strings.HasPrefix(arg, "-") {
					continue
				} else {

					var path string
					file, _ := checkPathType(arg)
					if file == "symlink" {
						wd, _ := os.Getwd()
						Path = fmt.Sprintf("%s/%s", wd, arg)
						break
					}
					if file == "directory" {
						if isStandardLibrary(arg) {
							path = arg
						} else {
							PATH, _ := os.Getwd()
							path = fmt.Sprintf("%s/%s", PATH, arg)
							Path = path
						}

						totalBlocks, _ := calculateTotalBlocks(path, options)
						fmt.Printf("total %d\n", totalBlocks)
						break
					}
				}
			}
		}

	} else if len(os.Args) == 2 {
		var path string
		if isStandardLibrary(os.Args[1]) {
			path = os.Args[1]
		} else {
			PATH, _ := os.Getwd()
			path = PATH
			Path = path
		}
		totalBlocks, _ := calculateTotalBlocks(path, options)
		fmt.Printf("total %d\n", totalBlocks)
	}

	maxNlinkWidth := 0
	maxUserWidth := 0
	maxGroupWidth := 0
	maxSizeWidth := 0
	maxMajorWidth := 0
	maxMinorWidth := 0

	// printing the number of hardlinks of a specific file.
	for _, file := range files {

		nlinkWidth := len(fmt.Sprintf("%d", file.Nlink))
		if nlinkWidth > maxNlinkWidth {
			maxNlinkWidth = nlinkWidth
		}

		usr, _ := user.LookupId(fmt.Sprint(file.Uid))
		userWidth := len(usr.Username)
		if userWidth > maxUserWidth {
			maxUserWidth = userWidth
		}

		grp, _ := user.LookupGroupId(fmt.Sprint(file.Gid))
		groupWidth := len(grp.Name)
		if groupWidth > maxGroupWidth {
			maxGroupWidth = groupWidth
		}

		if file.Mode&os.ModeDevice != 0 {
			major := Major(file.Rdev)
			minor := Minor(file.Rdev)
			majorWidth := len(fmt.Sprintf("%d", major))
			minorWidth := len(fmt.Sprintf("%d", minor))
			if majorWidth > maxMajorWidth {
				maxMajorWidth = majorWidth
			}
			if minorWidth > maxMinorWidth {
				maxMinorWidth = minorWidth
			}
		} else {
			sizeWidth := len(fmt.Sprintf("%d", file.Size))
			if sizeWidth > maxSizeWidth {
				maxSizeWidth = sizeWidth
			}
		}
	}
	if val, _ := IsSymlink(Path); val {
		files, _ = GetSymlinksInDir(fmt.Sprintf("%s/..", Path))
		files = filterSymLink(files)
	}

	for _, file := range files {
		usr, _ := user.LookupId(fmt.Sprint(file.Uid))
		grp, _ := user.LookupGroupId(fmt.Sprint(file.Gid))

		modeStr := FormatFileMode(file.Mode)

		size := ""
		if file.Mode&os.ModeDevice != 0 {
			major := Major(file.Rdev)
			minor := Minor(file.Rdev)
			size = fmt.Sprintf("%*d, %*d", maxMajorWidth, major, maxMinorWidth, minor) // this is a device
		} else {
			size = fmt.Sprintf("%*d", maxSizeWidth, file.Size) // normal directory
		}

		fileName := FormatFileName(file, options)

		timeFormat := "Jan _2 15:04"
		sixMonthsAgo := time.Now().AddDate(0, -6, 0)
		if file.ModTime.Before(sixMonthsAgo) {
			timeFormat = "Jan _2  2006"
		}

		fmt.Printf("%s %*d %-*s %-*s %*s %s %s\n",
			modeStr,
			maxNlinkWidth, file.Nlink,
			maxUserWidth, usr.Username,
			maxGroupWidth, grp.Name,
			maxSizeWidth+maxMajorWidth+maxMinorWidth, size,
			file.ModTime.Format(timeFormat),
			fileName,
		)
	}
}

func filterSymLink(files []FI.FileInfo) []FI.FileInfo {
	link := os.Args[len(os.Args)-1]
	for _, fi := range files {
		if fi.Name == link {
			return []FI.FileInfo{fi}
		}
	}
	return files
}

//This function will format the files in the terminal correctly, based on the column width
func PrintColumnar(files []FI.FileInfo, options OP.Options) {
	termWidth := T.GetTerminalWidth()
	if termWidth < 1 {
		termWidth = 80
	}

	maxWidth := 0
	for _, file := range files {
		width := len(FormatFileName(file, options))
		if width > maxWidth {
			maxWidth = width
		}
	}

	colWidth := maxWidth + 2
	numCols := termWidth / colWidth
	if numCols == 0 {
		numCols = 1
	}

	numRows := (len(files) + numCols - 1) / numCols

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			idx := j*numRows + i
			if idx < len(files) {
				fileName := FormatFileName(files[idx], options)
				fmt.Printf("%-*s", colWidth, fileName)
			}
		}
		fmt.Println()
	}
}

// This functions lists entries to the console based on the option long format
func PrintFiles(files []FI.FileInfo, options OP.Options) {
	if options.LongFormat {
		PrintLongFormat(files, options)
	} else if options.OnePerLine {
		for _, file := range files {
			fmt.Println(FormatFileName(file, options))
		}
	} else {
		PrintColumnar(files, options)
	}
}

/*
Extract major and minor device numbers from a uintptr,
which typically represents the underlying system's device information.
*/
func Major(dev uint64) uint64 {
	return (dev >> 8) & 0xFF
}

func Minor(dev uint64) uint64 {
	return dev & 0xFF
}

// isStandardLibrary checks if the given path is a standard library directory.
func isStandardLibrary(path string) bool {
	standardLibs := []string{"/usr/bin", "/etc", "/dev", "/usr/lib", "/usr/local/bin", "/bin", "/sbin"}

	for _, lib := range standardLibs {
		if CleanPath(path) == lib {
			return true
		}
	}
	return false
}

// CleanPath normalizes a given path by removing redundant elements like "." and "..".
func CleanPath(path string) string {

	if path == "" {
		return "."
	}

	parts := strings.Split(path, "/")
	var cleanedParts []string

	for _, part := range parts {
		switch part {
		case "":
			continue
		case ".":
			continue
		case "..":
			if len(cleanedParts) > 0 {
				cleanedParts = cleanedParts[:len(cleanedParts)-1]
			}
		default:
			cleanedParts = append(cleanedParts, part)
		}
	}

	cleanedPath := strings.Join(cleanedParts, "/")

	if strings.HasPrefix(path, "/") {
		cleanedPath = "/" + cleanedPath
	}

	if cleanedPath == "" {
		return "."
	}

	return cleanedPath
}
