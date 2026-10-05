package api

import "github.com/minifyjs/minifyjs/core/internal/engine"

// Minify transforms a single JavaScript source string. It is the
// entry point for the common case; bundle() is the entry point for
// multi-file projects.
func Minify(source string, opts Options) (Result, error) {
	r, err := engine.Transform(source, engine.Options{
		MinifyWhitespace:  opts.MinifyWhitespace,
		MinifyIdentifiers: opts.MinifyIdentifiers,
		MinifySyntax:      opts.MinifySyntax,
		Target:            opts.Target,
		Format:            opts.Format,
		Sourcemap:         opts.Sourcemap,
		Banner:            opts.Banner,
		Footer:            opts.Footer,
		LegalComments:     opts.LegalComments,
		SourceName:        opts.SourceName,
	})
	if err != nil {
		return Result{}, err
	}
	return Result{
		Code:          r.Code,
		Map:           r.Map,
		Diagnostics:   r.Diagnostics,
		OriginalBytes: r.OriginalBytes,
		MinifiedBytes: r.MinifiedBytes,
	}, nil
}