# fx Documentation

> Source: https://uber-go.github.io/fx/
> Fetched: 2026-01-31T10:55:43.225207+00:00
> Content-Hash: f7166ae14004b69d
> Type: html

---

[ ](https://github.com/uber-go/fx/edit/master/docs/src/index.md "Edit this page")

# Fx¶

Fx is **a dependency injection system for Go**.

  * **Eliminate globals**

* * *

By using Fx-managed singletons, you can eliminate global state from your application. With Fx, you don't have to rely on `init()` functions for setup, instead relying on Fx to manage the lifecycle of your application.

  * **Reduce boilerplate**

* * *

Fx reduces the amount of code copy-pasted across your services. It lets you define shared application setup in a single place, and then reuse it across all your services.

  * **Automatic plumbing**

* * *

Fx automatically constructs your application's dependency graph. A component added to the application can be used by any other component without any additional configuration.

[Learn more about the dependency container ](container.html)

  * **Code reuse**

* * *

Fx lets teams within your organization build loosely-coupled and well-integrated shareable components referred to as modules.

[Learn more about modules ](modules.html)

  * **Battle-tested**

Fx is the backbone of nearly all Go services at Uber.




[Get started ](get-started/index.html)

May 13, 2025 May 13, 2025
