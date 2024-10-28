package utils

import (
	"fmt"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	T "my-ls-1/cmd/terminal"
	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
)

// implementation of the Unix command (ls -l)
func PrintLongFormat(files []FI.FileInfo, options OP.Options) {
	if len(os.Args) > 2 { // path is declared
		// check for the path to display size
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-") {
				continue
			} else {
				file, _ := checkPathType(arg)
				if file == "directory" {
					var path string
					if isStandardLibrary(arg) {
						path = arg
					} else {
						PATH, _ := os.Getwd()
						path = fmt.Sprintf("%s/%s", PATH, arg)
					}

					totalBlocks, _ := calculateTotalBlocks(path, options)
					fmt.Printf("total %d\n", totalBlocks)
				}
			}
		}
	} else if len(os.Args) == 2 {
		var path string
		if isStandardLibrary(os.Args[1]) {
			path = os.Args[1]
		} else {
			PATH, _ := os.Getwd()
			path = fmt.Sprintf("%s/%s", PATH, os.Args[1])
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

	for _, file := range files {
		usr, _ := user.LookupId(fmt.Sprint(file.Uid))
		grp, _ := user.LookupGroupId(fmt.Sprint(file.Gid))

		modeStr := FormatFileMode(file.Mode)

		size := ""
		if file.Mode&os.ModeDevice != 0 {
			major := Major(file.Rdev)
			minor := Minor(file.Rdev)
			size = fmt.Sprintf("%*d, %*d", maxMajorWidth, major, maxMinorWidth, minor)
		} else {
			size = fmt.Sprintf("%*d", maxSizeWidth, file.Size)
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

func PrintColumnar(files []FI.FileInfo, options OP.Options) {
	termWidth := T.GetTerminalWidth()

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

	numRows := int(math.Ceil(float64(len(files)) / float64(numCols)))

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			idx := j*numRows + i
			if idx < len(files) {
				fmt.Printf("%-*s", colWidth, FormatFileName(files[idx], options))
			}
		}
		fmt.Println()
	}
}

func PrintFiles(files []FI.FileInfo, options OP.Options) {
	if options.LongFormat {
		// fullFiles := AddCurrentAndParentDirectory(files)
		PrintLongFormat(files, options)
	} else if options.OnePerLine {
		for _, file := range files {
			fmt.Println(FormatFileName(file, options))
		}
	} else {
		PrintColumnar(files, options)
	}
}

// implement the third party package of unix.
func Major(dev uint64) uint64 {
	return (dev >> 32) & 0xFFFFFFFF
}

func Minor(dev uint64) uint64 {
	return dev & 0xFFFFFFFF
}

// isStandardLibrary checks if the given path is a standard library directory.
func isStandardLibrary(path string) bool {
	standardLibs := []string{"/usr/bin", "/etc", "/dev", "/usr/lib", "/usr/local/bin", "/bin", "/sbin"}

	for _, lib := range standardLibs {
		if filepath.Clean(path) == lib {
			return true
		}
	}
	return false
}
