package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unsafe"

	"golang.org/x/sys/unix"
)

type Options struct {
	LongFormat bool
	Recursive  bool
	ShowHidden bool
	Reverse    bool
	SortByTime bool
	SortBySize bool
	OnePerLine bool
	NoColor    bool
}

type FileInfo struct {
	Name       string
	Size       int64
	Mode       os.FileMode
	ModTime    time.Time
	IsDir      bool
	Nlink      uint64
	Uid        uint32
	Gid        uint32
	IsLink     bool
	LinkTarget string
	Rdev       uint64
	Blocks     int64
}

const (
	ColorReset = "\033[0m"
)

var colorMap map[string]string

func initColorMap() {
	colorMap = make(map[string]string)
	lsColors := os.Getenv("LS_COLORS")
	if lsColors == "" {
		return
	}

	pairs := strings.Split(lsColors, ":")
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			colorMap[parts[0]] = "\033[" + parts[1] + "m"
		}
	}
}

func main() {
	initColorMap()
	options, args := parseFlags()

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
				listRecursive(arg, options)
			} else {
				listDir(arg, options)
			}
		} else {
			listSingleFile(arg, options)
		}
	}
}

func listSingleFile(path string, options Options) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		fmt.Printf("ls: cannot access '%s': %v\n", path, err)
		return
	}

	file := FileInfo{
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
		printLongFormat([]FileInfo{file}, options)
	} else {
		fmt.Println(formatFileName(file, options))
	}
}

func readDir(path string, options Options) ([]FileInfo, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	entries, err := dir.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	files := make([]FileInfo, 0, len(entries))

	if options.ShowHidden {
		addSpecialEntry(path, ".", &files)
		parentPath := filepath.Dir(path)
		addSpecialEntry(parentPath, "..", &files)
	}

	for _, entry := range entries {
		if !options.ShowHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		fileInfo := createFileInfo(path, info)
		files = append(files, fileInfo)
	}

	sortFiles(files, options)
	return files, nil
}

func createFileInfo(path string, info os.FileInfo) FileInfo {
	fileInfo := FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
		IsLink:  info.Mode()&os.ModeSymlink != 0,
	}

	if fileInfo.IsLink {
		linkTarget, err := os.Readlink(filepath.Join(path, info.Name()))
		if err == nil {
			fileInfo.LinkTarget = linkTarget
		}
	}

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		fileInfo.Nlink = stat.Nlink
		fileInfo.Uid = stat.Uid
		fileInfo.Gid = stat.Gid
		fileInfo.Rdev = stat.Rdev
		fileInfo.Blocks = stat.Blocks
	}

	return fileInfo
}

func addSpecialEntry(path, name string, files *[]FileInfo) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	fileInfo := createFileInfo(filepath.Dir(path), info)
	fileInfo.Name = name
	*files = append(*files, fileInfo)
}

func sortFiles(files []FileInfo, options Options) {
	sort.Slice(files, func(i, j int) bool {
		if options.SortByTime {
			if !files[i].ModTime.Equal(files[j].ModTime) {
				return files[i].ModTime.After(files[j].ModTime)
			}
		} else if options.SortBySize {
			if files[i].Size != files[j].Size {
				return files[i].Size > files[j].Size
			}
		}

		return compareFilenamesAlphanumeric(files[i].Name, files[j].Name)
	})

	if options.Reverse {
		reverseSlice(files)
	}
}

func compareFilenamesAlphanumeric(a, b string) bool {
	aRunes := []rune(a)
	bRunes := []rune(b)
	aLen := len(aRunes)
	bLen := len(bRunes)

	for i, j := 0, 0; i < aLen && j < bLen; {
		// Skip non-alphanumeric characters
		for i < aLen && !isAlphanumeric(aRunes[i]) {
			i++
		}
		for j < bLen && !isAlphanumeric(bRunes[j]) {
			j++
		}

		// If we've reached the end of either string, compare lengths
		if i == aLen || j == bLen {
			return aLen < bLen
		}

		// If both characters are digits, compare the whole number
		if unicode.IsDigit(aRunes[i]) && unicode.IsDigit(bRunes[j]) {
			aNum, aEnd := extractNumber(aRunes[i:])
			bNum, bEnd := extractNumber(bRunes[j:])

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

func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func extractNumber(runes []rune) (int, int) {
	num := 0
	i := 0
	for i < len(runes) && unicode.IsDigit(runes[i]) {
		digit, _ := strconv.Atoi(string(runes[i]))
		num = num*10 + digit
		i++
	}
	return num, i
}

func reverseSlice(slice []FileInfo) {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - 1 - i
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func formatFileName(file FileInfo, options Options) string {
	name := file.Name
	if !options.NoColor {
		name = colorize(file, name)
	}
	if file.IsLink {
		name += " -> " + file.LinkTarget
	}
	return name
}

func colorize(file FileInfo, name string) string {
	var colorCode string

	if file.IsDir {
		colorCode = colorMap["di"]
	} else if file.IsLink {
		colorCode = colorMap["ln"]
	} else if file.Mode&0o111 != 0 {
		colorCode = colorMap["ex"]
	} else if file.Mode&os.ModeNamedPipe != 0 {
		colorCode = colorMap["pi"]
	} else if file.Mode&os.ModeSocket != 0 {
		colorCode = colorMap["so"]
	} else if file.Mode&os.ModeDevice != 0 {
		colorCode = colorMap["bd"]
	} else {
		ext := filepath.Ext(name)
		if ext != "" {
			colorCode = colorMap["*"+ext]
		}
	}

	if colorCode == "" {
		return name
	}

	return colorCode + name + ColorReset
}

func parseFlags() (Options, []string) {
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
			dirs = append(dirs, arg)
		}
	}

	return options, dirs
}

func listRecursive(path string, options Options) {
	filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("ls: cannot access '%s': %v\n", p, err)
			return nil
		}
		if d.IsDir() {
			if p != path {
				fmt.Printf("\n%s:\n", p)
			}
			files, err := readDir(p, options)
			if err != nil {
				fmt.Printf("ls: cannot access '%s': %v\n", p, err)
				return nil
			}
			printFiles(files, options)
		}
		return nil
	})
}

func listDir(path string, options Options) {
	files, err := readDir(path, options)
	if err != nil {
		fmt.Printf("ls: cannot access '%s': %v\n", path, err)
		return
	}
	printFiles(files, options)
}

func printFiles(files []FileInfo, options Options) {
	if options.LongFormat {
		printLongFormat(files, options)
	} else if options.OnePerLine {
		for _, file := range files {
			fmt.Println(formatFileName(file, options))
		}
	} else {
		printColumnar(files, options)
	}
}

func formatFileMode(mode os.FileMode) string {
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
	result.WriteString(formatPermissions(mode))

	return result.String()
}

func formatPermissions(mode os.FileMode) string {
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

func printLongFormat(files []FileInfo, options Options) {
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

		modeStr := formatFileMode(file.Mode)

		size := ""
		if file.Mode&os.ModeDevice != 0 {
			major := unix.Major(file.Rdev)
			minor := unix.Minor(file.Rdev)
			size = fmt.Sprintf("%*d, %*d", maxMajorWidth, major, maxMinorWidth, minor)
		} else {
			size = fmt.Sprintf("%*d", maxSizeWidth, file.Size)
		}

		fileName := formatFileName(file, options)

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

func printColumnar(files []FileInfo, options Options) {
	termWidth := getTerminalWidth()

	maxWidth := 0
	for _, file := range files {
		width := len(formatFileName(file, options))
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
				fmt.Printf("%-*s", colWidth, formatFileName(files[idx], options))
			}
		}
		fmt.Println()
	}
}

func getTerminalWidth() int {
	defaultWidth := 80

	var size [4]uint16
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&size))); err == 0 {
		return int(size[1])
	}

	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil {
			return width
		}
	}

	return defaultWidth
}
