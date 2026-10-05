package config

// Defaults returns the built-in configuration used when no config
// file is found and no flags are given.
//
// Design notes:
//
//   - Minification is ON by default. A user who runs `minifyjs app.js`
//     wants a smaller app.js, not a copy. Users who want a pure
//     passthrough can pass --no-minify.
//
//   - Target is empty (esnext) by default. Lowering is a compatibility
//     choice the user makes; guessing a target silently produces
//     output that may not run where they expect.
//
//   - Format is empty (preserve) by default. The module wrapper is
//     determined by the input.
//
//   - Sourcemap is empty (none) by default. Source maps cost a file
//     and are a deliberate choice.
//
//   - Cache is off by default. It writes to disk; users should opt in
//     the first time they see it.
func Defaults() Config {
	return Config{
		Minify: MinifyOptions{
			Whitespace:  true,
			Identifiers: true,
			Syntax:      true,
		},
		Target:  "",
		Format:  "",
		Cache:   false,
		Quiet:   false,
		Verbose: false,
	}
}