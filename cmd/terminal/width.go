package terminal

import (
	"os"
	"strconv"
)

func GetTerminalWidth() int {
	defaultWidth := 80

	// Try getting width using COLUMNS environment variable
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			return width
		}
	}

	// Try getting width using TERM_COLUMNS as an alternative
	if cols := os.Getenv("TERM_COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			return width
		}
	}

	return defaultWidth
}
