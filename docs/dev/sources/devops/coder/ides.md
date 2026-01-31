# Coder IDEs

> Source: https://coder.com/docs/ides
> Fetched: 2026-01-31T11:05:31.934747+00:00
> Content-Hash: 6ec3996787ad0df6
> Type: html

---

[Home](/docs "Home")[User Guides](/docs/user-guides "User Guides")Access Workspaces

There are many ways to connect to your workspace, the options are only limited by the template configuration.

Deployment operators can learn more about different types of workspace connections and performance in our [networking docs](/docs/admin/infrastructure).

You can see the primary methods of connecting to your workspace in the workspace dashboard.

## Web Terminal

The Web Terminal is a browser-based terminal that provides instant access to your workspace's shell environment. It uses [xterm.js](https://xtermjs.org/) and WebSocket technology for a responsive terminal experience with features like persistent sessions, Unicode support, and clickable URLs.

Read the complete [Web Terminal documentation](/docs/user-guides/workspace-access/web-terminal) for customization options, keyboard shortcuts, and troubleshooting guides.

## SSH

### Through with the CLI

Coder will use the optimal path for an SSH connection (determined by your deployment's [networking configuration](/docs/admin/infrastructure)) when using the CLI:

``coder ssh my-workspace ``

Or, you can configure plain SSH on your client below.

Note

The `coder ssh` command does not have full parity with the standard SSH command. For users who need the full functionality of SSH, use the configuration method below.

### Configure SSH

Coder generates [SSH key pairs](/docs/admin/security/secrets#ssh-keys) for each user to simplify the setup process.

  1. Use your terminal to authenticate the CLI with Coder web UI and your workspaces:

``coder login <accessURL> ``

  2. Access Coder via SSH:

``coder config-ssh ``

  3. Run `coder config-ssh --dry-run` if you'd like to see the changes that will be before you proceed:

``coder config-ssh --dry-run ``

  4. Confirm that you want to continue by typing **yes** and pressing enter. If




successful, you'll see the following message:

``You should now be able to ssh into your workspace. For example, try running: $ ssh coder.<workspaceName> ``

Your workspace is now accessible via `ssh coder.<workspace_name>` (for example, `ssh coder.myEnv` if your workspace is named `myEnv`).

## Visual Studio Code

You can develop in your Coder workspace remotely with [VS Code](https://code.visualstudio.com/download). We support connecting with the desktop client and VS Code in the browser with code-server.

Read more details on [using VS Code in your workspace](/docs/user-guides/workspace-access/vscode).

## Cursor

[Cursor](https://cursor.sh/) is an IDE built on VS Code with enhanced AI capabilities. Cursor connects using the Coder extension.

Read more about [using Cursor with your workspace](/docs/user-guides/workspace-access/cursor).

## Windsurf

[Windsurf](/docs/user-guides/workspace-access/windsurf) is Codeium's code editor designed for AI-assisted development. Windsurf connects using the Coder extension.

## JetBrains IDEs

We support JetBrains IDEs using [Gateway](https://www.jetbrains.com/remote-development/gateway/). The following IDEs are supported for remote development:

  * IntelliJ IDEA
  * CLion
  * GoLand
  * PyCharm
  * Rider
  * RubyMine
  * WebStorm
  * [JetBrains Fleet](/docs/user-guides/workspace-access/jetbrains/fleet)



Read our [docs on JetBrains](/docs/user-guides/workspace-access/jetbrains) for more information on connecting your JetBrains IDEs.

## code-server

[code-server](https://github.com/coder/code-server) is our supported method of running VS Code in the web browser. Learn more about [what makes code-server different from VS Code web](/docs/user-guides/workspace-access/code-server) or visit the [documentation for code-server](https://coder.com/docs/code-server/latest).

## Other Web IDEs

We support a variety of other browser IDEs and tools to interact with your workspace. Each of these can be configured by your template admin using our [Web IDE guides](/docs/admin/templates/extending-templates/web-ides).

Supported IDEs:

  * VS Code Web
  * JupyterLab
  * RStudio
  * Airflow
  * File Browser



Our [Module Registry](https://registry.coder.com/modules) also hosts a variety of tools for extending the capability of your workspace. If you have a request for a new IDE or tool, please file an issue in our [Modules repo](https://github.com/coder/registry/issues).

## Ports and Port forwarding

You can manage listening ports on your workspace page through with the listening ports window in the dashboard. These ports are often used to run internal services or preview environments.

You can also [share ports](/docs/user-guides/workspace-access/port-forwarding#sharing-ports) with other users, or [port-forward](/docs/user-guides/workspace-access/port-forwarding#the-coder-port-forward-command) through the CLI with `coder port forward`. Read more in the [docs on workspace ports](/docs/user-guides/workspace-access/port-forwarding).

## Remote Desktops

Coder also supports connecting with an RDP solution, see our [RDP guide](/docs/user-guides/workspace-access/remote-desktops) for details.
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
