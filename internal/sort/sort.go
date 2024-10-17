package sort

import (
	"sort"

	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
	U "my-ls-1/pkg/utils"
)

func SortFiles(files []FI.FileInfo, options OP.Options) {
	sort.Slice(files, func(i, j int) bool {
		if options.SortByTime {
			if !files[i].ModTime.Equal(files[j].ModTime) {
				return files[i].ModTime.After(files[j].ModTime)
			}
		} else if options.SortBySize {
			if files[i].Size != files[j].Size {
				return files[i].Size > files[j].Size
			}
		}

		return U.CompareFilenamesAlphanumeric(files[i].Name, files[j].Name)
	})

	if options.Reverse {
		ReverseSlice(files)
	}
}

func ReverseSlice(slice []FI.FileInfo) {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - 1 - i
		slice[i], slice[j] = slice[j], slice[i]
	}
}
