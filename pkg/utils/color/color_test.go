package color

import (
	"os"
	"testing"

	FI "my-ls-1/pkg/fileinfo"
)

func TestInitColorMap(t *testing.T) {
	// Save original LS_COLORS
	originalColors := os.Getenv("LS_COLORS")
	defer os.Setenv("LS_COLORS", originalColors)

	tests := []struct {
		name      string
		lsColors  string
		wantColor string
		colorKey  string
	}{
		{
			name:      "directory color",
			lsColors:  "di=01;34:ln=01;36",
			wantColor: "\033[01;34m",
			colorKey:  "di",
		},
		{
			name:      "symlink color",
			lsColors:  "di=01;34:ln=01;36",
			wantColor: "\033[01;36m",
			colorKey:  "ln",
		},
		{
			name:      "empty LS_COLORS",
			lsColors:  "",
			wantColor: "",
			colorKey:  "di",
		},
		{
			name:      "malformed color entry",
			lsColors:  "di=01;34:invalid:ln=01;36",
			wantColor: "\033[01;34m",
			colorKey:  "di",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LS_COLORS", tt.lsColors)
			InitColorMap()
			if got := colorMap[tt.colorKey]; got != tt.wantColor {
				t.Errorf("InitColorMap() colorMap[%v] = %v, want %v", tt.colorKey, got, tt.wantColor)
			}
		})
	}
}

func TestColorize(t *testing.T) {
	// Setup color map for testing
	os.Setenv("LS_COLORS", "di=01;34:ln=01;36:ex=01;32:pi=40;33:so=01;35:bd=40;33;01:*.txt=01;31")
	InitColorMap()

	type args struct {
		file FI.FileInfo
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "directory",
			args: args{
				file: FI.FileInfo{
					Name:  "testdir",
					IsDir: true,
				},
				name: "testdir",
			},
			want: "\033[01;34mtestdir" + ColorReset,
		},
		{
			name: "symlink",
			args: args{
				file: FI.FileInfo{
					Name:   "testlink",
					IsLink: true,
				},
				name: "testlink",
			},
			want: "\033[01;36mtestlink" + ColorReset,
		},
		{
			name: "executable file",
			args: args{
				file: FI.FileInfo{
					Name: "testexec",
					Mode: 0o755,
				},
				name: "testexec",
			},
			want: "\033[01;32mtestexec" + ColorReset,
		},
		{
			name: "named pipe",
			args: args{
				file: FI.FileInfo{
					Name: "testpipe",
					Mode: os.ModeNamedPipe,
				},
				name: "testpipe",
			},
			want: "\033[40;33mtestpipe" + ColorReset,
		},
		{
			name: "socket",
			args: args{
				file: FI.FileInfo{
					Name: "testsocket",
					Mode: os.ModeSocket,
				},
				name: "testsocket",
			},
			want: "\033[01;35mtestsocket" + ColorReset,
		},
		{
			name: "device",
			args: args{
				file: FI.FileInfo{
					Name: "testdevice",
					Mode: os.ModeDevice,
				},
				name: "testdevice",
			},
			want: "\033[40;33;01mtestdevice" + ColorReset,
		},
		{
			name: "text file",
			args: args{
				file: FI.FileInfo{
					Name: "test.txt",
					Mode: 0o644,
				},
				name: "test.txt",
			},
			want: "\033[01;31mtest.txt" + ColorReset,
		},
		{
			name: "regular file no color",
			args: args{
				file: FI.FileInfo{
					Name: "test.unknown",
					Mode: 0o644,
				},
				name: "test.unknown",
			},
			want: "test.unknown",
		},
		{
			name: "empty name",
			args: args{
				file: FI.FileInfo{},
				name: "",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Colorize(tt.args.file, tt.args.name); got != tt.want {
				t.Errorf("Colorize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorizeWithoutLSColors(t *testing.T) {
	// Save original LS_COLORS
	originalColors := os.Getenv("LS_COLORS")
	defer os.Setenv("LS_COLORS", originalColors)

	// Clear LS_COLORS
	os.Setenv("LS_COLORS", "")
	InitColorMap()

	file := FI.FileInfo{
		Name:  "test.txt",
		IsDir: false,
		Mode:  0o644,
	}

	got := Colorize(file, "test.txt")
	want := "test.txt"

	if got != want {
		t.Errorf("Colorize() with no LS_COLORS = %v, want %v", got, want)
	}
}
