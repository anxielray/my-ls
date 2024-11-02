package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDirectory(t *testing.T) (string, func()) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "block-test")
	if err != nil {
		t.Fatal(err)
	}

	// Create test file structure
	files := map[string]int64{
		"empty.txt":         0,            // 0 blocks
		"small.txt":         100,          // 1 block
		"exact.txt":         512,          // 1 block
		"large.txt":         1024,         // 2 blocks
		"subdir/nested.txt": 2048,         // 4 blocks
		"uneven.txt":        600,          // 2 blocks (rounds up)
		"verylarge.txt":     512 * 100,    // 100 blocks
		"subdir/big.txt":    512*10 + 100, // 11 blocks
	}

	for path, size := range files {
		fullPath := filepath.Join(tmpDir, path)

		// Create directory if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}

		// Create and write to file
		f, err := os.Create(fullPath)
		if err != nil {
			t.Fatal(err)
		}

		// Create file with specific size
		if err := f.Truncate(size); err != nil {
			f.Close()
			t.Fatal(err)
		}
		f.Close()
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func Test_calculateTotalBlocks(t *testing.T) {
	tmpDir, cleanup := setupTestDirectory(t)
	defer cleanup()

	// Create a non-readable directory for error testing
	nonReadableDir := filepath.Join(tmpDir, "noaccess")
	if err := os.MkdirAll(nonReadableDir, 0o000); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		dir     string
		want    int64
		wantErr bool
	}{
		{
			name:    "single empty file",
			dir:     filepath.Join(tmpDir, "empty.txt"),
			want:    0,
			wantErr: false,
		},
		{
			name:    "single small file",
			dir:     filepath.Join(tmpDir, "small.txt"),
			want:    1, // 100 bytes = 1 block
			wantErr: false,
		},
		{
			name:    "exact block size file",
			dir:     filepath.Join(tmpDir, "exact.txt"),
			want:    1, // 512 bytes = 1 block
			wantErr: false,
		},
		{
			name:    "large file",
			dir:     filepath.Join(tmpDir, "large.txt"),
			want:    2, // 1024 bytes = 2 blocks
			wantErr: false,
		},
		{
			name:    "uneven size file",
			dir:     filepath.Join(tmpDir, "uneven.txt"),
			want:    2, // 600 bytes = 2 blocks (rounds up)
			wantErr: false,
		},
		{
			name:    "directory with nested files",
			dir:     filepath.Join(tmpDir, "subdir"),
			want:    15, // nested.txt(4) + big.txt(11) = 15 blocks
			wantErr: false,
		},
		{
			name:    "non-existent directory",
			dir:     filepath.Join(tmpDir, "nonexistent"),
			want:    0,
			wantErr: true,
		},
		{
			name:    "non-readable directory",
			dir:     nonReadableDir,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateTotalBlocks(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateTotalBlocks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateTotalBlocks() = %v, want %v", got, tt.want)
			}
		})
	}
}
