package main

import (
	"fmt"
	"strings"
)

// parseFlags parses args into a flags struct.
//
// The parser is a hand-written state machine because:
//
//   - Go's flag package stops at the first positional argument,
//     which makes `minifyjs app.js -o out.js` behave differently
//     from `minifyjs -o out.js app.js`. Users expect both to work.
//   - We want to support both --flag and --flag=value.
//   - We want to support --no-<flag> negation uniformly.
//
// The grammar is: any number of flags and positionals, in any order,
// except that flags that take a value must be followed by that value
// (either as the next arg or as --flag=value).
func parseFlags(args []string) (*flags, error) {
	f := newFlags()
	i := 0
	for i < len(args) {
		arg := args[i]
		i++

		// Positional.
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			f.inputs = append(f.inputs, arg)
			continue
		}

		// Long form with = value.
		var name, inlineValue string
		hasInline := false
		if strings.HasPrefix(arg, "--") {
			if eq := strings.IndexByte(arg, '='); eq >= 0 {
				name = arg[:eq]
				inlineValue = arg[eq+1:]
				hasInline = true
			} else {
				name = arg
			}
		} else {
			// Short flag: only single-letter forms, no bundling.
			name = arg
		}

		// consumeValue returns the value for the current flag,
		// either the inline value or the next argument.
		consumeValue := func() (string, error) {
			if hasInline {
				return inlineValue, nil
			}
			if i >= len(args) {
				return "", fmt.Errorf("flag %s requires a value", name)
			}
			v := args[i]
			i++
			return v, nil
		}

		switch name {
		case "-h", "--help":
			f.showHelp = true
		case "-V", "--version":
			f.showVersion = true

		case "-o", "--output", "--outfile":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			f.output = v
			f.outputSet = true

		case "--compress":
			f.compress = true
			f.compressSet = true
		case "--no-compress":
			f.compress = false
			f.compressSet = true
		case "--mangle":
			f.mangle = true
			f.mangleSet = true
		case "--no-mangle":
			f.mangle = false
			f.mangleSet = true
		case "--no-minify":
			f.noMinify = true

		case "--target":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			f.target = v
			f.targetSet = true
		case "--format":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			f.format = v
			f.formatSet = true
		case "--sourcemap":
			// Accept both "--sourcemap" (implicit "external") and
			// "--sourcemap=inline".
			if hasInline {
				f.sourcemap = inlineValue
			} else {
				f.sourcemap = "external"
			}
			f.smSet = true

		case "--banner":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			f.banner = v
			f.bannerSet = true
		case "--footer":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			f.footer = v
			f.footerSet = true
		case "--legal-comments":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			f.legalComments = v
			f.lcSet = true

		case "--define":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			k, val, ok := strings.Cut(v, "=")
			if !ok {
				return nil, fmt.Errorf("--define requires KEY=VALUE, got %q", v)
			}
			f.defines[k] = val

		case "--drop":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			if v != "console" && v != "debugger" {
				return nil, fmt.Errorf("invalid --drop kind %q", v)
			}
			f.drops = append(f.drops, v)

		case "--pure":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			f.pures = append(f.pures, v)

		case "--bundle":
			f.bundle = true
			f.bundleSet = true
		case "--no-bundle":
			f.bundle = false
			f.bundleSet = true
		case "--outdir":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			f.outdir = v
			f.outdirSet = true
		case "--platform":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			f.platform = v
			f.platformSet = true
		case "--splitting":
			f.splitting = true
			f.splitSet = true
		case "--no-splitting":
			f.splitting = false
			f.splitSet = true

		case "--cache":
			f.cache = true
			f.cacheSet = true
		case "--no-cache":
			f.cache = false
			f.cacheSet = true
		case "--cache-dir":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			f.cacheDir = v
			f.cacheDirSet = true

		case "--config":
			v, err := consumeValue()
			if err != nil {
				return nil, err
			}
			f.configPath = v
			f.configPathSet = true
		case "--no-config":
			f.noConfig = true

		case "-q", "--quiet":
			f.quiet = true
			f.quietSet = true
		case "--no-quiet":
			f.quiet = false
			f.quietSet = true
		case "--verbose":
			f.verbose = true
			f.verboseSet = true
		case "--no-verbose":
			f.verbose = false
			f.verboseSet = true

		case "--watch":
			f.watch = true

		case "--":
			// Everything after -- is positional.
			for i < len(args) {
				f.inputs = append(f.inputs, args[i])
				i++
			}
			return f, nil

		default:
			return nil, fmt.Errorf("unknown flag %s", name)
		}
	}
	if f.outputSet && f.outdirSet {
		return nil, fmt.Errorf("--outfile and --outdir cannot be used together")
	}
	return f, nil
}