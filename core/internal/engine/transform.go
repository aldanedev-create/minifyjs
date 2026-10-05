package engine

import (
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

	if opts.SourceName != "" {
	    tOpts.Sourcefile = opts.SourceName
	}

	tOpts := api.TransformOptions{
		Loader: api.LoaderJS,

		MinifyWhitespace:  opts.MinifyWhitespace,
		MinifyIdentifiers: opts.MinifyIdentifiers,
		MinifySyntax:      opts.MinifySyntax,

		// Never let esbuild inject its own runtime helpers for a
		// single-file transform; there is nowhere for them to live.
		// Callers who need helpers must use Build.
		LegalComments: mapLegalComments(opts.LegalComments),
	}

	if opts.Target != "" {
		tOpts.Target = api.Target(opts.Target)
	}
	if opts.Format != "" {
		tOpts.Format = mapFormat(opts.Format)
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

	res := api.Transform(source, tOpts)

	out := Result{
		Code:          string(res.Code),
		OriginalBytes: len(source),
		MinifiedBytes: len(res.Code),
	}
	if len(res.Map) > 0 {
		out.Map = string(res.Map)
	}
	out.Diagnostics = convertMessages(res.Errors, res.Warnings)

	return out, nil
}

// mapFormat converts MinifyJS's format string into esbuild's enum.
// The empty string is handled by the caller (it means "preserve").
func mapFormat(f string) api.Format {
	switch f {
	case "iife":
		return api.FormatIIFE
	case "cjs":
		return api.FormatCommonJS
	case "esm":
		return api.FormatESModule
	default:
		return api.FormatDefault
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
		return api.SourceMapBoth
	default:
		return api.SourceMapNone
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