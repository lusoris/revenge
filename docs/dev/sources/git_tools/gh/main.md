# GitHub CLI

> Source: https://cli.github.com/manual/
> Fetched: 2026-02-01T11:55:49.308994+00:00
> Content-Hash: 4bcdcf915e099295
> Type: html

---

# GitHub CLI manual

GitHub CLI, or `gh`, is a command-line interface to GitHub for use in your terminal or your scripts.

- [Available commands](./gh)

- [Usage examples](./examples)

- [Community extensions](https://github.com/topics/gh-extension)

## Installation

You can find installation instructions on our [README](https://github.com/cli/cli#installation).

## Configuration

- Run [`gh auth login`](./gh_auth_login) to authenticate with your GitHub account. Alternatively, `gh` will respect the `GITHUB_TOKEN` [environment variable](./gh_help_environment).

- To set your preferred editor, use `gh config set editor <editor>`. Read more about [`gh config`](./gh_config) and [environment variables](./gh_help_environment).

- Declare your aliases for often-used commands with [`gh alias set`](./gh_alias_set).

## GitHub Enterprise

GitHub CLI supports GitHub Enterprise Server 2.20 and above. To authenticate with a GitHub instance, run:

    gh auth login --hostname <hostname>
    

To define this host as a default for all GitHub CLI commands, set the GH_HOST environment variable:

    export GH_HOST=<hostname>
    

Finally, to authenticate commands in scripting mode or automation, set the GH_ENTERPRISE_TOKEN:

    export GH_ENTERPRISE_TOKEN=<access-token>
    

## Support

- Ask usage questions and send us feedback in [Discussions](https://github.com/cli/cli/discussions)

- Report bugs or search for existing feature requests in our [issue tracker](https://github.com/cli/cli/issues)

  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
