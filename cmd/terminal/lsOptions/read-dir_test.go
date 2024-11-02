package lsOptions

import (
	"os"
	"reflect"
	"testing"

	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
)

func TestReadDirectory(t *testing.T) {
	type args struct {
		path    string
		options OP.Options
	}
	tests := []struct {
		name    string
		args    args
		want    []FI.FileInfo
		wantErr bool
	}{
		{
			name: "Non-existing directory",
			args: args{
				path:    "./nonexistent",
				options: OP.Options{ShowHidden: false},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDirectory(tt.args.path, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddSpecialEntry(t *testing.T) {
	type args struct {
		path  string
		name  string
		files *[]FI.FileInfo
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Add current directory entry",
			args: args{
				path:  "./testdir",
				name:  ".",
				files: &[]FI.FileInfo{},
			},
		},
		{
			name: "Add parent directory entry",
			args: args{
				path:  "./testdir",
				name:  "..",
				files: &[]FI.FileInfo{},
			},
		},
	}

	// Setup: Create a temporary directory for testing
	testDir := "./testdir"
	os.MkdirAll(testDir, 0o755)
	defer os.RemoveAll(testDir) // Clean up after tests

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalLen := len(*tt.args.files)
			AddSpecialEntry(tt.args.path, tt.args.name, tt.args.files)
			if len(*tt.args.files) != originalLen+1 {
				t.Errorf("Expected file slice length to increase by 1, got %d", len(*tt.args.files))
			}

			// Verify that the last added entry is correct
			if (*tt.args.files)[originalLen].Name != tt.args.name {
				t.Errorf("Expected last entry name to be %s, got %s", tt.args.name, (*tt.args.files)[originalLen].Name)
			}
		})
	}
}
