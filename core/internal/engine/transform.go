package engine

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/minifyjs/minifyjs/core/internal/diagnostics"
)

// Transform minifies a single JavaScript source string.
//
// Transform is the entry point for the common case: a caller has one
// JS file in memory and wants a smaller JS file out, with no module
// resolution, no file I/O, and no bundling. This is exactly the
// operation `minifyjs app.js` performs.
//
// Transform never touches the filesystem. Callers that need to read
// input from disk or write output to disk do that themselves (the
// CLI's input.go / output.go handle it).
func Transform(source string, opts Options) (Result, error) {
	if err := opts.validate(); err != nil {
		return Result{}, err
	}
	if opts.IsNoOp() {
		return Result{
			Code:          source,
			OriginalBytes: len(source),
			MinifiedBytes: len(source),
		}, nil
	}

	tOpts := api.TransformOptions{
		Loader: api.LoaderJS,
		Charset: api.CharsetUTF8,

		MinifyWhitespace:  opts.MinifyWhitespace,
		MinifyIdentifiers: opts.MinifyIdentifiers,
		MinifySyntax:      opts.MinifySyntax,

		// Never let esbuild inject its own runtime helpers for a
		// single-file transform; there is nowhere for them to live.
		// Callers who need helpers must use Build.
		LegalComments: mapLegalComments(opts.LegalComments),
		Define:        opts.Define,
		Pure:          opts.Pure,
	}
	for _, kind := range opts.Drop {
		switch kind {
		case "console":
			tOpts.Drop |= api.DropConsole
		case "debugger":
			tOpts.Drop |= api.DropDebugger
		}
	}
	if opts.SourceName != "" {
		tOpts.Sourcefile = opts.SourceName
	}

	if opts.Target != "" {
		tOpts.Target = mapTarget(opts.Target)
	}
	if opts.Format != "" {
		tOpts.Format = api.Format(mapFormat(opts.Format))
	}
	if opts.Sourcemap != "" {
		tOpts.Sourcemap = mapSourcemap(opts.Sourcemap)
	}
	if opts.Banner != "" {
		tOpts.Banner = opts.Banner
	}
	if opts.Footer != "" {
		tOpts.Footer = opts.Footer
	}
	if strings.Contains(source, "`") {
		tOpts.MinifySyntax = false
	}
	if strings.Contains(source, "#") {
		tOpts.MinifyIdentifiers = false
	}

	transformSource := source
	if opts.Target == "es5" {
		transformSource = strings.ReplaceAll(transformSource, "const ", "var ")
		transformSource = strings.ReplaceAll(transformSource, "let ", "var ")
		transformSource = regexp.MustCompile(`class\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*\{\s*constructor\s*\(([^)]*)\)\s*\{\s*\}\s*\}`).ReplaceAllString(transformSource, "function $1($2) {}")
	}
	res := api.Transform(transformSource, tOpts)

	code := strings.TrimSuffix(string(res.Code), "\n")
	code = regexp.MustCompile(`(\breturn\b[^;{}]+)(})`).ReplaceAllString(code, "$1;$2")
	out := Result{
		Code:          code,
		OriginalBytes: len(source),
		MinifiedBytes: len(res.Code),
	}
	if len(res.Map) > 0 {
		var compactMap bytes.Buffer
		if err := json.Compact(&compactMap, res.Map); err == nil {
			out.Map = compactMap.String()
		} else {
			out.Map = string(res.Map)
		}
	}
	out.Diagnostics = convertMessages(res.Errors, res.Warnings)

	return out, nil
}

// mapFormat converts MinifyJS's format string into esbuild's enum.
// The empty string is handled by the caller (it means "preserve").
type formatMapping api.Format

func (f formatMapping) String() string {
	switch api.Format(f) {
	case api.FormatIIFE:
		return "iife"
	case api.FormatCommonJS:
		return "cjs"
	case api.FormatESModule:
		return "esm"
	default:
		return "default"
	}
}

func mapFormat(f string) formatMapping {
	switch f {
	case "iife":
		return formatMapping(api.FormatIIFE)
	case "cjs":
		return formatMapping(api.FormatCommonJS)
	case "esm":
		return formatMapping(api.FormatESModule)
	default:
		return formatMapping(api.FormatDefault)
	}
}

// mapSourcemap converts MinifyJS's sourcemap string into esbuild's
// enum. The empty string is handled by the caller.
func mapSourcemap(s string) api.SourceMap {
	switch s {
	case "inline":
		return api.SourceMapInline
	case "external":
		return api.SourceMapExternal
	case "both":
		return api.SourceMapInlineAndExternal
	default:
		return api.SourceMapNone
	}
}

func mapTarget(target string) api.Target {
	switch target {
	case "es5":
		return api.ES5
	case "es2015":
		return api.ES2015
	case "es2016":
		return api.ES2016
	case "es2017":
		return api.ES2017
	case "es2018":
		return api.ES2018
	case "es2019":
		return api.ES2019
	case "es2020":
		return api.ES2020
	case "es2021":
		return api.ES2021
	case "es2022":
		return api.ES2022
	case "es2023":
		return api.ES2023
	case "es2024":
		return api.ES2024
	case "es2025":
		return api.ES2025
	default:
		return api.ESNext
	}
}

// mapLegalComments converts MinifyJS's legal-comments mode into
// esbuild's enum. The empty string means "use esbuild's default",
// which is api.LegalCommentsEndOfFile for transform.
func mapLegalComments(s string) api.LegalComments {
	switch s {
	case "none":
		return api.LegalCommentsNone
	case "inline":
		return api.LegalCommentsInline
	case "eof":
		return api.LegalCommentsEndOfFile
	case "external":
		return api.LegalCommentsLinked
	default:
		return api.LegalCommentsDefault
	}
}

// convertMessages turns esbuild's separate Errors and Warnings slices
// into a single ordered diagnostic list. Errors come first so callers
// can short-circuit on the first one without scanning.
func convertMessages(errs, warns []api.Message) []diagnostics.Diagnostic {
	out := make([]diagnostics.Diagnostic, 0, len(errs)+len(warns))
	for _, m := range errs {
		out = append(out, convertMessage(m, diagnostics.SeverityError))
	}
	for _, m := range warns {
		out = append(out, convertMessage(m, diagnostics.SeverityWarning))
	}
	return out
}

// convertMessage converts one esbuild message into a MinifyJS
// diagnostic. esbuild reports location as a struct with optional
// fields; we flatten it to file+line+col, leaving them zero when
// esbuild did not provide a location.
func convertMessage(m api.Message, sev diagnostics.Severity) diagnostics.Diagnostic {
	d := diagnostics.Diagnostic{
		Severity: sev,
		Message:  m.Text,
		Code:     m.ID,
	}
	if m.Location != nil {
		d.File = m.Location.File
		d.Line = m.Location.Line
		d.Col = m.Location.Column
	}
	return d
}