package pkg

import (
	"os"
	"time"
)

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
