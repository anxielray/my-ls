package internal

import (
	"os"
	"path/filepath"
	"testing"

	OP "my-ls-1/pkg/options"
)

func setupTestFiles(t *testing.T) (string, func()) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "ls-test")
	if err != nil {
		t.Fatal(err)
	}

	// Create test files and directories
	files := []string{
		"file1.txt",
		"file2.txt",
		"hidden.txt",
		"dir1",
		"dir1/subfile1.txt",
		"dir2",
	}

	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		if filepath.Base(f) == "dir1" || filepath.Base(f) == "dir2" {
			if err := os.MkdirAll(path, 0o755); err != nil {
				t.Fatal(err)
			}
		} else {
			if err := os.WriteFile(path, []byte("test content"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}

	// Create a symlink
	if err := os.Symlink(filepath.Join(tmpDir, "file1.txt"), filepath.Join(tmpDir, "symlink.txt")); err != nil {
		t.Fatal(err)
	}

	// Create cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestListSingleFile(t *testing.T) {
	tmpDir, cleanup := setupTestFiles(t)
	defer cleanup()

	type args struct {
		path    string
		options OP.Options
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "regular file",
			args: args{
				path: filepath.Join(tmpDir, "file1.txt"),
				options: OP.Options{
					LongFormat: false,
				},
			},
		},
		{
			name: "regular file with long format",
			args: args{
				path: filepath.Join(tmpDir, "file1.txt"),
				options: OP.Options{
					LongFormat: true,
				},
			},
		},
		{
			name: "symlink",
			args: args{
				path: filepath.Join(tmpDir, "symlink.txt"),
				options: OP.Options{
					LongFormat: false,
				},
			},
		},
		{
			name: "nonexistent file",
			args: args{
				path: filepath.Join(tmpDir, "nonexistent.txt"),
				options: OP.Options{
					LongFormat: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ListSingleFile(tt.args.path, tt.args.options)
		})
	}
}

func TestListDir(t *testing.T) {
	tmpDir, cleanup := setupTestFiles(t)
	defer cleanup()

	type args struct {
		path    string
		options OP.Options
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "list directory",
			args: args{
				path: tmpDir,
				options: OP.Options{
					LongFormat: false,
				},
			},
		},
		{
			name: "list directory with long format",
			args: args{
				path: tmpDir,
				options: OP.Options{
					LongFormat: true,
				},
			},
		},
		{
			name: "list directory with all files",
			args: args{
				path: tmpDir,
				options: OP.Options{
					ShowHidden: true,
				},
			},
		},
		{
			name: "list directory reverse order",
			args: args{
				path: tmpDir,
				options: OP.Options{
					Reverse: true,
				},
			},
		},
		{
			name: "list directory by time",
			args: args{
				path: tmpDir,
				options: OP.Options{
					SortByTime: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ListDir(tt.args.path, tt.args.options)
		})
	}
}

func TestListRecursive(t *testing.T) {
	tmpDir, cleanup := setupTestFiles(t)
	defer cleanup()

	type args struct {
		path    string
		options OP.Options
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "recursive listing",
			args: args{
				path: tmpDir,
				options: OP.Options{
					Recursive: true,
				},
			},
		},
		{
			name: "recursive listing with hidden files",
			args: args{
				path: tmpDir,
				options: OP.Options{
					Recursive:  true,
					ShowHidden: true,
				},
			},
		},
		{
			name: "recursive listing with long format",
			args: args{
				path: tmpDir,
				options: OP.Options{
					Recursive:  true,
					LongFormat: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ListRecursive(tt.args.path, tt.args.options)
		})
	}
}
