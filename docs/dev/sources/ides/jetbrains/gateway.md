# JetBrains Gateway

> Source: https://www.jetbrains.com/help/idea/remote-development-overview.html
> Fetched: 2026-01-31T11:04:54.379843+00:00
> Content-Hash: f63f595ab8410b81
> Type: html

---

# Remote development overview

Remote development lets you use a remote machine, development container, WSL, or various providers to check out and load your project, index, analyze, build, run, debug, and test your code.

With Remote development you can do the following:

  * Edit, build, or debug applications on a different OS than you are running locally.

  * Use larger or more powerful hardware than your local machine for development.

  * Use a laptop as a thin client, no source code needs to be hosted on your local machine.

  * Work from anywhere, while the sensitive intelligence is hosted on the companyâs servers only.




## Connections

The remote host is a physical or virtual machine hosting the source code and running a backend IDE. You connect to the backend that transparently provides full access to all IDE features.

A connection to the remote machine can be established using various scenarios:

### SSH connection

SSH connection from a local machine into a remote server or vice versa (from an already installed IDE on your server to your local machine).

You can use one of the following ways:

  * JetBrains Toolbox App: supports the connection on Linux, macOS, and Windows. For more information, refer to the [Toolbox App](https://www.jetbrains.com/help/toolbox-app/about-the-instance.html) page.

  * IntelliJ IDEA: connects to your remote project from the IntelliJ IDEA welcome screen. For more information, refer to [Connect to a remote server from IntelliJ IDEA](remote-development-starting-page.html).

  * JetBrains Gateway: you can use JetBrains Gateway for the SSH connection to a Linux machine. You can also connect to various development environments. 

For more information, refer to [Connect and work with JetBrains Gateway](remote-development-a.html).




### Dev Container connection

Development Container connection when starting a Dev Container on the remote machine for the project with the JSON file located in the remote file system or for the project cloned from a Git repository.

Refer to [Start Dev Container for a remote project](start-dev-container-for-a-remote-project.html) for this workflow description.

### WSL connection

WSL connection when configuring your IDE backend to launch directly in WSL2. JetBrains Gateway offers native WSL support for such a scenario.

Refer to [Connect to a project running on WSL2](remote-development-a.html#run_in_wsl) for more information.

### Dev Environments

Connections to various development environments running on JetBrains CodeCanvas, Gitpod, Google Cloud, GitHub Codespaces, Amazon CodeCatalyst, and Coder are also available through JetBrains Gateway.

For more information on how to connect to each of the environments, refer to [Connect and work with JetBrains Gateway](remote-development-a.html).

## Extensibility

### The IDE backend

The backend can be extended with all the diversity of IntelliJ IDEA plugins in the following ways:

  * By unpacking required plugins into the [appropriate directories](https://intellij-support.jetbrains.com/hc/en-us/articles/206544519-Directories-used-by-the-IDE-to-store-settings-caches-plugins-and-logs?page=3)

  * By running the following code (requires network connection to [JetBrains Marketplace](https://plugins.jetbrains.com/)): 

remote-dev-server installPlugins <PLUGIN_ID1> <PLUGIN_ID2> ... 

Check the following example:

remote-dev-server installPlugins IdeaVIM




If a plugin provides a new set of inspections and features, all of those will be shown on JetBrains Client.

For more information, refer to [Install plugins](work-inside-remote-project.html#plugins).

### JetBrains Gateway SDK

JetBrains Gateway can be extended like any other IntelliJ platform-based product.

You can use one of the following ways:

  * Set up a new project with <https://github.com/JetBrains/gradle-intellij-plugin/> (`gradle-intellij-plugin` should be >= 1.1.4)

  * Use the following settings to build your plugin against JetBrains Gateway:

intellij { version.set("213.2667-CUSTOM-SNAPSHOT") type.set("GW") instrumentCode.set(false) } 

You may see available versions at <https://www.jetbrains.com/intellij-repository/snapshots> (see group com.jetbrains.gateway)




### Orchestration

Apart from the basic SSH and Code With Me connections, a vendor can customize JetBrains Gateway for its own orchestration service. This can be done within the custom deal between the JetBrains and the vendor.

JetBrains Gateway is based on the IntelliJ platform, it has APIs for connections and interactions with JetBrains Client. 

Check the following example:

A big organization wants to write its own orchestration. The basic SSH flow is not enough due to security reasons. The organization writes an internal plugin and delivers it to its developers. Developers can install this plugin in JetBrains Gateway or in IntelliJ IDEA on their laptops.

This is a very brief introduction of the APIs (they are not yet final, and indeed, this is not the full scope, but they explain the overall idea).

### JetBrains Client

JetBrains Client is not designed to be extensible for the connection part. However, you can develop and install all the variety of IntelliJ IDEA plugins, which modify the UI, keyboard shortcuts, themes, and other parts that touch the IDE UI interaction, but not its functionality.

17 October 2025

[Remote development](remote.html)[System requirements for remote development](prerequisites.html)
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
