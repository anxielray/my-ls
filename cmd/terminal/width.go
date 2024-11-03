package terminal

import (
	"os"
	"strconv"
)

// This function will retrieve a logical terminal width columns with a default of 80
func GetTerminalWidth() int {
	defaultWidth := 80

	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			return width
		}
	}

	if cols := os.Getenv("TERM_COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			return width
		}
	}

	return defaultWidth
}
