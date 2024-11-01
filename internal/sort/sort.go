package sort

import (
	"strconv"

	FI "my-ls-1/pkg/fileinfo"
	OP "my-ls-1/pkg/options"
)

/*This function will take an array of fileInfo and sort them based on the conditions
set by the flags passed on the command line. If the option for reverse was set to true
this action takes place right after the sorting*/
func SortFiles(files []FI.FileInfo, options OP.Options) {
	CustomSort(files, func(i, j int) bool {

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

// CustomSort sorts a slice using a custom less function
func CustomSort(slice interface{}, less func(i, j int) bool) {
	switch s := slice.(type) {
	case []FI.FileInfo:
		bubbleSortFileInfo(s, less)
	default:
		panic("unsupported slice type")
	}
}

// bubbleSortFileInfo implements bubble sort for a slice of FileInfo
func bubbleSortFileInfo(slice []FI.FileInfo, less func(i, j int) bool) {
	n := len(slice)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if less(j+1, j) {
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
}

func ReverseSlice(slice []FI.FileInfo) {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - 1 - i
		slice[i], slice[j] = slice[j], slice[i]
	}
}

//sorts the files rune by rune, pioritizing the special characters and numerical strings
func CompareFilenamesAlphanumeric(a, b string) bool {
	aRunes := []rune(a)
	bRunes := []rune(b)
	aLen := len(aRunes)
	bLen := len(bRunes)

	for i, j := 0, 0; i < aLen && j < bLen; {

		if i == aLen || j == bLen {
			return aLen < bLen
		}

		if IsSpecialCharacter(aRunes[i]) && !IsSpecialCharacter(bRunes[j]) {
			return true
		}
		if !IsSpecialCharacter(aRunes[i]) && IsSpecialCharacter(bRunes[j]) {
			return false
		}

		if IsDigit(aRunes[i]) && IsDigit(bRunes[j]) {
			aNum, aEnd := ExtractNumber(aRunes[i:])
			bNum, bEnd := ExtractNumber(bRunes[j:])

			if aNum != bNum {
				return aNum < bNum
			}

			i += aEnd
			j += bEnd
		} else {
			aLower := ToLower(aRunes[i])
			bLower := ToLower(bRunes[j])
			if aLower != bLower {
				return aLower < bLower
			}
		}
		i++
		j++
	}

	return aLen < bLen
}

// IsAlphanumeric checks if a rune is a letter (A-Z, a-z) or a digit (0-9).
func IsAlphanumeric(r rune) bool {
	return IsLetter(r) || IsDigit(r)
}

//This function gets the int from the rune provided
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

//checks if the rune is is a digit(0 -9)
func IsDigit(r rune) bool {
	return (r >= '0' && r <= '9')
}

//converts a rune to a lowercase string if is in uppercase
func ToLower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		r = r + ('a' - 'A')
	}
	return r
}

//check if the rune is a letter
func IsLetter(r rune) bool {
	return ((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'))
}

//checks if the rune is a special character, despite the ascii value
func IsSpecialCharacter(r rune) bool {
	allowedSpecialChars := "!#$%&'()*+,-./:;<=>?@[]^_`{|}~"

	if r == ' ' || containsRune(allowedSpecialChars, r) {
		return true
	}

	return false
}

// containsRune checks if a string contains a specific rune.
func containsRune(s string, r rune) bool {
	for _, char := range s {
		if char == r {
			return true
		}
	}
	return false
}
