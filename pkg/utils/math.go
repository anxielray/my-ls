package utils

import (
	"math"
	"os"
	"path/filepath"
)

//this function calculates the size of the files in the passed path directory and returns the size in MBs
const blockSize = 512 // Default block size in bytes

func calculateTotalBlocks(dir string) (int64, error) {
	var totalBlocks int64

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileSize := info.Size()
			numBlocks := int64(math.Ceil(float64(fileSize) / float64(blockSize)))
			totalBlocks += numBlocks
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return totalBlocks, nil
}
