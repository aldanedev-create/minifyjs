package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/minifyjs/minifyjs/core/api"
	"github.com/minifyjs/minifyjs/core/internal/config"
	"github.com/minifyjs/minifyjs/core/internal/diagnostics"
	miniio "github.com/minifyjs/minifyjs/core/internal/io"
)

// writer is the type of stderr/stdout. Aliased so tests can swap it.
type writer = io.Writer

// Exit codes. These are stable and documented in
// docs/reference/exit-codes.md.
const (
	exitOK       = 0 // success
	exitError    = 1 // processing error (bad input, esbuild error, I/O error)
	exitUsage    = 2 // bad command-line usage
	exitNotFound = 3 // input file not found
	exitInternal = 4 // unexpected internal error
)

// run is the real entry point. It returns an exit code rather than
// calling os.Exit so tests can invoke it directly.
func run(args []string) int {
	return runWith(args, os.Stdin, os.Stdout, os.Stderr, os.Getenv)
}

// runWith is run with explicit streams and environment, for tests.
func runWith(
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
	getenv func(string) string,
) int {
	f, err := parseFlags(args)
	if err != nil {
		fmt.Fprintf(stderr, "minifyjs: %v\n", err)
		fmt.Fprintln(stderr, "Run 'minifyjs --help' for usage.")
		return exitUsage
	}

	if f.showHelp {
		printHelp(stdout)
		return exitOK
	}
	if f.showVersion {
		printVersion(stdout)
		return exitOK
	}

	cfg, code := resolveConfig(f, getenv, stderr)
	if code != exitOK {
		return code
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(stderr, "minifyjs: %v\n", err)
		return exitUsage
	}

	if f.watch {
		return runWatch(cfg, stdout, stderr)
	}

	if cfg.Bundle {
		return runBundle(cfg, stderr)
	}
	return runTransform(cfg, stdin, stdout, stderr)
}

// resolveConfig merges, in order of increasing precedence:
// built-in defaults, a discovered config file, environment variables,
// and the parsed flags.
func resolveConfig(f *flags, getenv func(string) string, stderr io.Writer) (config.Config, int) {
	cfg := config.Defaults()

	// 1. Config file (unless --no-config).
	if !f.noConfig {
		searchDir := "."
		if f.configPathSet {
			searchDir = filepath.Dir(f.configPath)
		}
		var cfgPath string
		var found bool
		var err error
		if f.configPathSet {
			cfgPath = f.configPath
			found = true
		} else {
			cfgPath, found, err = config.Find(searchDir)
			if err != nil {
				fmt.Fprintf(stderr, "minifyjs: config: %v\n", err)
				return config.Config{}, exitError
			}
		}
		if found {
			fileCfg, err := config.Load(cfgPath)
			if err != nil {
				fmt.Fprintf(stderr, "minifyjs: %v\n", err)
				return config.Config{}, exitError
			}
			cfg = overlay(cfg, fileCfg)
		}
	}

	// 2. Environment.
	cfg.EnvOverlay(getenv)

	// 3. Flags.
	applyFlags(&cfg, f)

	return cfg, exitOK
}

// overlay applies any field in from that is non-zero onto base.
// Because config.Load returns a Config with only the file's explicit
// fields set, non-zero is a reliable signal for every field except
// those whose valid zero value is meaningful (bools, strings). For
// those, the loader uses pointer fields to distinguish "unset" from
// "false/empty"; by the time we get here, unset means zero.
func overlay(base, from config.Config) config.Config {
	if from.Inputs != nil {
		base.Inputs = from.Inputs
	}
	if from.Output != "" {
		base.Output = from.Output
	}
	base.Minify = from.Minify // file minify block fully replaces
	if from.Target != "" {
		base.Target = from.Target
	}
	if from.Format != "" {
		base.Format = from.Format
	}
	if from.Sourcemap != "" {
		base.Sourcemap = from.Sourcemap
	}
	if from.Banner != "" {
		base.Banner = from.Banner
	}
	if from.Footer != "" {
		base.Footer = from.Footer
	}
	if from.LegalComments != "" {
		base.LegalComments = from.LegalComments
	}
	if from.Bundle {
		base.Bundle = true
	}
	if from.BundleOptions.EntryPoints != nil {
		base.BundleOptions.EntryPoints = from.BundleOptions.EntryPoints
	}
	if from.BundleOptions.Platform != "" {
		base.BundleOptions.Platform = from.BundleOptions.Platform
	}
	if from.BundleOptions.Splitting {
		base.BundleOptions.Splitting = true
	}
	if from.Cache {
		base.Cache = true
	}
	if from.CacheDir != "" {
		base.CacheDir = from.CacheDir
	}
	if from.Quiet {
		base.Quiet = true
	}
	if from.Verbose {
		base.Verbose = true
	}
	return base
}

// applyFlags overlays the parsed flags onto cfg. Every flag that has
// a corresponding *Set boolean is only applied when that boolean is
// true; this is what allows --no-compress on the command line to
// override a config file's "compress: true".
func applyFlags(cfg *config.Config, f *flags) {
	if len(f.inputs) > 0 {
		cfg.Inputs = f.inputs
	}
	if f.outputSet {
		cfg.Output = f.output
	}
	if f.noMinify {
		cfg.Minify = config.MinifyOptions{}
	} else {
		if f.compressSet {
			if f.compress {
				cfg.Minify.Whitespace = true
				cfg.Minify.Syntax = true
			} else {
				cfg.Minify.Whitespace = false
				cfg.Minify.Syntax = false
			}
		}
		if f.mangleSet {
			cfg.Minify.Identifiers = f.mangle
		}
	}
	if f.targetSet {
		cfg.Target = f.target
	}
	if f.formatSet {
		cfg.Format = f.format
	}
	if f.smSet {
		cfg.Sourcemap = f.sourcemap
	}
	if f.bannerSet {
		cfg.Banner = f.banner
	}
	if f.footerSet {
		cfg.Footer = f.footer
	}
	if f.lcSet {
		cfg.LegalComments = f.legalComments
	}
	if f.bundleSet {
		cfg.Bundle = f.bundle
	}
	if f.outdirSet {
		cfg.BundleOptions.EntryPoints = f.inputs
		cfg.Output = f.outdir
	}
	if f.platformSet {
		cfg.BundleOptions.Platform = f.platform
	}
	if f.splitSet {
		cfg.BundleOptions.Splitting = f.splitting
	}
	if f.cacheSet {
		cfg.Cache = f.cache
	}
	if f.cacheDirSet {
		cfg.CacheDir = f.cacheDir
	}
	if f.quietSet {
		cfg.Quiet = f.quiet
	}
	if f.verboseSet {
		cfg.Verbose = f.verbose
	}
}

// runTransform handles the single-file / stdin path.
func runTransform(cfg config.Config, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(cfg.Inputs) > 1 {
		fmt.Fprintln(stderr, "minifyjs: only one input file allowed (did you mean --bundle?)")
		return exitUsage
	}

	var (
		source    []byte
		sourceName string
	)
	switch {
	case len(cfg.Inputs) == 0:
		// stdin
		b, err := miniio.ReadAll(stdin)
		if err != nil {
			fmt.Fprintf(stderr, "minifyjs: read stdin: %v\n", err)
			return exitError
		}
		source = b
		sourceName = "<stdin>"
	case cfg.Inputs[0] == "-":
		b, err := miniio.ReadAll(stdin)
		if err != nil {
			fmt.Fprintf(stderr, "minifyjs: read stdin: %v\n", err)
			return exitError
		}
		source = b
		sourceName = "<stdin>"
	default:
		b, err := miniio.ReadFile(cfg.Inputs[0])
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(stderr, "minifyjs: %s: no such file\n", cfg.Inputs[0])
				return exitNotFound
			}
			fmt.Fprintf(stderr, "minifyjs: %v\n", err)
			return exitError
		}
		source = b
		sourceName = cfg.Inputs[0]
	}

	opts := api.Options{
		MinifyWhitespace:  cfg.Minify.Whitespace,
		MinifyIdentifiers: cfg.Minify.Identifiers,
		MinifySyntax:      cfg.Minify.Syntax,
		Target:            cfg.Target,
		Format:            cfg.Format,
		Sourcemap:         cfg.Sourcemap,
		Banner:            cfg.Banner,
		Footer:            cfg.Footer,
		LegalComments:     cfg.LegalComments,
		SourceName:        sourceName,
	}

	result, err := api.Minify(string(source), opts)
	if err != nil {
		fmt.Fprintf(stderr, "minifyjs: %v\n", err)
		return exitError
	}

	if len(result.Diagnostics) > 0 {
		renderDiagnostics(stderr, result.Diagnostics, cfg.Quiet)
	}
	if diagnostics.HasErrors(result.Diagnostics) {
		return exitError
	}

	if err := writeOutput(cfg.Output, []byte(result.Code), stdout); err != nil {
		fmt.Fprintf(stderr, "minifyjs: %v\n", err)
		return exitError
	}

	if !cfg.Quiet && cfg.Output != "" {
		reportSavings(stderr, sourceName, cfg.Output, result.OriginalBytes, result.MinifiedBytes)
	}
	return exitOK
}

// runBundle handles --bundle.
func runBundle(cfg config.Config, stderr io.Writer) int {
	entries := cfg.BundleOptions.EntryPoints
	if len(entries) == 0 {
		entries = cfg.Inputs
	}
	if len(entries) == 0 {
		fmt.Fprintln(stderr, "minifyjs: --bundle requires at least one entry point")
		return exitUsage
	}
	if cfg.Output == "" {
		fmt.Fprintln(stderr, "minifyjs: --bundle requires --outdir or --outfile")
		return exitUsage
	}

	platform := cfg.BundleOptions.Platform
	if platform == "" {
		platform = "browser"
	}

	bopts := api.BundleOptions{
		EntryPoints:   entries,
		Bundle:        true,
		Platform:      platform,
		Splitting:     cfg.BundleOptions.Splitting,
		AbsWorkingDir: mustGetwd(),
		Options: api.Options{
			MinifyWhitespace:  cfg.Minify.Whitespace,
			MinifyIdentifiers: cfg.Minify.Identifiers,
			MinifySyntax:      cfg.Minify.Syntax,
			Target:            cfg.Target,
			Format:            cfg.Format,
			Sourcemap:         cfg.Sourcemap,
			Banner:            cfg.Banner,
			Footer:            cfg.Footer,
			LegalComments:     cfg.LegalComments,
		},
	}

	// Decide whether Output is a dir or a file.
	if isDirLike(cfg.Output) {
		bopts.OutDir = cfg.Output
	} else {
		bopts.OutFile = cfg.Output
	}

	result, err := api.Bundle(bopts)
	if err != nil {
		fmt.Fprintf(stderr, "minifyjs: %v\n", err)
		return exitError
	}
	if len(result.Diagnostics) > 0 {
		renderDiagnostics(stderr, result.Diagnostics, cfg.Quiet)
	}
	if diagnostics.HasErrors(result.Diagnostics) {
		return exitError
	}

	if !cfg.Quiet {
		fmt.Fprintf(stderr, "bundled %d entry point(s) -> %s\n", len(entries), cfg.Output)
	}
	return exitOK
}

// writeOutput sends the minified code to the configured destination.
func writeOutput(path string, code []byte, stdout io.Writer) error {
	if path == "" || path == "-" {
		_, err := stdout.Write(code)
		return err
	}
	return miniio.WriteFileAtomic(path, code, 0o644)
}

// renderDiagnostics prints diagnostics to stderr unless quiet.
func renderDiagnostics(w io.Writer, diags []diagnostics.Diagnostic, quiet bool) {
	if quiet {
		// Still print errors, but skip warnings/info.
		filtered := make([]diagnostics.Diagnostic, 0, len(diags))
		for _, d := range diags {
			if d.IsError() {
				filtered = append(filtered, d)
			}
		}
		if len(filtered) == 0 {
			return
		}
		diags = filtered
	}
	f := diagnostics.TerminalFormatter{Prefix: "minifyjs: "}
	_ = f.Format(w, diags)
}

func reportSavings(w io.Writer, in, out string, origBytes, minBytes int) {
	saved := origBytes - minBytes
	pct := 0.0
	if origBytes > 0 {
		pct = float64(saved) * 100 / float64(origBytes)
	}
	fmt.Fprintf(w, "%s -> %s: %d -> %d bytes (%.1f%% smaller)\n",
		in, out, origBytes, minBytes, pct)
}

func isDirLike(path string) bool {
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, string(os.PathSeparator)) {
		return true
	}
	if info, err := os.Stat(path); err == nil {
		return info.IsDir()
	}
	return false
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

// printHelp writes the --help text. It is generated from a static
// string rather than reflected from flags so the wording can be
// curated; a future tool (tools/docs/generate_cli_reference.py)
// cross-checks the two.
func printHelp(w io.Writer) {
	fmt.Fprint(w, `MinifyJS - a native JavaScript minifier (no Node.js required)

Usage:
  minifyjs [input.js] [flags]
  cat input.js | minifyjs [flags]
  minifyjs --bundle src/main.js --outdir dist [flags]

Input:
  input.js              Input file. Omit or use "-" to read stdin.

Output:
  -o, --output PATH     Write output to PATH. Default: stdout.

Minify (default: on):
  --compress            Enable whitespace + syntax minification.
  --no-compress         Disable whitespace + syntax minification.
  --mangle              Enable identifier mangling.
  --no-mangle           Disable identifier mangling.
  --no-minify           Disable all minification passes.

Target:
  --target VERSION      Output target: es5, es2015..es2024, esnext,
                        or a comma list (es2020,chrome90).
  --format FORMAT       Output format: iife, cjs, esm.
  --sourcemap MODE      Source map: inline, external, both.

Content:
  --banner TEXT         Prepend TEXT to output (e.g. license header).
  --footer TEXT         Append TEXT to output.
  --legal-comments MODE Legal comments: none, inline, eof, external.
  --define K=V          Replace identifier K with expression V.
  --drop KIND           Remove statements: console, debugger.
  --pure FUNC           Mark FUNC as side-effect-free for tree-shaking.

Bundling:
  --bundle              Resolve imports and bundle into one file.
  --outdir DIR          Output directory for bundled files.
  --platform P          Target platform: browser, node, neutral.
  --splitting           Enable code splitting (requires --format esm).

Caching:
  --cache               Enable on-disk result cache.
  --no-cache            Disable cache.
  --cache-dir DIR       Override cache directory.

Config:
  --config PATH         Use PATH as config file.
  --no-config           Ignore any config file.

Mode:
  --watch               Re-run on input file changes.
  -q, --quiet           Suppress non-error output.
  -v, --verbose         Print extra information.
  -h, --help            Print this help and exit.
  -V, --version         Print version and exit.

Exit codes:
  0  success
  1  processing error (bad input, esbuild error, I/O error)
  2  bad command-line usage
  3  input file not found
  4  unexpected internal error

Examples:
  minifyjs app.js -o app.min.js
  minifyjs app.js --compress --mangle -o app.min.js
  cat app.js | minifyjs --target es2015
  minifyjs --bundle src/main.js --outdir dist --format esm
`)
}