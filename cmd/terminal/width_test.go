package terminal

import (
	"os"
	"testing"
)

func TestGetTerminalWidth(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want int
	}{
		{
			name: "Default width when no env vars set",
			env:  nil,
			want: 80,
		},
		{
			name: "Width from COLUMNS",
			env:  map[string]string{"COLUMNS": "100"},
			want: 100,
		},
		{
			name: "Width from TERM_COLUMNS",
			env:  map[string]string{"TERM_COLUMNS": "120"},
			want: 120,
		},
		{
			name: "Invalid COLUMNS value",
			env:  map[string]string{"COLUMNS": "not_a_number"},
			want: 80,
		},
		{
			name: "Invalid TERM_COLUMNS value",
			env:  map[string]string{"TERM_COLUMNS": "not_a_number"},
			want: 80,
		},
		{
			name: "Zero width in COLUMNS",
			env:  map[string]string{"COLUMNS": "0"},
			want: 80,
		},
		{
			name: "Zero width in TERM_COLUMNS",
			env:  map[string]string{"TERM_COLUMNS": "0"},
			want: 80,
		},
		{
			name: "Negative width in COLUMNS",
			env:  map[string]string{"COLUMNS": "-50"},
			want: 80,
		},
		{
			name: "Negative width in TERM_COLUMNS",
			env:  map[string]string{"TERM_COLUMNS": "-50"},
			want: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables for the test
			for key, value := range tt.env {
				os.Setenv(key, value)
			}
			// Clean up environment variables after the test
			defer func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
			}()

			if got := GetTerminalWidth(); got != tt.want {
				t.Errorf("GetTerminalWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}
