package utils

import (
	"os"
	"strings"

	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
	C "my-ls-1/pkg/utils/color"
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
