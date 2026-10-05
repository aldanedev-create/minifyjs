// Package config loads and validates MinifyJS configuration from a
// JSON file (conventionally minifyjs.config.json) and provides the
// precedence rules for combining file config, environment variables,
// and CLI flags.
//
// The CLI consults this package once, at startup, before parsing any
// source. Every other package receives a Config value; nothing else
// reads files from disk for options.
//
// Precedence (highest wins):
//
//  1. CLI flags passed to `minifyjs`
//  2. Environment variables (MINIFYJS_*)
//  3. A config file found via the search path
//  4. Built-in defaults (defaults.go)
//
// This matches the convention used by eslint, prettier, and similar
// tools, so users do not have to learn a new override model.
package config

// Config is the fully-resolved configuration used by the CLI and the
// public Go API. Every field has a zero value that means "use the
// default from Defaults()", except where noted.
type Config struct {
	// Input files. Empty means read from stdin.
	Inputs []string

	// Output path. Empty means write to stdout.
	Output    string
	OutputDir bool

	// Minify controls the three esbuild minification passes.
	Minify    MinifyOptions
	MinifySet bool

	// Target is the ECMAScript version for output. Empty means
	// esnext (no lowering).
	Target string

	// Format is the module wrapper for output. Empty means
	// "preserve".
	Format string

	// Sourcemap controls source-map emission. Valid values:
	// "", "inline", "external", "both".
	Sourcemap string

	// Banner is prepended verbatim to the output.
	Banner string

	// Footer is appended verbatim to the output.
	Footer string

	// LegalComments controls what happens to /*! ... */ comments.
	// Valid values: "", "none", "inline", "eof", "external".
	LegalComments string
	Define        map[string]string
	Drop          []string
	Pure          []string

	// Bundle enables module resolution and bundling. Requires
	// BundleOptions.EntryPoints or a single entry in Inputs.
	Bundle bool

	// BundleOptions are only consulted when Bundle is true.
	BundleOptions BundleOptions

	// Cache enables the on-disk result cache. Off by default so a
	// first-time user's runs are as simple as possible.
	Cache bool

	// CacheDir overrides the default cache location. Empty means
	// the OS-specific default from cache.DefaultDir().
	CacheDir string

	// Quiet suppresses progress output on stderr. Errors still
	// print.
	Quiet bool

	// Verbose prints extra information on stderr. Mutually
	// exclusive with Quiet; Verbose wins if both are set.
	Verbose bool
}

// MinifyOptions mirrors engine.Options's three minify switches. Kept
// as its own struct so config does not import engine (engine is
// imported by the CLI, not by config).
type MinifyOptions struct {
	Whitespace  bool
	Identifiers bool
	Syntax      bool
}

// BundleOptions mirrors engine.BuildOptions for the subset that makes
// sense in a config file. File paths are resolved relative to the
// config file's directory so a config checked into a repo works no
// matter where the user invokes `minifyjs` from.
type BundleOptions struct {
	EntryNames  string
	ChunkNames  string
	AssetNames  string
	Metafile    string
	External    []string
	Packages    string
	TreeShaking string
	Charset     string

	WorkingDir string

	EntryPoints []string
	Platform    string // "browser" | "node" | "neutral"
	Splitting   bool
}

// AllMinify is a convenience that turns on all three minification
// passes. It is what `--compress` sets.
func (c *Config) AllMinify() {
	c.Minify.Whitespace = true
	c.Minify.Identifiers = true
	c.Minify.Syntax = true
}

// IsMinifying reports whether any minification pass is enabled.
func (c Config) IsMinifying() bool {
	return c.Minify.Whitespace || c.Minify.Identifiers || c.Minify.Syntax
}
