package utils

import (
	"os"
	"testing"

	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
)

func TestFormatFileName(t *testing.T) {
	type args struct {
		file    FI.FileInfo
		options OP.Options
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal file without color",
			args: args{
				file: FI.FileInfo{
					Name:   "file.txt",
					IsLink: false,
				},
				options: OP.Options{NoColor: true},
			},
			want: "file.txt",
		},
		{
			name: "normal file with color",
			args: args{
				file: FI.FileInfo{
					Name:   "file.txt",
					IsLink: false,
				},
				options: OP.Options{NoColor: false},
			},
			want: "file.txt", // Adjust this based on actual color implementation
		},
		{
			name: "symlink file",
			args: args{
				file: FI.FileInfo{
					Name:       "link.txt",
					IsLink:     true,
					LinkTarget: "target.txt",
				},
				options: OP.Options{NoColor: true},
			},
			want: "link.txt -> target.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatFileName(tt.args.file, tt.args.options)
			if got != tt.want {
				t.Errorf("FormatFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatPermissions(t *testing.T) {
	type args struct {
		mode os.FileMode
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "read/write/execute permissions",
			args: args{mode: 0o755},
			want: "rwxr-xr-x",
		},
		{
			name: "no permissions",
			args: args{mode: 0o000},
			want: "---------",
		},
		{
			name: "setuid permission",
			args: args{mode: 0o4755},
			want: "rwxr-xr-x",
		},
		{
			name: "setgid permission",
			args: args{mode: 0o2755},
			want: "rwxr-xr-x",
		},
		{
			name: "sticky bit",
			args: args{mode: 0o1777},
			want: "rwxrwxrwx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatPermissions(tt.args.mode)
			if got != tt.want {
				t.Errorf("FormatPermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatFileMode(t *testing.T) {
	type args struct {
		mode os.FileMode
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "regular file",
			args: args{mode: 0o644},
			want: "-rw-r--r--",
		},
		{
			name: "directory",
			args: args{mode: os.ModeDir | 0o755},
			want: "drwxr-xr-x",
		},
		{
			name: "symlink",
			args: args{mode: os.ModeSymlink},
			want: "l---------",
		},
		{
			name: "character device",
			args: args{mode: os.ModeDevice | os.ModeCharDevice},
			want: "c---------",
		},
		{
			name: "block device",
			args: args{mode: os.ModeDevice},
			want: "b---------",
		},
		{
			name: "named pipe",
			args: args{mode: os.ModeNamedPipe},
			want: "p---------",
		},
		{
			name: "socket",
			args: args{mode: os.ModeSocket},
			want: "s---------",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatFileMode(tt.args.mode)
			if got != tt.want {
				t.Errorf("FormatFileMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
