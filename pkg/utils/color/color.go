package color

import (
	"os"
	"path/filepath"
	"strings"

	FI "my-ls-1/pkg/fileinfo"
)

const (
	ColorReset = "\033[0m"
)

var colorMap map[string]string

func InitColorMap() {
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

func Colorize(file FI.FileInfo, name string) string {
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
