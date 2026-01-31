# Black Formatter

> Source: https://black.readthedocs.io/en/stable/
> Fetched: 2026-01-31T16:07:53.047579+00:00
> Content-Hash: d2684f3677ee744b
> Type: html

---

# The uncompromising code formatter¶

> “Any color you like.”

By using _Black_ , you agree to cede control over minutiae of hand-formatting. In return, _Black_ gives you speed, determinism, and freedom from `pycodestyle` nagging about formatting. You will save time and mental energy for more important matters.

_Black_ makes code review faster by producing the smallest diffs possible. Blackened code looks the same regardless of the project you’re reading. Formatting becomes transparent after a while and you can focus on the content instead.

Try it out now using the [Black Playground](https://black.vercel.app).

Note - Black is now stable!

_Black_ is [successfully used](https://github.com/psf/black#used-by) by many projects, small and big. _Black_ has a comprehensive test suite, with efficient parallel tests, our own auto formatting and parallel Continuous Integration runner. Now that we have become stable, you should not expect large changes to formatting in the future. Stylistic changes will mostly be responses to bug reports and support for new Python syntax.

Also, as a safety measure which slows down processing, _Black_ will check that the reformatted code still produces a valid AST that is effectively equivalent to the original (see the [Pragmatism](the_black_code_style/current_style.html#pragmatism) section for details). If you’re feeling confident, use `\--fast`.

Note

[Black is licensed under the MIT license](license.html).

## Testimonials¶

**Mike Bayer** , author of [SQLAlchemy](https://www.sqlalchemy.org/):

> _I can’t think of any single tool in my entire programming career that has given me a bigger productivity increase by its introduction. I can now do refactorings in about 1% of the keystrokes that it would have taken me previously when we had no way for code to format itself._

**Dusty Phillips** , [writer](https://smile.amazon.com/s/ref=nb_sb_noss?url=search-alias%3Daps&amp;field-keywords=dusty+phillips):

> _Black is opinionated so you don’t have to be._

**Hynek Schlawack** , creator of [attrs](https://www.attrs.org/), core developer of Twisted and CPython:

> _An auto-formatter that doesn’t suck is all I want for Xmas!_

**Carl Meyer** , [Django](https://www.djangoproject.com/) core developer:

> _At least the name is good._

**Kenneth Reitz** , creator of [requests](http://python-requests.org/) and [pipenv](https://docs.pipenv.org/):

> _This vastly improves the formatting of our code. Thanks a ton!_

## Show your style¶

Use the badge in your project’s README.md:
    
    
    [![Code style: black](https://img.shields.io/badge/code%20style-black-000000.svg)](https://github.com/psf/black)
    

Using the badge in README.rst:
    
    
    .. image:: https://img.shields.io/badge/code%20style-black-000000.svg
       :target: https://github.com/psf/black
    

Looks like this:

[ ](https://github.com/psf/black)

## Contents¶

  * [The Black Code Style](the_black_code_style/index.html)
    * [Current style](the_black_code_style/current_style.html)
      * [Code style](the_black_code_style/current_style.html#code-style)
      * [Pragmatism](the_black_code_style/current_style.html#pragmatism)
    * [Future style](the_black_code_style/future_style.html)
      * [Preview style](the_black_code_style/future_style.html#preview-style)
      * [Unstable style](the_black_code_style/future_style.html#unstable-style)
    * [Stability Policy](the_black_code_style/index.html#stability-policy)



User Guide

  * [Getting Started](getting_started.html)
    * [Do you like the _Black_ code style?](getting_started.html#do-you-like-the-black-code-style)
    * [Try it out online](getting_started.html#try-it-out-online)
    * [Installation](getting_started.html#installation)
    * [Basic usage](getting_started.html#basic-usage)
    * [Next steps](getting_started.html#next-steps)
  * [Usage and Configuration](usage_and_configuration/index.html)
    * [The basics](usage_and_configuration/the_basics.html)
      * [Usage](usage_and_configuration/the_basics.html#usage)
      * [Configuration via a file](usage_and_configuration/the_basics.html#configuration-via-a-file)
      * [Next steps](usage_and_configuration/the_basics.html#next-steps)
    * [File collection and discovery](usage_and_configuration/file_collection_and_discovery.html)
      * [Ignoring unmodified files](usage_and_configuration/file_collection_and_discovery.html#ignoring-unmodified-files)
      * [.gitignore](usage_and_configuration/file_collection_and_discovery.html#gitignore)
    * [Black as a server (blackd)](usage_and_configuration/black_as_a_server.html)
      * [Usage](usage_and_configuration/black_as_a_server.html#usage)
      * [Protocol](usage_and_configuration/black_as_a_server.html#protocol)
    * [Black Docker image](usage_and_configuration/black_docker_image.html)
      * [Usage](usage_and_configuration/black_docker_image.html#usage)
  * [Integrations](integrations/index.html)
    * [Editor integration](integrations/editors.html)
      * [Emacs](integrations/editors.html#emacs)
      * [PyCharm/IntelliJ IDEA](integrations/editors.html#pycharm-intellij-idea)
      * [Wing IDE](integrations/editors.html#wing-ide)
      * [Vim](integrations/editors.html#vim)
      * [Gedit](integrations/editors.html#gedit)
      * [Visual Studio Code](integrations/editors.html#visual-studio-code)
      * [SublimeText](integrations/editors.html#sublimetext)
      * [Python LSP Server](integrations/editors.html#python-lsp-server)
      * [Atom/Nuclide](integrations/editors.html#atom-nuclide)
      * [Gradle (the build tool)](integrations/editors.html#gradle-the-build-tool)
      * [Kakoune](integrations/editors.html#kakoune)
      * [Thonny](integrations/editors.html#thonny)
    * [GitHub Actions integration](integrations/github_actions.html)
      * [Compatibility](integrations/github_actions.html#compatibility)
      * [Usage](integrations/github_actions.html#usage)
    * [Version control integration](integrations/source_version_control.html)
      * [Jupyter Notebooks](integrations/source_version_control.html#jupyter-notebooks)
      * [Excluding files with pre-commit](integrations/source_version_control.html#excluding-files-with-pre-commit)
  * [Guides](guides/index.html)
    * [Introducing _Black_ to your project](guides/introducing_black_to_your_project.html)
      * [Avoiding ruining git blame](guides/introducing_black_to_your_project.html#avoiding-ruining-git-blame)
    * [Using _Black_ with other tools](guides/using_black_with_other_tools.html)
      * [Black compatible configurations](guides/using_black_with_other_tools.html#black-compatible-configurations)
  * [Frequently Asked Questions](faq.html)
    * [Why spaces? I prefer tabs](faq.html#why-spaces-i-prefer-tabs)
    * [Does Black have an API?](faq.html#does-black-have-an-api)
    * [Is Black safe to use?](faq.html#is-black-safe-to-use)
    * [How stable is Black’s style?](faq.html#how-stable-is-black-s-style)
    * [Why is my file not formatted?](faq.html#why-is-my-file-not-formatted)
    * [Why is my Jupyter Notebook cell not formatted?](faq.html#why-is-my-jupyter-notebook-cell-not-formatted)
    * [Why does Flake8 report warnings?](faq.html#why-does-flake8-report-warnings)
    * [Which Python versions does Black support?](faq.html#which-python-versions-does-black-support)
    * [Why does my linter or typechecker complain after I format my code?](faq.html#why-does-my-linter-or-typechecker-complain-after-i-format-my-code)
    * [Can I run Black with PyPy?](faq.html#can-i-run-black-with-pypy)
    * [Why does Black not detect syntax errors in my code?](faq.html#why-does-black-not-detect-syntax-errors-in-my-code)
    * [What is `compiled: yes/no` all about in the version output?](faq.html#what-is-compiled-yes-no-all-about-in-the-version-output)
    * [Why are emoji not displaying correctly on Windows?](faq.html#why-are-emoji-not-displaying-correctly-on-windows)



Development

  * [Contributing](contributing/index.html)
    * [The basics](contributing/the_basics.html)
    * [Gauging changes](contributing/gauging_changes.html)
    * [Issue triage](contributing/issue_triage.html)
    * [Release process](contributing/release_process.html)
  * [Change Log](change_log.html)
    * [26.1.0](change_log.html#id1)
    * [25.12.0](change_log.html#id2)
    * [25.11.0](change_log.html#id5)
    * [25.9.0](change_log.html#id12)
    * [25.1.0](change_log.html#id17)
    * [24.10.0](change_log.html#id24)
    * [24.8.0](change_log.html#id30)
    * [24.4.2](change_log.html#id36)
    * [24.4.1](change_log.html#id39)
    * [24.4.0](change_log.html#id44)
    * [24.3.0](change_log.html#id48)
    * [24.2.0](change_log.html#id53)
    * [24.1.1](change_log.html#id59)
    * [24.1.0](change_log.html#id62)
    * [23.12.1](change_log.html#id68)
    * [23.12.0](change_log.html#id70)
    * [23.11.0](change_log.html#id77)
    * [23.10.1](change_log.html#id84)
    * [23.10.0](change_log.html#id90)
    * [23.9.1](change_log.html#id97)
    * [23.9.0](change_log.html#id100)
    * [23.7.0](change_log.html#id106)
    * [23.3.0](change_log.html#id118)
    * [23.1.0](change_log.html#id125)
    * [22.12.0](change_log.html#id134)
    * [22.10.0](change_log.html#id139)
    * [22.8.0](change_log.html#id147)
    * [22.6.0](change_log.html#id158)
    * [22.3.0](change_log.html#id164)
    * [22.1.0](change_log.html#id172)
    * [21.12b0](change_log.html#b0)
    * [21.11b1](change_log.html#b1)
    * [21.11b0](change_log.html#id184)
    * [21.10b0](change_log.html#id187)
    * [21.9b0](change_log.html#id191)
    * [21.8b0](change_log.html#id193)
    * [21.7b0](change_log.html#id197)
    * [21.6b0](change_log.html#id200)
    * [21.5b2](change_log.html#b2)
    * [21.5b1](change_log.html#id207)
    * [21.5b0](change_log.html#id210)
    * [21.4b2](change_log.html#id213)
    * [21.4b1](change_log.html#id216)
    * [21.4b0](change_log.html#id219)
    * [20.8b1](change_log.html#id222)
    * [20.8b0](change_log.html#id224)
    * [19.10b0](change_log.html#id227)
    * [19.3b0](change_log.html#id228)
    * [18.9b0](change_log.html#id229)
    * [18.6b4](change_log.html#b4)
    * [18.6b3](change_log.html#b3)
    * [18.6b2](change_log.html#id230)
    * [18.6b1](change_log.html#id231)
    * [18.6b0](change_log.html#id232)
    * [18.5b1](change_log.html#id233)
    * [18.5b0](change_log.html#id234)
    * [18.4a4](change_log.html#a4)
    * [18.4a3](change_log.html#a3)
    * [18.4a2](change_log.html#a2)
    * [18.4a1](change_log.html#a1)
    * [18.4a0](change_log.html#a0)
    * [18.3a4](change_log.html#id235)
    * [18.3a3](change_log.html#id236)
    * [18.3a2](change_log.html#id237)
    * [18.3a1](change_log.html#id238)
    * [18.3a0](change_log.html#id239)
  * [Authors](authors.html)



# Indices and tables¶

  * [Index](genindex.html)

  * [Search Page](search.html)



  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
