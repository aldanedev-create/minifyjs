package config

import (
	"fmt"
	"strings"
)

// Validate checks that the resolved config is internally consistent.
// It runs after defaults, file, env, and flags have all been merged,
// so it sees the final values the CLI will act on.
//
// Validate does not check that input files exist — that is the CLI's
// job when it opens them — but it does check that a target string is
// well-formed enough to pass to esbuild, and that mutually exclusive
// options are not both set.
func (c Config) Validate() error {
	if err := validateTarget(c.Target); err != nil {
		return err
	}
	switch c.Format {
	case "", "iife", "cjs", "esm":
	default:
		return fmt.Errorf("config: invalid format %q (want iife, cjs, or esm)", c.Format)
	}
	switch c.Sourcemap {
	case "", "inline", "external", "both":
	default:
		return fmt.Errorf("config: invalid sourcemap %q (want inline, external, or both)", c.Sourcemap)
	}
	switch c.LegalComments {
	case "", "none", "inline", "eof", "external":
	default:
		return fmt.Errorf("config: invalid legalComments %q", c.LegalComments)
	}
	if c.Quiet && c.Verbose {
		return fmt.Errorf("config: quiet and verbose cannot both be set")
	}
	if c.Bundle && len(c.BundleOptions.EntryPoints) == 0 && len(c.Inputs) == 0 {
		return fmt.Errorf("config: bundle requires bundle.entryPoints or inputs")
	}
	if c.BundleOptions.Splitting && (!c.Bundle || c.Format != "esm") {
		return fmt.Errorf("config: bundle.splitting requires bundle=true and format=esm")
	}
	return nil
}

// validateTarget accepts the forms esbuild accepts: an ECMAScript
// keyword (es5, es2015..es2024, esnext), a browser+version pair
// (chrome58), or a comma-separated list of those.
func validateTarget(t string) error {
	if t == "" {
		return nil
	}
	for _, part := range strings.Split(t, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			return fmt.Errorf("config: empty target component in %q", t)
		}
		for _, r := range part {
			if !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') && r != '.' {
				return fmt.Errorf("config: invalid target component %q", part)
			}
		}
	}
	return nil
}