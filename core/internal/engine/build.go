package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

// BuildOptions configures a multi-file build (bundling). It is
// separate from Options because bundling has fundamentally different
// inputs and outputs than a single-file transform: it needs an entry
// point on disk, it produces one or more output files, and it may
// resolve modules from node_modules even though MinifyJS does not
// require Node.js at install time (module resolution is a filesystem
// operation, not a Node runtime dependency).
//
// Build is the Phase 9 feature from the blueprint. It is not used by
// the default `minifyjs file.js` path; it is reached only via
// `minifyjs --bundle src/main.js`.
type BuildOptions struct {
	// EntryPoints are the files to bundle. At least one is required.
	EntryPoints []string
	EntryNames  string
	ChunkNames  string
	AssetNames  string
	Metafile    string
	External    []string
	Packages    string
	TreeShaking string
	Charset     string

	// OutDir is the directory to write output files into. Either
	// OutDir or OutFile must be set.
	OutDir string

	// OutFile is the path of the single output file. Cannot be set
	// together with OutDir.
	OutFile string

	// AbsWorkingDir is the directory all relative paths are resolved
	// against. Defaults to the current working directory.
	AbsWorkingDir string

	// Platform is "browser", "node", or "neutral".
	Platform string

	// Bundle enables module resolution and dependency graph
	// construction. Without it, Build behaves like Transform with
	// the extra file-level options.
	Bundle bool

	// Splitting enables code splitting. Only valid with
	// Format == "esm" and Bundle == true.
	Splitting bool

	// Format controls the module wrapper for the build.
	Format string

	// Options carries the same minify/target/format/sourcemap
	// settings as a single-file transform.
	Options
}

// Build runs a bundling build. Unlike Transform, Build writes to
// the filesystem; callers are responsible for creating OutDir first
// if they set OutDir rather than OutFile.
//
// The returned Result.Code is empty for Build because the output is
// on disk. Callers who need the bundled code in memory should read
// the output file themselves or use Transform on each module.
func Build(opts BuildOptions) (Result, error) {
	if len(opts.EntryPoints) == 0 {
		return Result{}, fmt.Errorf("engine: Build requires at least one entry point")
	}
	if opts.OutDir == "" && opts.OutFile == "" {
		return Result{}, fmt.Errorf("engine: Build requires OutDir or OutFile")
	}
	if opts.OutDir != "" && opts.OutFile != "" {
		return Result{}, fmt.Errorf("engine: Build cannot set both OutDir and OutFile")
	}
	options := opts.Options
	if opts.Format != "" {
		options.Format = opts.Format
	}
	if err := options.validate(); err != nil {
		return Result{}, err
	}
	if opts.Splitting && (options.Format != "esm" || !opts.Bundle) {
		return Result{}, fmt.Errorf("engine: Splitting requires Format=\"esm\" and Bundle=true")
	}

	bOpts := api.BuildOptions{
		EntryPoints:   opts.EntryPoints,
		Outdir:        opts.OutDir,
		Outfile:       opts.OutFile,
		AbsWorkingDir: opts.AbsWorkingDir,
		Bundle:        opts.Bundle,
		Splitting:     opts.Splitting,
		Write:         true,
		EntryNames:    opts.EntryNames,
		ChunkNames:    opts.ChunkNames,
		AssetNames:    opts.AssetNames,
		External:      opts.External,
		Metafile:      opts.Metafile != "",
		Charset:       api.CharsetUTF8,

		MinifyWhitespace:  opts.MinifyWhitespace,
		MinifyIdentifiers: opts.MinifyIdentifiers,
		MinifySyntax:      opts.MinifySyntax,
		Define:            opts.Define,
		Pure:              opts.Pure,

		LogLevel: api.LogLevelSilent, // we surface diagnostics ourselves
	}
	switch opts.Packages {
	case "", "bundle":
	case "external":
		bOpts.Packages = api.PackagesExternal
	default:
		return Result{}, fmt.Errorf("invalid packages mode: %s", opts.Packages)
	}
	switch opts.TreeShaking {
	case "", "default":
	case "true":
		bOpts.TreeShaking = api.TreeShakingTrue
	case "false":
		bOpts.TreeShaking = api.TreeShakingFalse
	default:
		return Result{}, fmt.Errorf("invalid tree-shaking mode: %s", opts.TreeShaking)
	}
	switch opts.Charset {
	case "", "utf8":
	case "ascii":
		bOpts.Charset = api.CharsetASCII
	default:
		return Result{}, fmt.Errorf("invalid charset: %s", opts.Charset)
	}
	for _, kind := range opts.Drop {
		switch kind {
		case "console":
			bOpts.Drop |= api.DropConsole
		case "debugger":
			bOpts.Drop |= api.DropDebugger
		}
	}

	if opts.Platform != "" {
		switch opts.Platform {
		case "browser":
			bOpts.Platform = api.PlatformBrowser
		case "node":
			bOpts.Platform = api.PlatformNode
		case "neutral":
			bOpts.Platform = api.PlatformNeutral
		}
	}
	if opts.Target != "" {
		bOpts.Target = mapTarget(opts.Target)
	}
	if options.Format != "" {
		bOpts.Format = api.Format(mapFormat(options.Format))
	}
	if opts.Sourcemap != "" {
		bOpts.Sourcemap = mapSourcemap(opts.Sourcemap)
	}
	if opts.Banner != "" {
		bOpts.Banner = map[string]string{"js": opts.Banner}
	}
	if opts.Footer != "" {
		bOpts.Footer = map[string]string{"js": opts.Footer}
	}

	if options.LegalComments != "" {
		bOpts.LegalComments = mapLegalComments(options.LegalComments)
	}
	res := api.Build(bOpts)
	if len(res.Errors) == 0 && opts.Metafile != "" {
		path := opts.Metafile
		if !filepath.IsAbs(path) {
			path = filepath.Join(opts.AbsWorkingDir, path)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return Result{}, err
		}
		if err := os.WriteFile(path, []byte(res.Metafile), 0644); err != nil {
			return Result{}, err
		}
	}

	out := Result{
		Diagnostics: convertMessages(res.Errors, res.Warnings),
	}
	if len(res.OutputFiles) > 0 {
		// When OutFile is used, esbuild returns the output in
		// OutputFiles[0].Contents even though it also wrote it to
		// disk. We surface that so callers can compare against
		// what is on disk if they want.
		for _, f := range res.OutputFiles {
			if f.Path == opts.OutFile {
				out.Code = string(f.Contents)
				out.MinifiedBytes = len(f.Contents)
			}
		}
	}
	return out, nil
}
