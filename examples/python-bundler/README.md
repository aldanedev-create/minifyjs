# Python: bundle a multi-file project

Uses `minifyjs.bundle()` to resolve imports, tree-shake unused
exports, and produce a single minified output file with a source map.

## Run

```console
pip install minifyjs
python main.py
```

## What it demonstrates

- `minifyjs.bundle()` for multi-file projects
- Passing `entry_points`, `outfile`, `format`, `target`, `sourcemap`
- Handling `Result.has_errors` and iterating `Result.diagnostics`
- That MinifyJS writes the output files itself in bundle mode (unlike
  `minify()`, which returns the code as a string)

## Output

After running, `dist/` will contain:

```
dist/
├── bundle.js
└── bundle.js.map
```

The bundle contains the code from `main.js`, `utils.js`, and
`math.js`, minified and with identifiers renamed. Unused exports
(like `formatDate` from `utils.js` and `subtract` and `divide` from
`math.js`) are removed by tree shaking.

## Editing the example

Try removing the `import { add }` line from `main.js` and rerunning.
`add` will be tree-shaken out of the bundle entirely, and the output
will be smaller.