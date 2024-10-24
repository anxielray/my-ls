package sort

import (
	"testing"
	"time"

	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
)

func TestSortFiles(t *testing.T) {
	baseTime := time.Now()

	type args struct {
		files   []FI.FileInfo
		options OP.Options
	}

	tests := []struct {
		name      string
		args      args
		wantOrder []string
	}{
		{
			name: "Default sort (alphanumeric)",
			args: args{
				files: []FI.FileInfo{
					{Name: "file2.txt"},
					{Name: "file1.txt"},
					{Name: "file10.txt"},
				},
				options: OP.Options{},
			},
			wantOrder: []string{"file1.txt", "file2.txt", "file10.txt"},
		},
		{
			name: "Sort by time",
			args: args{
				files: []FI.FileInfo{
					{Name: "old.txt", ModTime: baseTime.Add(-2 * time.Hour)},
					{Name: "new.txt", ModTime: baseTime},
					{Name: "medium.txt", ModTime: baseTime.Add(-1 * time.Hour)},
				},
				options: OP.Options{SortByTime: true},
			},
			wantOrder: []string{"new.txt", "medium.txt", "old.txt"},
		},
		{
			name: "Sort by size",
			args: args{
				files: []FI.FileInfo{
					{Name: "small.txt", Size: 100},
					{Name: "large.txt", Size: 300},
					{Name: "medium.txt", Size: 200},
				},
				options: OP.Options{SortBySize: true},
			},
			wantOrder: []string{"large.txt", "medium.txt", "small.txt"},
		},
		{
			name: "Reverse alphanumeric sort",
			args: args{
				files: []FI.FileInfo{
					{Name: "b.txt"},
					{Name: "a.txt"},
					{Name: "c.txt"},
				},
				options: OP.Options{Reverse: true},
			},
			wantOrder: []string{"c.txt", "b.txt", "a.txt"},
		},
		{
			name: "Sort mixed case alphanumeric",
			args: args{
				files: []FI.FileInfo{
					{Name: "B.txt"},
					{Name: "a.txt"},
					{Name: "C.txt"},
				},
				options: OP.Options{},
			},
			wantOrder: []string{"a.txt", "B.txt", "C.txt"},
		},
		{
			name: "Sort with numbers in names",
			args: args{
				files: []FI.FileInfo{
					{Name: "file100.txt"},
					{Name: "file20.txt"},
					{Name: "file3.txt"},
				},
				options: OP.Options{},
			},
			wantOrder: []string{"file3.txt", "file20.txt", "file100.txt"},
		},
		{
			name: "Reverse size sort",
			args: args{
				files: []FI.FileInfo{
					{Name: "small.txt", Size: 100},
					{Name: "large.txt", Size: 300},
					{Name: "medium.txt", Size: 20},
				},
				options: OP.Options{SortBySize: true, Reverse: true},
			},
			wantOrder: []string{"small.txt", "medium.txt", "large.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalFiles := make([]FI.FileInfo, len(tt.args.files))
			copy(originalFiles, tt.args.files)

			SortFiles(tt.args.files, tt.args.options)

			// Check if the order matches expected
			if len(tt.args.files) != len(tt.wantOrder) {
				t.Errorf("SortFiles() resulted in different length: got %v, want %v", len(tt.args.files), len(tt.wantOrder))
				return
			}

			for i, want := range tt.wantOrder {
				if tt.args.files[i].Name != want {
					t.Errorf("SortFiles() incorrect order at position %d: got %v, want %v", i, tt.args.files[i].Name, want)
				}
			}
		})
	}
}

func TestReverseSlice(t *testing.T) {
	tests := []struct {
		name  string
		files []FI.FileInfo
		want  []string
	}{
		{
			name: "Reverse normal slice",
			files: []FI.FileInfo{
				{Name: "a.txt"},
				{Name: "b.txt"},
				{Name: "c.txt"},
			},
			want: []string{"c.txt", "b.txt", "a.txt"},
		},
		{
			name: "Reverse single element",
			files: []FI.FileInfo{
				{Name: "a.txt"},
			},
			want: []string{"a.txt"},
		},
		{
			name:  "Reverse empty slice",
			files: []FI.FileInfo{},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReverseSlice(tt.files)
			for i, want := range tt.want {
				if tt.files[i].Name != want {
					t.Errorf("ReverseSlice() incorrect order at position %d: got %v, want %v", i, tt.files[i].Name, want)
				}
			}
		})
	}
}

func TestCompareFilenamesAlphanumeric(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected bool
	}{
		{
			name:     "Simple alphabetic",
			a:        "a.txt",
			b:        "b.txt",
			expected: true,
		},
		{
			name:     "Numeric comparison",
			a:        "file2",
			b:        "file10",
			expected: true,
		},
		{
			name:     "Mixed case",
			a:        "A.txt",
			b:        "b.txt",
			expected: true,
		},
		{
			name:     "Special characters",
			a:        ".hidden",
			b:        "visible",
			expected: true,
		},
		{
			name:     "Same string",
			a:        "file.txt",
			b:        "file.txt",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareFilenamesAlphanumeric(tt.a, tt.b); got != tt.expected {
				t.Errorf("CompareFilenamesAlphanumeric() = %v, want %v", got, tt.expected)
			}
		})
	}
}
