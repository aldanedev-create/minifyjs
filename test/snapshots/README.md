# Golden snapshots

Frozen outputs of internal stages that don't have their own public
API. These catch unintended changes to the engine's behavior even
when the final output happens to be correct.

Four kinds of snapshot live here, mirroring the pipeline stages:

- `parser/` — the AST that MinifyJS's internal parser would produce
  for a given input. (Currently delegated to esbuild; these snapshot
  esbuild's normalized AST representation as MinifyJS sees it.)
- `ast/` — post-analysis AST after scope resolution and before
  optimization. Frozen so a change to the analysis passes is visible.
- `printer/` — the exact byte sequence the printer emits for a
  given AST. Frozen so a whitespace change is caught immediately.
- `optimizer/` — the AST after each optimization pass, in order.
  Frozen so a change to one pass does not silently affect another.

## Regenerating

Snapshots are regenerated with:

    UPDATE_GOLDEN=1 go test ./...

Every regenerated snapshot appears as a diff in the pull request.
Reviewers approve the diff explicitly, which is the point: the
snapshot system makes every output change visible and reviewed.

Snapshots are not the primary correctness mechanism. They are a
change-detection mechanism. Real correctness is verified by the
integration tests, which execute the original and minified code and
compare behavior.