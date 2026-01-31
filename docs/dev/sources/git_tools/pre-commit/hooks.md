# pre-commit Hooks

> Source: https://pre-commit.com/hooks.html
> Fetched: 2026-01-31T16:08:24.167189+00:00
> Content-Hash: 72f44770b9ce478c
> Type: html

---

# Supported hooks

##  featured hooks ¶

here are a few hand-picked repositories which provide pre-commit integrations.

these are fairly popular and are generally known to work well in most setups!

_this list is not intended to be exhaustive_

provided by the pre-commit team:

  * [pre-commit/pre-commit-hooks](https://github.com/pre-commit/pre-commit-hooks): a handful of language-agnostic hooks which are universally useful!
  * [pre-commit/pygrep-hooks](https://github.com/pre-commit/pygrep-hooks): a few quick regex-based hooks for a handful of quick syntax checks
  * [pre-commit/sync-pre-commit-deps](https://github.com/pre-commit/sync-pre-commit-deps): sync pre-commit hook dependencies based on other installed hooks
  * [pre-commit/mirrors-*](https://github.com/orgs/pre-commit/repositories?language=&q=%22mirrors-%22+archived%3AFalse&sort=): pre-commit mirrors of a handful of popular tools



for python projects:

  * [asottile/pyupgrade](https://github.com/asottile/pyupgrade): automatically upgrade syntax for newer versions of the language
  * [asottile/(others)](https://sourcegraph.com/search?q=context:global+file:%5E%5C.pre-commit-hooks%5C.yaml%24+repo:%5Egithub.com/asottile/): a few other repos by the pre-commit creator
  * [psf/black](https://github.com/psf/black): The uncompromising Python code formatter
  * [hhatto/autopep8](https://github.com/hhatto/autopep8): automatically fixes PEP8 violations
  * [astral-sh/ruff-pre-commit](https://github.com/astral-sh/ruff-pre-commit): the ruff linter and formatter for python
  * [google/yapf](https://github.com/google/yapf): a highly configurable python formatter
  * [PyCQA/flake8](https://github.com/PyCQA/flake8): a linter framework for python
  * [PyCQA/isort](https://github.com/PyCQA/isort): an import sorter for python
  * [PyCQA/(others)](https://sourcegraph.com/search?q=context:global+file:%5E%5C.pre-commit-hooks%5C.yaml%24+repo:%5Egithub.com/PyCQA/): a few other python code quality tools
  * [adamchainz/django-upgrade](https://github.com/adamchainz/django-upgrade): automatically upgrade your Django project code



for shell scripts:

  * [shellcheck-py/shellcheck-py](https://github.com/shellcheck-py/shellcheck-py): runs shellcheck on your scripts
  * [openstack/bashate](https://github.com/openstack/bashate): code style enforcement for bash programs



for the web:

  * [biomejs/pre-commit](https://github.com/biomejs/pre-commit): a fast formatter / fixer written in rust
  * [standard/standard](https://github.com/standard/standard): linter / fixer
  * [oxipng/oxipng](https://github.com/oxipng/oxipng): optimize png files



for configuration files:

  * [python-jsonschema/check-jsonschema](https://github.com/python-jsonschema/check-jsonschema): check many common configurations with jsonschema
  * [rhysd/actionlint](https://github.com/rhysd/actionlint): lint your GitHub Actions workflow files
  * [google/yamlfmt](https://github.com/google/yamlfmt): a formatter for yaml files
  * [adrienverge/yamllint](https://github.com/adrienverge/yamllint): a linter for YAML files



for text / docs / prose:

  * [crate-ci/typos](https://github.com/crate-ci/typos): find and fix common typographical errors
  * [thlorenz/doctoc](https://github.com/thlorenz/doctoc): generate a table-of-contents in markdown files
  * [amperser/proselint](https://github.com/amperser/proselint): A linter for prose.
  * [markdownlint/markdownlint](https://github.com/markdownlint/markdownlint): a Markdown lint tool in Ruby
  * [DavidAnson/markdownlint-cli2](https://github.com/DavidAnson/markdownlint-cli2): a Markdown lint tool in Node
  * [codespell-project/codespell](https://github.com/codespell-project/codespell): check code for common misspellings



for linting commit messages:

  * [jorisroovers/gitlint](https://github.com/jorisroovers/gitlint)
  * [commitizen-tools/commitizen](https://github.com/commitizen-tools/commitizen)



for secret scanning / security:

  * [gitleaks/gitleaks](https://github.com/gitleaks/gitleaks)
  * [trufflesecurity/truffleHog](https://github.com/trufflesecurity/truffleHog)
  * [thoughtworks/talisman](https://github.com/thoughtworks/talisman)



for other programming languages:

  * [realm/SwiftLint](https://github.com/realm/SwiftLint): enforce Swift style and conventions
  * [nicklockwood/SwiftFormat](https://github.com/nicklockwood/SwiftFormat): a formatter for Swift
  * [AleksaC/terraform-py](https://github.com/AleksaC/terraform-py): format and validate terraform syntax
  * [rubocop/rubocop](https://github.com/rubocop/rubocop): static analysis and formatting for Ruby
  * [bufbuild/buf](https://github.com/bufbuild/buf): tooling for Protocol Buffers
  * [sqlfluff/sqlfluff](https://github.com/sqlfluff/sqlfluff): a modular linter and auto formatter for SQL
  * [aws-cloudformation/cfn-lint](https://github.com/aws-cloudformation/cfn-lint): aws CloudFormation linter
  * [google/go-jsonnet](https://github.com/google/go-jsonnet): linter / formatter for jsonnet
  * [JohnnyMorganz/StyLua](https://github.com/JohnnyMorganz/StyLua): an opinionated Lua code formatter
  * [Koihik/LuaFormatter](https://github.com/Koihik/LuaFormatter): a formatter for Lua code
  * [mrtazz/checkmake](https://github.com/mrtazz/checkmake): linter for Makefile syntax
  * [nbqa-dev/nbqa](https://github.com/nbQA-dev/nbQA): run common linters on Jupyter Notebooks



##  finding hooks ¶

it's recommended to use your favorite searching tool to find existing hooks to use in your project.

for example, here's some searches you may find useful using [sourcegraph](https://sourcegraph.com/search):

  * hooks which run on python files: [`file:^\\.pre-commit-hooks\\.yaml$ "types: [python]"`](https://sourcegraph.com/search?q=context:global+file:%5E.pre-commit-hooks.yaml%24+%22types:+%5Bpython%5D%22)
  * hooks which run on shell files: [`file:^\\.pre-commit-hooks\\.yaml$ "types: [shell]"`](https://sourcegraph.com/search?q=context:global+file:%5E.pre-commit-hooks.yaml%24+%22types:+%5Bshell%5D%22)
  * pre-commit configurations in popular projects: [`file:^\\.pre-commit-config\\.yaml$`](https://sourcegraph.com/search?q=context:global+file:%5E.pre-commit-hooks.yaml)



you may also find [github's search](https://github.com/search) useful as well, though its querying and sorting capabilities are quite limited plus it requires a login:

  * repositories providing hooks: [`path:.pre-commit-hooks.yaml language:YAML`](https://github.com/search?q=path%3A.pre-commit-hooks.yaml+language%3AYAML&type=code&l=YAML)



##  adding to this page ¶

the previous iteration of this page was a laundry list of hooks and maintaining quality of the listed tools was cumbersome.

**this page is not intended to be exhaustive**

you may send [a pull request](https://github.com/pre-commit/pre-commit.com/blob/main/sections/hooks.md) to expand this list however there are a few requirements you _must_ follow or your PR will be closed without comment:

  * the tool must already be fairly popular (>500 stars)
  * the tool must use a managed language (no `unsupported` / `unsupported_script` / `docker` hooks)
  * the tool must operate on files


  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
