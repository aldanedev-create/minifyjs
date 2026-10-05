# End-to-end tests

Cross-cutting scenarios that exercise MinifyJS as a whole. These are
not unit tests and not single-CLI-invocation tests (those live in
`core/integration/`). They simulate a real workflow:

- Build a small project with the Python API, serve it with a
  framework, fetch the served assets, and assert the served bytes
  match what MinifyJS produced.
- Run `minifyjs` from a fresh shell environment with no `PATH`
  tricks and no Node.js on the machine, and verify that a complete
  build works.
- Take a real-world JavaScript file (a library, a framework bundle)
  and verify that `minifyjs` succeeds, that the output is smaller,
  and that executing the output produces the same observable
  behavior as the input.

Each scenario lives in its own subdirectory with its own README,
inputs, and expected outputs. They are wired into the CI workflow
`integration-test.yml`.