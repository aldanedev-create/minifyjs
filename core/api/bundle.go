package api

import "github.com/minifyjs/minifyjs/core/internal/engine"

// BundleOptions mirrors engine.BuildOptions.
type BundleOptions struct {
	EntryPoints   []string
	OutDir        string
	OutFile       string
	AbsWorkingDir string
	Platform      string
	Bundle        bool
	Splitting     bool
	Options
}

// Bundle runs a bundling build. The output is written to disk; the
// returned Result carries diagnostics and, when OutFile is used, the
// output code.
func Bundle(opts BundleOptions) (Result, error) {
	r, err := engine.Build(engine.BuildOptions{
		EntryPoints:   opts.EntryPoints,
		OutDir:        opts.OutDir,
		OutFile:       opts.OutFile,
		AbsWorkingDir: opts.AbsWorkingDir,
		Platform:      opts.Platform,
		Bundle:        opts.Bundle,
		Splitting:     opts.Splitting,
		Format:        opts.Format,
		Options: engine.Options{
			MinifyWhitespace:  opts.MinifyWhitespace,
			MinifyIdentifiers: opts.MinifyIdentifiers,
			MinifySyntax:      opts.MinifySyntax,
			Target:            opts.Target,
			Sourcemap:         opts.Sourcemap,
			Banner:            opts.Banner,
			Footer:            opts.Footer,
			LegalComments:     opts.LegalComments,
			Define:            opts.Define,
			Drop:              opts.Drop,
			Pure:              opts.Pure,
		},
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