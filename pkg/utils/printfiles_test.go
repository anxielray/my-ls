package utils

import (
	"testing"
	"time"

	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
)

func TestPrintColumnar(t *testing.T) {
	type args struct {
		files   []FI.FileInfo
		options OP.Options
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Print files in columnar format",
			args: args{
				files: []FI.FileInfo{
					{Name: "file1.txt"},
					{Name: "file2.txt"},
					{Name: "file3.txt"},
				},
				options: OP.Options{OnePerLine: false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintColumnar(tt.args.files, tt.args.options)
		})
	}
}

func TestPrintFiles(t *testing.T) {
	type args struct {
		files   []FI.FileInfo
		options OP.Options
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Print files in long format",
			args: args{
				files: []FI.FileInfo{
					{Name: "file1.txt"},
				},
				options: OP.Options{LongFormat: true},
			},
		},
		{
			name: "Print files one per line",
			args: args{
				files: []FI.FileInfo{
					{Name: "file1.txt"},
					{Name: "file2.txt"},
				},
				options: OP.Options{OnePerLine: true},
			},
		},
		{
			name: "Print files in columnar format",
			args: args{
				files: []FI.FileInfo{
					{Name: "file1.txt"},
					{Name: "file2.txt"},
				},
				options: OP.Options{OnePerLine: false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintFiles(tt.args.files, tt.args.options)
		})
	}
}

func TestMajor(t *testing.T) {
	type args struct {
		dev uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "Major device number",
			args: args{dev: 512},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Major(tt.args.dev); got != tt.want {
				t.Errorf("Major() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinor(t *testing.T) {
	type args struct {
		dev uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "Minor device number",
			args: args{dev: 512},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Minor(tt.args.dev); got != tt.want {
				t.Errorf("Minor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isStandardLibrary(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Standard library path",
			args: args{path: "/usr/bin"},
			want: true,
		},
		{
			name: "Non-standard library path",
			args: args{path: "/home/user"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isStandardLibrary(tt.args.path); got != tt.want {
				t.Errorf("isStandardLibrary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Clean normal path",
			args: args{path: "/usr//local/../bin/./file.txt"},
			want: "/usr/bin/file.txt",
		},
		{
			name: "Clean root path",
			args: args{path: "/.."},
			want: "/",
		},
		{
			name: "Clean current directory",
			args: args{path: "."},
			want: ".",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanPath(tt.args.path); got != tt.want {
				t.Errorf("CleanPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function for generating current time for testing
func now() time.Time {
	return time.Now()
}
