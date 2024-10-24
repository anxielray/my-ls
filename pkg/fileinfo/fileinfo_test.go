package fileinfo

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestFiles(t *testing.T) (string, func()) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "fileinfo_test")
	if err != nil {
		t.Fatal(err)
	}

	// Create a regular file
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test content"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a directory
	testDir := filepath.Join(tmpDir, "testdir")
	if err := os.Mkdir(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a symlink
	symlinkFile := filepath.Join(tmpDir, "symlink")
	if err := os.Symlink(regularFile, symlinkFile); err != nil {
		t.Fatal(err)
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestCreateFileInfo(t *testing.T) {
	// Setup test files and cleanup
	tmpDir, cleanup := setupTestFiles(t)
	defer cleanup()

	// Get FileInfo for test files
	regularFileInfo, _ := os.Stat(filepath.Join(tmpDir, "regular.txt"))
	dirInfo, _ := os.Stat(filepath.Join(tmpDir, "testdir"))
	symlinkInfo, _ := os.Lstat(filepath.Join(tmpDir, "symlink"))

	type args struct {
		path string
		info os.FileInfo
	}
	tests := []struct {
		name string
		args args
		want FileInfo
	}{
		{
			name: "Regular file",
			args: args{
				path: tmpDir,
				info: regularFileInfo,
			},
			want: FileInfo{
				Name:    "regular.txt",
				Size:    int64(len("test content")),
				Mode:    regularFileInfo.Mode(),
				ModTime: regularFileInfo.ModTime(),
				IsDir:   false,
				IsLink:  false,
			},
		},
		{
			name: "Directory",
			args: args{
				path: tmpDir,
				info: dirInfo,
			},
			want: FileInfo{
				Name:    "testdir",
				Size:    dirInfo.Size(),
				Mode:    dirInfo.Mode(),
				ModTime: dirInfo.ModTime(),
				IsDir:   true,
				IsLink:  false,
			},
		},
		{
			name: "Symlink",
			args: args{
				path: tmpDir,
				info: symlinkInfo,
			},
			want: FileInfo{
				Name:       "symlink",
				Size:       symlinkInfo.Size(),
				Mode:       symlinkInfo.Mode(),
				ModTime:    symlinkInfo.ModTime(),
				IsDir:      false,
				IsLink:     true,
				LinkTarget: filepath.Join(tmpDir, "regular.txt"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateFileInfo(tt.args.path, tt.args.info)

			// Compare fields that we can predict/control
			if got.Name != tt.want.Name {
				t.Errorf("CreateFileInfo().Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.IsDir != tt.want.IsDir {
				t.Errorf("CreateFileInfo().IsDir = %v, want %v", got.IsDir, tt.want.IsDir)
			}
			if got.IsLink != tt.want.IsLink {
				t.Errorf("CreateFileInfo().IsLink = %v, want %v", got.IsLink, tt.want.IsLink)
			}
			if got.IsLink && got.LinkTarget != tt.want.LinkTarget {
				t.Errorf("CreateFileInfo().LinkTarget = %v, want %v", got.LinkTarget, tt.want.LinkTarget)
			}
			if got.Mode.Perm() != tt.want.Mode.Perm() {
				t.Errorf("CreateFileInfo().Mode = %v, want %v", got.Mode, tt.want.Mode)
			}

			// Check if ModTime is within a reasonable range (1 second)
			if got.ModTime.Sub(tt.want.ModTime).Abs() > time.Second {
				t.Errorf("CreateFileInfo().ModTime = %v, want close to %v", got.ModTime, tt.want.ModTime)
			}

			// For regular files, check size
			if !got.IsDir && !got.IsLink && got.Size != tt.want.Size {
				t.Errorf("CreateFileInfo().Size = %v, want %v", got.Size, tt.want.Size)
			}

			// System-specific fields (Uid, Gid, etc.) are present but values depend on the system
			if got.Uid == 0 || got.Gid == 0 {
				t.Error("CreateFileInfo() system-specific fields not set")
			}
		})
	}
}
