package utils

import (
	"os"
	"reflect"
	"testing"

	FI "my-ls-1/pkg/fileinfo"
)

func Test_checkPathType(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Existing directory",
			args: args{path: "testdir"},
			want: "directory",
		},
		{
			name: "Existing file",
			args: args{path: "testfile.txt"},
			want: "file",
		},
		{
			name:    "Non-existing path",
			args:    args{path: "nonexistent"},
			wantErr: true,
		},
		{
			name: "Existing symlink",
			args: args{path: "testlink"},
			want: "symlink",
		},
	}

	// Setup: create test files, directories, and symlinks
	os.Mkdir("testdir", 0o755)
	defer os.RemoveAll("testdir") // Clean up all contents

	_, err := os.Create("testfile.txt")
	if err != nil {
		t.Fatalf("failed to create testfile.txt: %v", err)
	}
	defer os.Remove("testfile.txt") // Clean up

	// Create a symlink
	if err := os.Symlink("testfile.txt", "testlink"); err != nil {
		t.Fatalf("failed to create testlink: %v", err)
	}
	defer os.Remove("testlink") // Clean up

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkPathType(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkPathType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("checkPathType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSymlink(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Existing symlink",
			args: args{path: "testlink"},
			want: true,
		},
		{
			name: "Existing file (not symlink)",
			args: args{path: "testfile.txt"},
			want: false,
		},
		{
			name:    "Non-existing path",
			args:    args{path: "nonexistent"},
			wantErr: true,
		},
	}

	// Setup: create a test file and a symlink
	_, err := os.Create("testfile.txt")
	if err != nil {
		t.Fatalf("failed to create testfile.txt: %v", err)
	}
	defer os.Remove("testfile.txt") // Clean up

	if err := os.Symlink("testfile.txt", "testlink"); err != nil {
		t.Fatalf("failed to create testlink: %v", err)
	}
	defer os.Remove("testlink") // Clean up

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsSymlink(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsSymlink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsSymlink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSymlinksInDir(t *testing.T) {
	type args struct {
		dirPath string
	}
	tests := []struct {
		name    string
		args    args
		want    []FI.FileInfo
		wantErr bool
	}{
		{
			name:    "Non-existing directory",
			args:    args{dirPath: "nonexistent"},
			wantErr: true,
		},
	}

	// Setup: create a test directory and a symlink
	os.Mkdir("testdir", 0o755)
	defer os.RemoveAll("testdir") // Clean up all contents

	_, err := os.Create("testfile.txt")
	if err != nil {
		t.Fatalf("failed to create testfile.txt: %v", err)
	}
	defer os.Remove("testfile.txt") // Clean up

	if err := os.Symlink("testfile.txt", "testdir/testlink"); err != nil {
		t.Fatalf("failed to create symlink in testdir: %v", err)
	}
	defer os.Remove("testdir/testlink") // Clean up

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSymlinksInDir(tt.args.dirPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSymlinksInDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSymlinksInDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
