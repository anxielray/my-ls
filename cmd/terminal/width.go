package terminal

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

func GetTerminalWidth() int {
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
