package fileutil

import "os"

func DeleteFile(filePath string) error {
	// once one, remove the file
	return os.Remove(filePath)
}
