# Flask: minify static assets at startup

A minimal Flask app that minifies every `.js` file under `static/`
when the app starts, then serves the minified versions.

## Run

```console
pip install -r requirements.txt
python app.py
```

Then open <http://127.0.0.1:5000/>.

## What it demonstrates

- Subclassing `minifyjs.Adapter` for a Flask app
- Using `Adapter.minify_directory()` to process a folder
- Serving the minified files with Flask's `send_from_directory`

## Why minify at startup

For a small Flask app, running MinifyJS when the process starts is
the right balance: no build step, no stale artifacts, no separate
tool to run. For a larger app, you would run MinifyJS as part of the
deploy pipeline instead, and the adapter would be invoked from a
management command or a CI job.

## The `Adapter` subclass

```python
class FlaskStaticAdapter(Adapter):
    def output_name(self, path: Path) -> Path:
        return path.with_name(f"{path.stem}.min{path.suffix}")
```

The only override is `output_name`, which tells the adapter to write
`app.min.js` next to `app.js`. That is the entire integration.