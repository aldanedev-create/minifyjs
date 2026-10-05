package source

// SourceFile is the on-disk identity of a Source, kept separate from
// the Source itself so callers can pass around a lightweight handle
// without carrying the text.
//
// It exists mainly so the cache and any future bundler have a
// common, hashable key type for files they read from disk.
type SourceFile struct {
	// Path is the absolute or repo-relative path, as given by the
	// caller. Empty for stdin or in-memory sources.
	Path string
	// Size is the byte length. Kept so the cache can quickly reject
	// entries whose size has changed even before hashing.
	Size int
}

// IsStdin reports whether this SourceFile refers to stdin (no path).
func (f SourceFile) IsStdin() bool { return f.Path == "" }