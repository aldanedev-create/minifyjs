package main

// flags holds every command-line option. Fields mirror the config
// package's Config struct so the merge at the end of parsing is a
// simple field-by-field overlay.
//
// Every field has a "set" counterpart (e.g. Compress and compressSet)
// so the merge can distinguish "user did not pass this flag" from
// "user passed --no-compress". Without this, a config file's
// "compress: true" would be silently overridden by a false zero
// value.
type flags struct {
	// Positional
	inputs []string

	// Output
	output    string
	outputSet bool

	// Minify switches
	compress    bool
	compressSet bool
	mangle      bool
	mangleSet   bool
	noMinify    bool

	// Target / format / sourcemap
	target    string
	targetSet bool
	format    string
	formatSet bool
	sourcemap string
	smSet     bool

	// Banner / footer / legal comments
	banner        string
	bannerSet     bool
	footer        string
	footerSet     bool
	legalComments string
	lcSet         bool

	// Bundle
	bundle      bool
	bundleSet   bool
	outdir      string
	outdirSet   bool
	platform    string
	platformSet bool
	splitting   bool
	splitSet    bool

	bundleValues map[string]string
	externals    []string

	// Cache
	cache       bool
	cacheSet    bool
	cacheDir    string
	cacheDirSet bool

	// Config
	configPath    string
	configPathSet bool
	noConfig      bool

	// Verbosity
	quiet      bool
	quietSet   bool
	verbose    bool
	verboseSet bool

	// Watch
	watch bool

	// Meta
	showHelp    bool
	showVersion bool

	// Define / drop / pure (esbuild passthroughs)
	defines map[string]string
	drops   []string
	pures   []string
}

// newFlags returns an empty flags value. Defines is pre-allocated so
// the parser can write to it without checking nil.
func newFlags() *flags {
	return &flags{defines: map[string]string{}, bundleValues: map[string]string{}}
}
