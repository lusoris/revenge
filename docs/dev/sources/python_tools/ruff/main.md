# Ruff Documentation

> Source: https://docs.astral.sh/ruff/
> Fetched: 2026-01-31T16:07:27.962742+00:00
> Content-Hash: 2c150a3712495081
> Type: html

---

# Ruff

[](https://github.com/astral-sh/ruff) [](https://pypi.python.org/pypi/ruff) [](https://github.com/astral-sh/ruff/blob/main/LICENSE) [](https://pypi.python.org/pypi/ruff) [](https://github.com/astral-sh/ruff/actions) [](https://discord.com/invite/astral-sh)

[**Docs**](./) | [**Playground**](https://play.ruff.rs/)

An extremely fast Python linter and code formatter, written in Rust.

_Linting the CPython codebase from scratch._

  * ⚡️ 10-100x faster than existing linters (like Flake8) and formatters (like Black)
  * 🐍 Installable via `pip`
  * 🛠️ `pyproject.toml` support
  * 🤝 Python 3.14 compatibility
  * ⚖️ Drop-in parity with [Flake8](faq/#how-does-ruffs-linter-compare-to-flake8), isort, and [Black](faq/#how-does-ruffs-formatter-compare-to-black)
  * 📦 Built-in caching, to avoid re-analyzing unchanged files
  * 🔧 Fix support, for automatic error correction (e.g., automatically remove unused imports)
  * 📏 Over [800 built-in rules](rules/), with native re-implementations of popular Flake8 plugins, like flake8-bugbear
  * ⌨️ First-party [editor integrations](editors/) for [VS Code](https://github.com/astral-sh/ruff-vscode) and [more](editors/setup/)
  * 🌎 Monorepo-friendly, with [hierarchical and cascading configuration](configuration/#config-file-discovery)



Ruff aims to be orders of magnitude faster than alternative tools while integrating more functionality behind a single, common interface.

Ruff can be used to replace [Flake8](https://pypi.org/project/flake8/) (plus dozens of plugins), [Black](https://github.com/psf/black), [isort](https://pypi.org/project/isort/), [pydocstyle](https://pypi.org/project/pydocstyle/), [pyupgrade](https://pypi.org/project/pyupgrade/), [autoflake](https://pypi.org/project/autoflake/), and more, all while executing tens or hundreds of times faster than any individual tool.

Ruff is extremely actively developed and used in major open-source projects like:

  * [Apache Airflow](https://github.com/apache/airflow)
  * [Apache Superset](https://github.com/apache/superset)
  * [FastAPI](https://github.com/tiangolo/fastapi)
  * [Hugging Face](https://github.com/huggingface/transformers)
  * [Pandas](https://github.com/pandas-dev/pandas)
  * [SciPy](https://github.com/scipy/scipy)



...and [many more](https://github.com/astral-sh/ruff#whos-using-ruff).

Ruff is backed by [Astral](https://astral.sh), the creators of [uv](https://github.com/astral-sh/uv) and [ty](https://github.com/astral-sh/ty).

Read the [launch post](https://astral.sh/blog/announcing-astral-the-company-behind-ruff), or the original [project announcement](https://notes.crmarsh.com/python-tooling-could-be-much-much-faster).

## Testimonials

[**Sebastián Ramírez**](https://twitter.com/tiangolo/status/1591912354882764802), creator of [FastAPI](https://github.com/tiangolo/fastapi):

> Ruff is so fast that sometimes I add an intentional bug in the code just to confirm it's actually running and checking the code.

[**Nick Schrock**](https://twitter.com/schrockn/status/1612615862904827904), founder of [Elementl](https://www.elementl.com/), co-creator of [GraphQL](https://graphql.org/):

> Why is Ruff a gamechanger? Primarily because it is nearly 1000x faster. Literally. Not a typo. On our largest module (dagster itself, 250k LOC) pylint takes about 2.5 minutes, parallelized across 4 cores on my M1. Running ruff against our _entire_ codebase takes .4 seconds.

[**Bryan Van de Ven**](https://github.com/bokeh/bokeh/pull/12605), co-creator of [Bokeh](https://github.com/bokeh/bokeh/), original author of [Conda](https://docs.conda.io/en/latest/):

> Ruff is ~150-200x faster than flake8 on my machine, scanning the whole repo takes ~0.2s instead of ~20s. This is an enormous quality of life improvement for local dev. It's fast enough that I added it as an actual commit hook, which is terrific.

[**Timothy Crosley**](https://twitter.com/timothycrosley/status/1606420868514877440), creator of [isort](https://github.com/PyCQA/isort):

> Just switched my first project to Ruff. Only one downside so far: it's so fast I couldn't believe it was working till I intentionally introduced some errors.

[**Tim Abbott**](https://github.com/zulip/zulip/pull/23431#issuecomment-1302557034), lead developer of [Zulip](https://github.com/zulip/zulip) (also [here](https://github.com/astral-sh/ruff/issues/465#issuecomment-1317400028)):

> This is just ridiculously fast... `ruff` is amazing.

Back to top 
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
