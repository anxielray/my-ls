package utils

import (
	"testing"

	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
)

func Test_calculateTotalBlocks(t *testing.T) {
	type args struct {
		dir     string
		options OP.Options
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "Empty directory",
			args: args{
				dir:     "test_empty_dir",
				options: OP.Options{ShowHidden: false},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Directory with files",
			args: args{
				dir:     "test_dir_with_files",
				options: OP.Options{ShowHidden: false},
			},
			want:    0, // Adjust based on actual test file sizes
			wantErr: true,
		},
		{
			name: "Error reading directory",
			args: args{
				dir:     "non_existent_dir",
				options: OP.Options{ShowHidden: false},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Directory with hidden files",
			args: args{
				dir:     "test_dir_with_hidden_files",
				options: OP.Options{ShowHidden: true},
			},
			want:    0, // Adjust based on actual test file sizes including hidden
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateTotalBlocks(tt.args.dir, tt.args.options)
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
			name: "Add existing entry",
			args: args{
				path:  "test_dir/file.txt",
				name:  "file.txt",
				files: &[]FI.FileInfo{},
			},
		},
		{
			name: "Add non-existing entry",
			args: args{
				path:  "non_existent_file.txt",
				name:  "non_existent_file.txt",
				files: &[]FI.FileInfo{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddSpecialEntry(tt.args.path, tt.args.name, tt.args.files)
			// You can assert the length of the files slice here or check specific content
		})
	}
}

// func TestIsHidden(t *testing.T) {
// 	type args struct {
// 		entry os.DirEntry
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bool
// 	}{
// 		{
// 			name: "Hidden directory",
// 			args: args{
// 				entry: os.DirEntryMock{IsDirFunc: func() bool { return true }, NameFunc: func() string { return ".hidden" }},
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "Visible directory",
// 			args: args{
// 				entry: os.DirEntryMock{IsDirFunc: func() bool { return true }, NameFunc: func() string { return "visible" }},
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "Regular file",
// 			args: args{
// 				entry: os.DirEntryMock{IsDirFunc: func() bool { return false }, NameFunc: func() string { return "file.txt" }},
// 			},
// 			want: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := IsHidden(tt.args.entry); got != tt.want {
// 				t.Errorf("IsHidden() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestGetDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal path",
			args: args{path: "/usr/local/bin/file.txt"},
			want: "/usr/local/bin",
		},
		{
			name: "Root path",
			args: args{path: "./file.txt"},
			want: ".",
		},
		{
			name: "Current directory",
			args: args{path: "file.txt"},
			want: ".",
		},
		{
			name: "Trailing slash",
			args: args{path: "/usr/local/bin/"},
			want: "/usr/local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDir(tt.args.path); got != tt.want {
				t.Errorf("GetDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
