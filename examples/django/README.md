# Django: integrate MinifyJS with staticfiles

A minimal Django project that adds a `minifystatic` management
command. Running it minifies every `.js` file in `STATICFILES_DIRS`
into `STATIC_ROOT` with a `.min.js` suffix.

## Run

```console
pip install -r requirements.txt
python manage.py minifystatic
```

The command prints how many files it minified and how many bytes it
saved.

## What it demonstrates

- A `BaseCommand` subclass that invokes MinifyJS
- An `Adapter` subclass that mirrors Django's staticfiles layout
  (files come from `STATICFILES_DIRS`, output goes to `STATIC_ROOT`)
- Passing a `--no-mangle` flag through to MinifyJS options

## Fitting it into a real Django project

In a real project you would:

1. Copy `DjangoStaticAdapter` into your project (or your own app).
2. Copy `minifystatic.py` into `<yourapp>/management/commands/`.
3. Run `python manage.py minifystatic` before `collectstatic`.

The `minifystatic` step produces `.min.js` files. `collectstatic`
then copies both `app.js` and `app.min.js` into `STATIC_ROOT`. Your
templates reference `{% static 'app.min.js' %}`.

## Why not minify inside collectstatic

Django's `collectstatic` runs a pipeline of `STATICFILES_STORAGE`
backends. Hooking MinifyJS into that pipeline is possible but
surprising — collectstatic is not supposed to rewrite file contents.
Keeping minification as a separate command makes the dependency
explicit and lets CI run it once.

## The custom command

```python
class Command(BaseCommand):
    help = "Minify every JavaScript file in STATICFILES_DIRS into STATIC_ROOT"

    def handle(self, *args, **options):
        from minifyjs import Adapter, Options
        ...
```

Everything MinifyJS-specific is inside `handle()`. The rest of the
command is standard Django boilerplate.