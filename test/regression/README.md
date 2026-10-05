# Regression tests

One subdirectory per fixed GitHub issue, named `issue_NNNN/` where
`NNNN` is the issue number (zero-padded to 4 digits, e.g.
`issue_0001/`).

Each subdirectory contains:

- `input.js` — the smallest input that reproduces the original bug.
- `expected.min.js` — what MinifyJS should produce now that the bug
  is fixed.
- `README.md` — a link to the issue and a one-paragraph description
  of what went wrong.

The point of this directory is to make every fixed bug permanent.
A contributor who fixes an issue adds a subdirectory here. CI runs
every regression test on every pull request. If a future change
reintroduces the bug, the test fails with a clear link back to the
original issue.

A regression test is never deleted. If the expected output needs to
change (because esbuild changed its formatting, or because a new
minification pass changed the output), the `expected.min.js` is
regenerated and the `README.md` documents the reason.