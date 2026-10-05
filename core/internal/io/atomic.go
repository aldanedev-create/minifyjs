package io

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteFileAtomic writes data to path atomically by writing to a
// sibling temp file and renaming it into place.
//
// This prevents partial writes from being observed by another
// process (or by a build watcher) if the process crashes mid-write.
// The temp file is created in the same directory as path so the
// rename stays on one filesystem and is therefore atomic on POSIX.
//
// The file mode is applied to the final file, not the temp file.
func WriteFileAtomic(path string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".minifyjs-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tmpName := tmp.Name()

	// Best-effort cleanup if anything below fails.
	cleanup := func() {
		_ = os.Remove(tmpName)
	}

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("write temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return fmt.Errorf("close temp: %w", err)
	}
	if err := os.Chmod(tmpName, mode); err != nil {
		cleanup()
		return fmt.Errorf("chmod temp: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		cleanup()
		return fmt.Errorf("rename temp: %w", err)
	}
	return nil
}