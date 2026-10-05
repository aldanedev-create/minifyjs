package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/minifyjs/minifyjs/core/api"
	"github.com/minifyjs/minifyjs/core/internal/config"
)

// runWatch implements --watch. It polls the input files (and any
// files they import, once bundling is in play) for mtime changes and
// re-runs the transform on each change.
//
// Polling rather than fsnotify because:
//
//   - fsnotify does not work on every filesystem (notably some
//     network mounts and Docker volume drivers).
//   - Watch mode is a convenience for local development, not a
//     CI-critical feature; a 250ms poll is invisible to humans.
//
// If watch mode becomes performance-critical, swapping in fsnotify
// behind this interface is a small change.
func runWatch(cfg config.Config, stdout, stderr io.Writer) int {
	if len(cfg.Inputs) == 0 {
		fmt.Fprintln(stderr, "minifyjs: --watch requires at least one input file")
		return exitUsage
	}
	if cfg.Output == "" {
		fmt.Fprintln(stderr, "minifyjs: --watch requires --output (watched output cannot go to stdout)")
		return exitUsage
	}

	mtimes := make(map[string]time.Time)
	runOnce := func() int {
		var source []byte
		for _, in := range cfg.Inputs {
			b, err := readFileAndTouch(in, mtimes)
			if err != nil {
				fmt.Fprintf(stderr, "minifyjs: %v\n", err)
				return exitError
			}
			// For the simple non-bundle case we only support one input.
			source = b
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
		}
		res, err := api.Minify(string(source), opts)
		if err != nil {
			fmt.Fprintf(stderr, "minifyjs: %v\n", err)
			return exitError
		}
		if len(res.Diagnostics) > 0 {
			renderDiagnostics(stderr, res.Diagnostics, cfg.Quiet)
		}
		if res.HasErrors() {
			return exitError
		}
		if err := writeFile(cfg.Output, []byte(res.Code)); err != nil {
			fmt.Fprintf(stderr, "minifyjs: %v\n", err)
			return exitError
		}
		if !cfg.Quiet {
			fmt.Fprintf(stderr, "[watch] rebuilt %s at %s\n",
				cfg.Output, time.Now().Format("15:04:05"))
		}
		return exitOK
	}

	// Initial run.
	if code := runOnce(); code != exitOK {
		// Keep watching even if the first run failed; a user is
		// likely editing toward a correct state.
	}

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		changed := false
		for _, in := range cfg.Inputs {
			info, err := os.Stat(in)
			if err != nil {
				continue
			}
			prev, seen := mtimes[in]
			if !seen || info.ModTime().After(prev) {
				mtimes[in] = info.ModTime()
				changed = true
			}
		}
		if changed {
			runOnce()
		}
	}
	return exitOK
}

func readFileAndTouch(path string, mtimes map[string]time.Time) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	mtimes[path] = info.ModTime()
	return os.ReadFile(filepath.Clean(path))
}