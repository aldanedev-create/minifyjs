package main

import miniio "github.com/minifyjs/minifyjs/core/internal/io"

// writeFile is a small wrapper used by the watch loop. The main
// runTransform path calls miniio.WriteFileAtomic directly.
func writeFile(path string, data []byte) error {
	return miniio.WriteFileAtomic(path, data, 0o644)
}
