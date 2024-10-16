package utils

import (
	"fmt"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	T "my-ls-1/cmd/terminal"
	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
	C "my-ls-1/pkg/utils/color"

	"golang.org/x/sys/unix"
)

func FormatFileName(file FI.FileInfo, options OP.Options) string {
	name := file.Name
	if !options.NoColor {
		name = C.Colorize(file, name)
	}
	if file.IsLink {
		name += " -> " + file.LinkTarget
	}
	return name
}

func FormatPermissions(mode os.FileMode) string {
	const rwx = "rwxrwxrwx"
	perm := []byte("---------")

	for i := 0; i < 9; i++ {
		if mode&(1<<uint(8-i)) != 0 {
			perm[i] = rwx[i]
		}
	}

	if mode&os.ModeSetuid != 0 {
		if perm[2] == 'x' {
			perm[2] = 's'
		} else {
			perm[2] = 'S'
		}
	}
	if mode&os.ModeSetgid != 0 {
		if perm[5] == 'x' {
			perm[5] = 's'
		} else {
			perm[5] = 'S'
		}
	}
	if mode&os.ModeSticky != 0 {
		if perm[8] == 'x' {
			perm[8] = 't'
		} else {
			perm[8] = 'T'
		}
	}

	return string(perm)
}

func FormatFileMode(mode os.FileMode) string {
	var result strings.Builder

	// File type
	switch {
	case mode&os.ModeDir != 0:
		result.WriteRune('d')
	case mode&os.ModeSymlink != 0:
		result.WriteRune('l')
	case mode&os.ModeDevice != 0:
		if mode&os.ModeCharDevice != 0 {
			result.WriteRune('c')
		} else {
			result.WriteRune('b')
		}
	case mode&os.ModeNamedPipe != 0:
		result.WriteRune('p')
	case mode&os.ModeSocket != 0:
		result.WriteRune('s')
	default:
		result.WriteRune('-')
	}

	// Permission bits
	result.WriteString(FormatPermissions(mode))

	return result.String()
}

func PrintLongFormat(files []FI.FileInfo, options OP.Options) {
	var totalBlocks int64
	for _, file := range files {
		totalBlocks += file.Blocks
	}
	fmt.Printf("total %d\n", totalBlocks/2)

	maxNlinkWidth := 0
	maxUserWidth := 0
	maxGroupWidth := 0
	maxSizeWidth := 0
	maxMajorWidth := 0
	maxMinorWidth := 0

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

func IsAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func ExtractNumber(runes []rune) (int, int) {
	num := 0
	i := 0
	for i < len(runes) && unicode.IsDigit(runes[i]) {
		digit, _ := strconv.Atoi(string(runes[i]))
		num = num*10 + digit
		i++
	}
	return num, i
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

func AddSpecialEntry(path, name string, files *[]FI.FileInfo) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	fileInfo := FI.CreateFileInfo(filepath.Dir(path), info)
	fileInfo.Name = name
	*files = append(*files, fileInfo)
}

func CompareFilenamesAlphanumeric(a, b string) bool {
	aRunes := []rune(a)
	bRunes := []rune(b)
	aLen := len(aRunes)
	bLen := len(bRunes)

	for i, j := 0, 0; i < aLen && j < bLen; {
		// Skip non-alphanumeric characters
		for i < aLen && !IsAlphanumeric(aRunes[i]) {
			i++
		}
		for j < bLen && !IsAlphanumeric(bRunes[j]) {
			j++
		}

		// If we've reached the end of either string, compare lengths
		if i == aLen || j == bLen {
			return aLen < bLen
		}

		// If both characters are digits, compare the whole number
		if unicode.IsDigit(aRunes[i]) && unicode.IsDigit(bRunes[j]) {
			aNum, aEnd := ExtractNumber(aRunes[i:])
			bNum, bEnd := ExtractNumber(bRunes[j:])

			if aNum != bNum {
				return aNum < bNum
			}

			i += aEnd
			j += bEnd
		} else {
			// Compare characters case-insensitively
			aLower := unicode.ToLower(aRunes[i])
			bLower := unicode.ToLower(bRunes[j])
			if aLower != bLower {
				return aLower < bLower
			}
			i++
			j++
		}
	}

	// If all compared characters are the same, shorter string comes first
	return aLen < bLen
}
