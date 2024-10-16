package utils

import (
	"os"
	"time"
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
