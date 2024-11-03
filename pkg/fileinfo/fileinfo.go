package fileinfo

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

type FileInfo struct {
	Name       string
	Size       int64
	Mode       os.FileMode
	ModTime    time.Time
	IsDir      bool
	Nlink      uint64
	Total      int64
	Uid        uint32
	Gid        uint32
	IsLink     bool
	LinkTarget string
	Rdev       uint64
	Blocks     int64
}

// This function creates a customized FileInfo structure from the standard Golang fileInfo object
func CreateFileInfo(path string, info os.FileInfo) FileInfo {
	fileInfo := FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
		IsLink:  info.Mode()&os.ModeSymlink != 0,
	}

	if fileInfo.IsLink {
		linkTarget, err := os.Readlink(fmt.Sprintf("%s/%s", path, info.Name()))
		if err == nil {
			fileInfo.LinkTarget = linkTarget
		}
	}

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {

		fileInfo.Nlink = stat.Nlink
		fileInfo.Uid = stat.Uid
		fileInfo.Gid = stat.Gid
		fileInfo.Rdev = stat.Rdev
		fileInfo.Blocks = stat.Blocks / 2
	}

	return fileInfo
}
