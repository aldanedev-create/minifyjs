package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

// Key computes a cache key from the input source, the options that
// affect output, and any per-run context (banner, footer).
//
// The key is a SHA-256 of a canonical serialization of those inputs.
// SHA-256 is used rather than a faster hash because a cache key
// collision would silently produce wrong output, which is the worst
// possible failure mode for a build tool.
//
// The serialization is deliberately simple: a fixed set of fields,
// each length-prefixed, joined by '\n'. Length prefixing prevents
// the classic ambiguity where ("ab", "c") and ("a", "bc") hash to
// the same value.
func Key(source string, opts KeyOptions) string {
	var b strings.Builder
	writeField(&b, "v", "1") // schema version; bump if the layout changes
	writeField(&b, "source", source)
	writeField(&b, "target", opts.Target)
	writeField(&b, "format", opts.Format)
	writeField(&b, "sourcemap", opts.Sourcemap)
	writeField(&b, "banner", opts.Banner)
	writeField(&b, "footer", opts.Footer)
	writeField(&b, "legalComments", opts.LegalComments)
	writeField(&b, "minifyWhitespace", boolStr(opts.MinifyWhitespace))
	writeField(&b, "minifyIdentifiers", boolStr(opts.MinifyIdentifiers))
	writeField(&b, "minifySyntax", boolStr(opts.MinifySyntax))
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

// KeyOptions is the subset of engine.Options that affects output.
// It is duplicated here (rather than importing engine) so the cache
// package has no dependency on esbuild.
type KeyOptions struct {
	MinifyWhitespace  bool
	MinifyIdentifiers bool
	MinifySyntax      bool
	Target            string
	Format            string
	Sourcemap         string
	Banner            string
	Footer            string
	LegalComments     string
}

func writeField(b *strings.Builder, name, value string) {
	b.WriteString(name)
	b.WriteByte('=')
	b.WriteString(strconv.Itoa(len(value)))
	b.WriteByte(':')
	b.WriteString(value)
	b.WriteByte('\n')
}

func boolStr(v bool) string {
	if v {
		return "1"
	}
	return "0"
}