package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term" // Third party library for displaying terminal information
)

type Options struct {
	LongFormat bool
	Recursive  bool
	ShowHidden bool
	Reverse    bool
	SortByTime bool
}

type FileInfo struct {
	Name      string
	Size      int64
	Mode      fs.FileMode
	ModTime   time.Time
	IsDir     bool
	User      string
	Group     string
	Links     int
	BlockSize int64
}

func main() {
	opts := parseFlags()
	args := flag.Args()

	if len(args) == 0 {
		args = []string{"."}
	}

	for _, path := range args {
		err := listDirectory(path, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}
}

func parseFlags() Options {
	var opts Options
	flag.BoolVar(&opts.LongFormat, "l", false, "Use long listing format")
	flag.BoolVar(&opts.Recursive, "R", false, "List subdirectories recursively")
	flag.BoolVar(&opts.ShowHidden, "a", false, "Do not ignore entries starting with .")
	flag.BoolVar(&opts.Reverse, "r", false, "Reverse order while sorting")
	flag.BoolVar(&opts.SortByTime, "t", false, "Sort by modification time, newest first")

	flag.Parse()

	return opts
}

func listDirectory(path string, opts Options) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var fileInfos []FileInfo
	var totalBlocks int64

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}

		// if !opts.ShowHidden && strings.HasPrefix(info.Name(), ".") && info.Name() != "." && info.Name() != ".." {
		// 	continue
		// }
		if !opts.ShowHidden && strings.HasPrefix(info.Name(), ".") {
			continue
		}

		stat := info.Sys().(*syscall.Stat_t)
		usr, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
		group, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))

		fi := FileInfo{
			Name:      info.Name(),
			Size:      info.Size(),
			Mode:      info.Mode(),
			ModTime:   info.ModTime(),
			IsDir:     info.IsDir(),
			User:      usr.Username,
			Group:     group.Name,
			Links:     int(stat.Nlink),
			BlockSize: stat.Blocks,
		}

		fileInfos = append(fileInfos, fi)
		totalBlocks += fi.BlockSize
	}

	sortEntries(fileInfos, opts)

	if opts.LongFormat {
		fmt.Printf("total %d\n", totalBlocks/2) // Convert to 1K-blocks
		for _, fi := range fileInfos {
			displayLongFormat(fi)
		}
	} else {
		displayColumns(fileInfos)
	}

	if opts.Recursive {
		for _, fi := range fileInfos {
			if fi.IsDir && fi.Name != "." && fi.Name != ".." {
				newPath := filepath.Join(path, fi.Name)
				fmt.Printf("\n%s:\n", newPath)
				if err := listDirectory(newPath, opts); err != nil {
					fmt.Println("ls: cannot open directory", err) // output: ls: cannot open directory '/etc/ssl/private': Permission denied
					continue
				}
			}
		}
	}

	return nil
}

func sortEntries(entries []FileInfo, opts Options) {
	sort.Slice(entries, func(i, j int) bool {
		if opts.SortByTime {
			if entries[i].ModTime.Equal(entries[j].ModTime) {
				return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
			}
			return entries[i].ModTime.After(entries[j].ModTime)
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	if opts.Reverse {
		for i := 0; i < len(entries)/2; i++ {
			j := len(entries) - 1 - i
			entries[i], entries[j] = entries[j], entries[i]
		}
	}
}

func displayLongFormat(fi FileInfo) {
	fmt.Printf("%s %2d %-8s %-8s %8d %s %s\n",
		fi.Mode,
		fi.Links,
		fi.User,
		fi.Group,
		fi.Size,
		fi.ModTime.Format("Jan _2 15:04"),
		fi.Name)
}

func displayColumns(files []FileInfo) {
	if len(files) == 0 {
		return
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80 // fallback to 80 if unable to get terminal width
	}

	// Calculate the maximum width of file names
	maxWidth := 0
	for _, file := range files {
		if len(file.Name) > maxWidth {
			maxWidth = len(file.Name)
		}
	}

	columns := width / (maxWidth + 2) // +2 for spacing between columns
	if columns == 0 {
		columns = 1
	}

	rows := (len(files) + columns - 1) / columns

	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			index := col*rows + row
			if index < len(files) {
				fmt.Printf("%-*s", maxWidth+2, files[index].Name)
			}
		}
		fmt.Println()
	}
}
