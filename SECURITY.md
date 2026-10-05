# Security Policy

## Supported versions

Only the latest release of MinifyJS receives security updates. If
you are on an older version, please upgrade before reporting a
vulnerability.

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| other   | :x:                |

## Reporting a vulnerability

**Do not open a public GitHub issue for security problems.**

Report vulnerabilities privately by emailing
security@minifyjs.example. Include:

- A description of the issue and its impact.
- The version of MinifyJS affected (`minifyjs --version`).
- The operating system and architecture.
- A minimal reproduction if one is possible.
- Any suggested mitigation or fix, if you have one.

We aim to acknowledge reports within three business days. Once a
report is acknowledged, we will:

1. Confirm the issue and determine its severity.
2. Develop a fix in a private branch.
3. Coordinate a release with the reporter.
4. Credit the reporter in the release notes, unless they ask to
   remain anonymous.

We will not take legal action against researchers who follow this
policy and act in good faith.

## Scope

MinifyJS is a wrapper around [esbuild](https://esbuild.github.io/).
Security issues in esbuild itself should be reported to the esbuild
project at https://github.com/evanw/esbuild/security.

Issues that are in scope for MinifyJS include:

- The Python package executing unexpected code at install time.
- The Python bridge passing untrusted input to the native binary in
  a way that allows command injection or argument injection.
- The native binary writing to unintended paths.
- The cache reading or writing files outside its configured
  directory.
- Any behavior that would let a `.js` file being minified execute
  code during minification.

## Out of scope

- Denial of service through extremely large inputs. MinifyJS is a
  local build tool; if you feed it a 10 GB file, it will use a lot
  of memory. That is expected.
- Bugs in JavaScript engines that execute the *output* of MinifyJS.
  If the minified output has a bug in it, report that as a normal
  bug, not a vulnerability.
- Issues that require the attacker to already have write access to
  the user's filesystem.

## Safe harbor

We consider security research conducted in good faith under this
policy to be authorized. We will not pursue legal action against
researchers who follow this policy, and we will work with them to
understand and fix any issues they report.

## PGP key

If you need to encrypt your report, use the key published at
https://minifyjs.example/.well-known/security.asc.