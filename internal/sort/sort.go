package sort

import (
	"sort"
	"strconv"

	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
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

		return CompareFilenamesAlphanumeric(files[i].Name, files[j].Name)
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

func CompareFilenamesAlphanumeric(a, b string) bool {
	aRunes := []rune(a)
	bRunes := []rune(b)
	aLen := len(aRunes)
	bLen := len(bRunes)

	for i, j := 0, 0; i < aLen && j < bLen; {
		// Skip non-alphanumeric characters
		for i < aLen && !IsAlphanumeric(aRunes[i]) {
			i++
		}
		for j < bLen && !IsAlphanumeric(bRunes[j]) {
			j++
		}

		// If we've reached the end of either string, compare lengths
		if i == aLen || j == bLen {
			return aLen < bLen
		}

		// If both characters are digits, compare the whole number
		if IsDigit(aRunes[i]) && IsDigit(bRunes[j]) {
			aNum, aEnd := ExtractNumber(aRunes[i:])
			bNum, bEnd := ExtractNumber(bRunes[j:])

			if aNum != bNum {
				return aNum < bNum
			}

			i += aEnd
			j += bEnd
		} else {
			// Compare characters case-insensitively
			aLower := ToLower(aRunes[i])
			bLower := ToLower(bRunes[j])
			if aLower != bLower {
				return aLower < bLower
			}
			i++
			j++
		}
	}

	// If all compared characters are the same, shorter string comes first
	return aLen < bLen
}

func IsAlphanumeric(r rune) bool {
	return IsLetter(r) || IsDigit(r)
}

func ExtractNumber(runes []rune) (int, int) {
	num := 0
	i := 0
	for i < len(runes) && IsDigit(runes[i]) {
		digit, _ := strconv.Atoi(string(runes[i]))
		num = num*10 + digit
		i++
	}
	return num, i
}

func IsDigit(r rune) bool {
	return (r >= '0' && r <= '9')
}

func ToLower( r rune) rune {
	if (r >= 'A' && r <= 'Z'){
		r  = r + ('a'- 'A')
	}
	return r
}

func IsLetter(r rune) bool {
	return ((r >= 'A' && r <= 'Z')|| (r >= 'a' && r <= 'z') )
}