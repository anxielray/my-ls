package utils

import (
	"fmt"
	"math"
	"os"
	"os/user"
	"strings"
	"time"

	T "my-ls-1/cmd/terminal"
	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"

	"golang.org/x/sys/unix"
)

func PrintLongFormat(files []FI.FileInfo, options OP.Options) {
	if len(os.Args) > 2 {
		//check for the path to display size
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-") {
				continue
			} else {
				file, _ := checkPathType(arg)
				if file == "directory" {
					totalBlocks, _ := calculateTotalBlocks(arg)
					fmt.Printf("total %d\n", totalBlocks/2)
				}
			}
		}
	} else if len(os.Args) == 2 {
		println("Checking")
		path, _ := os.Getwd()
		totalBlocks, _ := calculateTotalBlocks(path)
		fmt.Printf("total %d\n", totalBlocks/2)
	}

	maxNlinkWidth := 0
	maxUserWidth := 0
	maxGroupWidth := 0
	maxSizeWidth := 0
	maxMajorWidth := 0
	maxMinorWidth := 0

	//printing the number of hardlinks of a specific file.
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
			major := unix.Major(file.Rdev)
			minor := unix.Minor(file.Rdev)
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
			major := unix.Major(file.Rdev)
			minor := unix.Minor(file.Rdev)
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
		PrintLongFormat(files, options)
	} else if options.OnePerLine {
		for _, file := range files {
			fmt.Println(FormatFileName(file, options))
		}
	} else {
		PrintColumnar(files, options)
	}
}
