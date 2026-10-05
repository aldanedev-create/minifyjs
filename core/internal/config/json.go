package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FileName is the config file name MinifyJS searches for, in order of
// precedence, when walking from the current directory up to the
// filesystem root.
//
// Only one file is loaded; the first one found wins. This matches
// the convention used by ESLint (.eslintrc.*) and Prettier
// (.prettierrc.*) — a project's config is discovered, not merged.
var FileNames = []string{
	"minifyjs.config.json",
	".minifyjsrc.json",
	".minifyjsrc",
}

// jsonConfig mirrors Config but with pointer fields so the loader can
// distinguish "field absent from JSON" from "field set to the zero
// value". Without pointers, a config that says
// `"cache": false` would be indistinguishable from one that omits
// `cache` entirely.
type jsonConfig struct {
	Inputs        *[]string      `json:"inputs,omitempty"`
	Output        *string        `json:"output,omitempty"`
	Minify        *jsonMinify    `json:"minify,omitempty"`
	Target        *string        `json:"target,omitempty"`
	Format        *string        `json:"format,omitempty"`
	Sourcemap     *string        `json:"sourcemap,omitempty"`
	Banner        *string        `json:"banner,omitempty"`
	Footer        *string        `json:"footer,omitempty"`
	LegalComments *string        `json:"legalComments,omitempty"`
	Bundle        *jsonBundle    `json:"bundle,omitempty"`
	Cache         *jsonCache     `json:"cache,omitempty"`
	Quiet         *bool          `json:"quiet,omitempty"`
	Verbose       *bool          `json:"verbose,omitempty"`
}

type jsonMinify struct {
	Whitespace  *bool `json:"whitespace,omitempty"`
	Identifiers *bool `json:"identifiers,omitempty"`
	Syntax      *bool `json:"syntax,omitempty"`
}

type jsonBundle struct {
	Enabled     *bool    `json:"enabled,omitempty"`
	EntryPoints []string `json:"entryPoints,omitempty"`
	Platform    *string  `json:"platform,omitempty"`
	Splitting   *bool    `json:"splitting,omitempty"`
}

type jsonCache struct {
	Enabled *bool   `json:"enabled,omitempty"`
	Dir     *string `json:"dir,omitempty"`
}

// Find walks up from start looking for a config file. It returns the
// path and true if one was found, or "" and false otherwise. Errors
// other than fs.ErrNotExist are returned to the caller.
func Find(start string) (string, bool, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", false, err
	}
	for {
		for _, name := range FileNames {
			candidate := filepath.Join(dir, name)
			info, err := os.Stat(candidate)
			if err == nil && info.Mode().IsRegular() {
				return candidate, true, nil
			}
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				return "", false, err
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false, nil
		}
		dir = parent
	}
}

// Load reads a config file and returns the Config it describes.
// Unset fields in the file are left at their zero value; the caller
// is expected to have started from Defaults() and to overlay the
// returned Config on top.
//
// Paths inside the file are resolved relative to the file's own
// directory, which lets a config checked into a monorepo work
// regardless of where `minifyjs` is invoked from.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("config: read %s: %w", path, err)
	}
	return Parse(data, filepath.Dir(path))
}

// Parse decodes config JSON. baseDir is the directory relative paths
// in the config are resolved against — normally the directory of the
// config file, or the current working directory when the config came
// from elsewhere.
func Parse(data []byte, baseDir string) (Config, error) {
	var raw jsonConfig
	if err := json.Unmarshal(data, &raw); err != nil {
		return Config{}, fmt.Errorf("config: parse: %w", err)
	}
	return raw.toConfig(baseDir)
}

func (j jsonConfig) toConfig(baseDir string) (Config, error) {
	var c Config

	if j.Inputs != nil {
		c.Inputs = resolvePaths(*j.Inputs, baseDir)
	}
	if j.Output != nil {
		c.Output = resolvePath(*j.Output, baseDir)
	}
	if j.Minify != nil {
		if j.Minify.Whitespace != nil {
			c.Minify.Whitespace = *j.Minify.Whitespace
		}
		if j.Minify.Identifiers != nil {
			c.Minify.Identifiers = *j.Minify.Identifiers
		}
		if j.Minify.Syntax != nil {
			c.Minify.Syntax = *j.Minify.Syntax
		}
	}
	if j.Target != nil {
		c.Target = *j.Target
	}
	if j.Format != nil {
		c.Format = *j.Format
	}
	if j.Sourcemap != nil {
		c.Sourcemap = *j.Sourcemap
	}
	if j.Banner != nil {
		c.Banner = *j.Banner
	}
	if j.Footer != nil {
		c.Footer = *j.Footer
	}
	if j.LegalComments != nil {
		c.LegalComments = *j.LegalComments
	}
	if j.Bundle != nil {
		if j.Bundle.Enabled != nil {
			c.Bundle = *j.Bundle.Enabled
		}
		c.BundleOptions.EntryPoints = resolvePaths(j.Bundle.EntryPoints, baseDir)
		if j.Bundle.Platform != nil {
			c.BundleOptions.Platform = *j.Bundle.Platform
		}
		if j.Bundle.Splitting != nil {
			c.BundleOptions.Splitting = *j.Bundle.Splitting
		}
	}
	if j.Cache != nil {
		if j.Cache.Enabled != nil {
			c.Cache = *j.Cache.Enabled
		}
		if j.Cache.Dir != nil {
			c.CacheDir = resolvePath(*j.Cache.Dir, baseDir)
		}
	}
	if j.Quiet != nil {
		c.Quiet = *j.Quiet
	}
	if j.Verbose != nil {
		c.Verbose = *j.Verbose
	}
	return c, nil
}

func resolvePaths(paths []string, baseDir string) []string {
	if len(paths) == 0 {
		return nil
	}
	out := make([]string, len(paths))
	for i, p := range paths {
		out[i] = resolvePath(p, baseDir)
	}
	return out
}

func resolvePath(p, baseDir string) string {
	if p == "" || filepath.IsAbs(p) {
		return p
	}
	return filepath.Clean(filepath.Join(baseDir, p))
}

// EnvOverlay applies MINIFYJS_* environment variables on top of c.
// It is called by the CLI between loading the file config and
// applying CLI flags, which gives the standard precedence:
// flags > env > file > defaults.
//
// Recognized variables:
//
//	MINIFYJS_TARGET         e.g. es2015
//	MINIFYJS_FORMAT         iife | cjs | esm
//	MINIFYJS_SOURCEMAP      inline | external | both
//	MINIFYJS_CACHE          1 | true | yes
//	MINIFYJS_CACHE_DIR      path
//	MINIFYJS_QUIET          1 | true | yes
//	MINIFYJS_VERBOSE        1 | true | yes
//	MINIFYJS_NO_MINIFY      1 | true | yes  (turns off all minify passes)
func (c *Config) EnvOverlay(getenv func(string) string) {
	if v := strings.TrimSpace(getenv("MINIFYJS_TARGET")); v != "" {
		c.Target = v
	}
	if v := strings.TrimSpace(getenv("MINIFYJS_FORMAT")); v != "" {
		c.Format = v
	}
	if v := strings.TrimSpace(getenv("MINIFYJS_SOURCEMAP")); v != "" {
		c.Sourcemap = v
	}
	if truthy(getenv("MINIFYJS_CACHE")) {
		c.Cache = true
	}
	if v := strings.TrimSpace(getenv("MINIFYJS_CACHE_DIR")); v != "" {
		c.CacheDir = v
	}
	if truthy(getenv("MINIFYJS_QUIET")) {
		c.Quiet = true
	}
	if truthy(getenv("MINIFYJS_VERBOSE")) {
		c.Verbose = true
	}
	if truthy(getenv("MINIFYJS_NO_MINIFY")) {
		c.Minify = MinifyOptions{}
	}
}

func truthy(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}