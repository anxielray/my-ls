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
			name: "regular file without color",
			args: args{
				file: FI.FileInfo{
					Name: "test.txt",
				},
				options: OP.Options{
					NoColor: true,
				},
			},
			want: "test.txt",
		},
		{
			name: "directory without color",
			args: args{
				file: FI.FileInfo{
					Name:  "testdir",
					IsDir: true,
				},
				options: OP.Options{
					NoColor: true,
				},
			},
			want: "testdir",
		},
		{
			name: "symlink without color",
			args: args{
				file: FI.FileInfo{
					Name:       "testlink",
					IsLink:     true,
					LinkTarget: "target.txt",
				},
				options: OP.Options{
					NoColor: true,
				},
			},
			want: "testlink -> target.txt",
		},
		{
			name: "empty filename",
			args: args{
				file: FI.FileInfo{
					Name: "",
				},
				options: OP.Options{
					NoColor: true,
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatFileName(tt.args.file, tt.args.options); got != tt.want {
				t.Errorf("FormatFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatPermissions(t *testing.T) {
	tests := []struct {
		name string
		mode os.FileMode
		want string
	}{
		{
			name: "no permissions",
			mode: 0,
			want: "---------",
		},
		{
			name: "all permissions",
			mode: 0o777,
			want: "rwxrwxrwx",
		},
		{
			name: "read only",
			mode: 0o444,
			want: "r--r--r--",
		},
		{
			name: "write only",
			mode: 0o222,
			want: "-w--w--w-",
		},
		{
			name: "execute only",
			mode: 0o111,
			want: "--x--x--x",
		},
		{
			name: "setuid with x",
			mode: os.FileMode(0o4755),
			want: "rwsr-xr-x",
		},
		{
			name: "setuid without x",
			mode: os.FileMode(0o4644),
			want: "rwSr--r--",
		},
		{
			name: "setgid with x",
			mode: os.FileMode(0o2755),
			want: "rwxr-sr-x",
		},
		{
			name: "setgid without x",
			mode: os.FileMode(0o2644),
			want: "rw-r-Sr--",
		},
		{
			name: "sticky with x",
			mode: os.FileMode(0o1755),
			want: "rwxr-xr-t",
		},
		{
			name: "sticky without x",
			mode: os.FileMode(0o1644),
			want: "rw-r--r-T",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatPermissions(tt.mode); got != tt.want {
				t.Errorf("FormatPermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatFileMode(t *testing.T) {
	tests := []struct {
		name string
		mode os.FileMode
		want string
	}{
		{
			name: "regular file",
			mode: 0o644,
			want: "-rw-r--r--",
		},
		{
			name: "directory",
			mode: os.ModeDir | 0o755,
			want: "drwxr-xr-x",
		},
		{
			name: "symlink",
			mode: os.ModeSymlink | 0o777,
			want: "lrwxrwxrwx",
		},
		{
			name: "character device",
			mode: os.ModeDevice | os.ModeCharDevice | 0o644,
			want: "crw-r--r--",
		},
		{
			name: "block device",
			mode: os.ModeDevice | 0o644,
			want: "brw-r--r--",
		},
		{
			name: "named pipe",
			mode: os.ModeNamedPipe | 0o644,
			want: "prw-r--r--",
		},
		{
			name: "socket",
			mode: os.ModeSocket | 0o644,
			want: "srw-r--r--",
		},
		{
			name: "setuid file",
			mode: os.ModeSetuid | 0o755,
			want: "-rwsr-xr-x",
		},
		{
			name: "setgid file",
			mode: os.ModeSetgid | 0o755,
			want: "-rwxr-sr-x",
		},
		{
			name: "sticky directory",
			mode: os.ModeDir | os.ModeSticky | 0o755,
			want: "drwxr-xr-t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatFileMode(tt.mode); got != tt.want {
				t.Errorf("FormatFileMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
