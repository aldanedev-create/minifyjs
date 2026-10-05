# Security

How MinifyJS handles security, and how to report a vulnerability.

## Reporting a vulnerability

Do **not** open a public GitHub issue for security problems. See
[SECURITY.md](../SECURITY.md) at the repository root for the full
policy, including the PGP key and the response timeline.

## What MinifyJS does

MinifyJS reads a JavaScript file, runs it through esbuild, and
writes a minified version. It does not execute the JavaScript it
processes. It does not open network connections. It does not read
from any file other than the input and the config, unless bundling
is enabled.

The binary is compiled with `CGO_ENABLED=0`, so it does not link
against any system C library at runtime. Every dependency is a Go
module pinned in `core/go.mod`.

## What MinifyJS does not do

- **It does not sandbox the input.** A minifier is not an
  interpreter. If you feed it a 10 GB file, it will use 10 GB of
  memory. That is expected.
- **It does not run the code it minifies.** A malicious `.js` file
  is just bytes to MinifyJS. It cannot execute them.
- **It does not validate that the input is JavaScript before
  parsing.** It parses as esbuild parses. Malformed input is
  reported as a syntax error, not as a security incident.
- **It does not evaluate expressions at build time.** If your
  input contains `while (true) {}`, MinifyJS will not hang on it.
  If it contains `eval("...")`, MinifyJS will not run the eval.

## The threat model

MinifyJS is a build tool. Its threat model assumes:

- The **input files are trusted** to the same degree as any other
  file in the user's working tree.
- The **config file** is trusted. It is looked up by walking up
  from the current directory, so a malicious config file in a
  parent directory could change behavior. Do not run `minifyjs`
  from inside an untrusted directory tree.
- The **cache directory** is trusted. It is in the user's own
  cache directory. A user with write access to that directory can
  already run arbitrary code as the user.
- The **binary** is trusted. It is installed via `pip install`,
  which does not authenticate the source beyond HTTPS.

MinifyJS does not assume:

- That the input is well-formed.
- That the input is small.
- That the input is free of adversarial content.

The first two are handled by the error-reporting path. The third is
out of scope: MinifyJS parses and rewrites the input, it does not
execute it.

## What is in scope for a vulnerability report

- The Python package executing code at install time (it should not;
  the wheel contains a binary and Python source, nothing else).
- The Python bridge passing untrusted input to the native binary in
  a way that allows command injection or argument injection.
- The native binary writing to unintended paths.
- The cache reading or writing files outside its configured
  directory.
- Any behavior that would let a `.js` file being minified execute
  code during minification.

## What is out of scope

- **Denial of service through large inputs.** MinifyJS is a local
  build tool. If you feed it a huge file, it uses a lot of memory.
- **Bugs in the output that run in a browser.** If the minified
  output has a bug, report it as a normal bug, not a vulnerability.
  (The engine is esbuild's, so an actual esbuild bug should go to
  the esbuild project.)
- **Issues that require the attacker to already have write access
  to the user's filesystem.**
- **Issues in dependencies that are already fixed upstream.**
  Upgrade and rerun.

## Dependency policy

MinifyJS has one Go dependency (`github.com/evanw/esbuild`) and zero
Python runtime dependencies. The benchmark and test tooling has
Python dependencies (pytest, ruff, mypy, psutil) but those are not
installed by `pip install minifyjs`.

The one Go dependency is:

- Checked against `govulncheck` on every CI run.
- Pinned to a specific version in `core/go.mod`.
- Updated by Dependabot, with a manual review before merge.

## Signed releases

Wheel signatures via sigstore are planned for a future release. As
of 0.1.0 they are not enabled; the checksums on the GitHub release
page are the only integrity guarantee.

## Acknowledgments

Security researchers who report issues responsibly are credited in
the release notes, unless they ask to remain anonymous. See
[SECURITY.md](../SECURITY.md) for the disclosure process.