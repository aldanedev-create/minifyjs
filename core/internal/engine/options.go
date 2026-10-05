// Package engine wraps esbuild's Go API. It is the only package in
// MinifyJS that imports github.com/evanw/esbuild directly; every
// other package (cmd/minifyjs, api, internal/config) goes through
// the types and functions in this package.
//
// Why a wrapper at all? Three reasons:
//
//   1. Stability. esbuild's options struct is large and evolves.
//      MinifyJS exposes a small, stable Options type; if esbuild
//      changes, only this package needs to be updated.
//
//   2. Determinism. esbuild has many defaults that are reasonable
//      for bundling but not necessarily for one-shot minification.
//      The wrapper pins the behavior we want and hides the rest.
//
//   3. Diagnostics. esbuild reports diagnostics as its own message
//      structs. The wrapper converts them into
//      internal/diagnostics.Diagnostic so the CLI and the Python
//      bridge can present them uniformly.
package engine

import "fmt"

// Options is the subset of esbuild's transform/build configuration
// that MinifyJS exposes. It is intentionally small: every field here
// has a stable meaning across versions, and every field maps to at
// most one esbuild option.
//
// The zero value of Options means "no transformation at all" — the
// output equals the input (modulo esbuild's re-printing, which is
// itself a normalizing step). Callers who want minification must set
// at least one of MinifyWhitespace / MinifyIdentifiers / MinifySyntax.
type Options struct {
	// MinifyWhitespace removes whitespace that is not required to
	// preserve token boundaries.
	MinifyWhitespace bool

	// MinifyIdentifiers renames local variables and function
	// parameters to shorter names. It never renames names that are
	// observable from outside the file (globals, exports, property
	// names), so it is safe for typical application code.
	MinifyIdentifiers bool

	// MinifySyntax rewrites code into smaller equivalent forms:
	// constant folding, dead-code elimination, boolean
	// simplification, etc. This is the aggressive pass and is what
	// makes esbuild's output competitive with Terser's.
	MinifySyntax bool

	// Target is the ECMAScript version the output must be
	// compatible with. Empty means "esnext" (no lowering). Valid
	// values are es5, es2015, es2016, ..., es2024, esnext, and the
	// compound form "es2015,chrome58,firefox57,safari11".
	Target string

	// Format controls the module wrapper of the output. Empty means
	// "preserve" (esbuild picks based on input). Valid values:
	// "iife", "cjs", "esm".
	Format string

	// Sourcemap requests a source map. Valid values: "" (none),
	// "inline", "external", "both".
	Sourcemap string

	// Banner is prepended verbatim to the output. Useful for
	// license comments. Empty means no banner.
	Banner string

	// Footer is appended verbatim to the output. Empty means no
	// footer.
	Footer string

	// LegalComments controls what happens to /*! ... */ and
	// //! ... comments. Valid values: "" (default: "eof"),
	// "none", "inline", "eof", "external".
	LegalComments string
	Define        map[string]string
	Drop          []string
	Pure          []string

	// SourceName is the file name associated with the source, used
	// in diagnostics ("app.js:3:5: ..."). Empty means "<stdin>".
	SourceName string
}

// IsNoOp reports whether the options would produce output identical
// to the input (aside from esbuild's normalizing re-print). Useful
// for short-circuiting in the CLI when the user passes no flags.
func (o Options) IsNoOp() bool {
	return !o.MinifyWhitespace &&
		!o.MinifyIdentifiers &&
		!o.MinifySyntax &&
		o.Target == "" &&
		o.Format == "" &&
		o.Sourcemap == "" &&
		o.Banner == "" &&
		o.Footer == "" &&
		o.LegalComments == "" &&
		o.SourceName == "" &&
		len(o.Define) == 0 &&
		len(o.Drop) == 0 &&
		len(o.Pure) == 0
}

// validate returns an error if any field has an unrecognized value.
// It is called by Transform and Build before handing options to
// esbuild so the error message is MinifyJS's, not esbuild's.
func (o Options) validate() error {
	if o.Target != "" && !validTarget(o.Target) {
		return fmt.Errorf("engine: invalid target %q", o.Target)
	}
	switch o.Format {
	case "", "iife", "cjs", "esm":
	default:
		return fmt.Errorf("engine: invalid format %q", o.Format)
	}
	switch o.Sourcemap {
	case "", "inline", "external", "both":
	default:
		return fmt.Errorf("engine: invalid sourcemap %q", o.Sourcemap)
	}
	switch o.LegalComments {
	case "", "none", "inline", "eof", "external":
	default:
		return fmt.Errorf("engine: invalid legal comments mode %q", o.LegalComments)
	}
	return nil
}

// validTarget accepts the same forms esbuild accepts: an ECMAScript
// year keyword (es5, es2015..es2024, esnext), a browser+version pair
// (chrome58), or a comma-separated list of the above.
func validTarget(t string) bool {
	// Delegated to esbuild in transform.go via a package-level slice
	// of known ECMAScript keywords, plus a permissive check for the
	// browser+version form. Kept loose here because esbuild is the
	// real validator.
	if t == "" {
		return true
	}
	for _, r := range t {
		if !(r == ',' ||
			r == '.' ||
			(r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}