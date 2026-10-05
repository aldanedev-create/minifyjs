package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
		Charset:       api.CharsetUTF8,

		MinifyWhitespace:  opts.MinifyWhitespace,
		MinifyIdentifiers: opts.MinifyIdentifiers,
		MinifySyntax:      opts.MinifySyntax,
		Define:            opts.Define,
		Pure:              opts.Pure,

		LogLevel: api.LogLevelSilent, // we surface diagnostics ourselves
	}
	if options.Format == "esm" {
		bOpts.Packages = api.PackagesExternal
	}
	if hasReExport(opts.EntryPoints) {
		bOpts.TreeShaking = api.TreeShakingFalse
	}
	if !opts.Splitting {
		bOpts.External = dynamicImportPaths(opts.EntryPoints)
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

	res := api.Build(bOpts)

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

func dynamicImportPaths(entries []string) []string {
	pattern := regexp.MustCompile(`\bimport\s*\(\s*["']([^"']+)["']\s*\)`)
	seen := make(map[string]bool)
	var paths []string
	for _, entry := range entries {
		data, err := os.ReadFile(entry)
		if err != nil {
			continue
		}
		for _, match := range pattern.FindAllStringSubmatch(string(data), -1) {
			if strings.HasPrefix(match[1], ".") {
				target := filepath.Clean(filepath.Join(filepath.Dir(entry), match[1]))
				if _, err := os.Stat(target); err != nil {
					continue
				}
			}
			if !seen[match[1]] {
				seen[match[1]] = true
				paths = append(paths, match[1])
			}
		}
	}
	return paths
}

func hasReExport(entries []string) bool {
	pattern := regexp.MustCompile(`\bexport\s+(?:\*|\{[^}]+\})\s+from\s*["']`)
	seen := make(map[string]bool)
	for _, entry := range entries {
		root := filepath.Dir(entry)
		_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || seen[path] {
				return nil
			}
			seen[path] = true
			if filepath.Ext(path) != ".js" {
				return nil
			}
			data, readErr := os.ReadFile(path)
			if readErr == nil && pattern.Match(data) {
				seen["__found__"] = true
			}
			return nil
		})
		if seen["__found__"] {
			return true
		}
	}
	return false
}