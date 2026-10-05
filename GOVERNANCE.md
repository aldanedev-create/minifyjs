# Governance

This document describes how decisions are made in the MinifyJS
project.

## Project scope

MinifyJS is a small, focused tool. It is a native JavaScript
minifier and optimizer distributed as a Python wheel and a
standalone binary, with no Node.js dependency. It is a wrapper
around [esbuild](https://esbuild.github.io/)'s Go API, not a
from-scratch JavaScript toolchain.

Changes that are in scope:

- Bug fixes
- Documentation improvements
- New flags that map to existing esbuild capabilities
- Improvements to the Python package, CLI ergonomics, or error
  reporting
- Ports to new platforms
- Performance improvements to the wrapper layer

Changes that are **not** in scope:

- Replacing esbuild with a different engine
- Adding support for languages other than JavaScript
- Adding a plugin system
- Reimplementing features that esbuild already provides
- Vendoring dependencies

Proposals that fall outside scope are closed with a link to this
document. That is not a judgment on the proposal's quality; it is a
statement about what the project is.

## Roles

### Users

Anyone who uses MinifyJS. Users participate by filing issues,
answering questions in discussions, and reporting bugs.

### Contributors

Anyone who has had a pull request merged. Contributors review other
pull requests, triage issues, and participate in design
discussions.

### Maintainers

Contributors who have been given commit access. Maintainers merge
pull requests, cut releases, manage the security policy, and make
final decisions when consensus cannot be reached.

The current maintainers are listed in
[`.github/CODEOWNERS`](.github/CODEOWNERS) (once that file exists)
and on the project's website.

### Lead maintainer

One maintainer is the lead. The lead breaks ties, decides on scope
questions when the maintainer group is split, and is the final
authority on security disclosures.

## Becoming a maintainer

There is no formal process. A contributor becomes a maintainer when
the existing maintainers agree that they have demonstrated:

- Sustained, high-quality contributions over several months.
- Good judgment on scope and design questions.
- Constructive participation in code review.
- Adherence to the Code of Conduct.

Maintainers are added by consensus of the existing maintainers. Any
maintainer may step down at any time by opening an issue or sending
an email to the lead.

## Decision making

Most decisions are made by consensus in a pull request or issue
thread. When consensus cannot be reached:

1. Any participant may propose a decision by writing a short summary
   of the options and their trade-offs in the relevant thread.
2. Maintainers discuss for at least one week.
3. If no consensus emerges, the lead maintainer decides.

Decisions that affect the public API, the on-disk format of the
cache or config files, or the exit codes are documented in the
relevant file under `docs/` and, when they are user-visible, in
`CHANGELOG.md`.

## Releases

Releases are cut by a maintainer. The release process is documented
in [`docs/development/release-process.md`](docs/development/release-process.md).

## Changing this document

Changes to this document are proposed via pull request. They are
merged only when every active maintainer has approved or has been
silent for two weeks.