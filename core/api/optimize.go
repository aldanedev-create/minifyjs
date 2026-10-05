package api

// Optimize is a convenience wrapper that enables all three
// minification passes. It exists as a named function because the
// Python API exposes optimize() separately from minify().
func Optimize(source string, opts Options) (Result, error) {
	opts.MinifyWhitespace = true
	opts.MinifyIdentifiers = true
	opts.MinifySyntax = true
	return Minify(source, opts)
}