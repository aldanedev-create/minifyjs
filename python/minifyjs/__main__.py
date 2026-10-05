"""Enables ``python -m minifyjs ...`` as an alternative to the
``minifyjs`` console script installed by pip."""

from .cli import main

raise SystemExit(main())
