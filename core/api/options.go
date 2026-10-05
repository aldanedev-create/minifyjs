// Package api is the public Go surface of MinifyJS. It is the only
// package external Go callers should import; every internal package
// is free to change shape between minor versions.
package api

// Options mirrors engine.Options. Kept as a separate type so the
// public API does not leak internal package names into its
// signatures.
type Options struct {
	MinifyWhitespace  bool
	MinifyIdentifiers bool
	MinifySyntax      bool
	Target            string
	Format            string
	Sourcemap         string
	Banner            string
	Footer            string
	LegalComments     string
	SourceName        string
	Define            map[string]string
	Drop              []string
	Pure              []string
}