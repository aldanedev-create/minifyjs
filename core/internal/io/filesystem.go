package io

import "os"

// Exists reports whether path exists on disk (file or directory).
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsFile reports whether path exists and is a regular file.
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// IsDir reports whether path exists and is a directory.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}