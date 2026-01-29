# fx Documentation

> Auto-fetched from [https://uber-go.github.io/fx/](https://uber-go.github.io/fx/)
> Last Updated: 2026-01-28T21:43:08.936937+00:00

---

Fx
uber-go/fx
Home
Get started
Get started
Create a minimal application
Add an HTTP server
Register a handler
Add a logger
Decouple registration
Register another handler
Register many handlers
Conclusion
Concepts
Concepts
Container
Lifecycle
Modules
Features
Features
Parameter Objects
Result Objects
Annotations
Value groups
Value groups
Feeding value groups
Consuming value groups
FAQ
Community
Community
Contributing
Release notes
API Reference
Fx
¶
Fx is
a dependency injection system for Go
.
Eliminate globals
By using Fx-managed singletons,
you can eliminate global state from your application.
With Fx, you don't have to rely on
init()
functions for setup,
instead relying on Fx to manage the lifecycle of your application.
Reduce boilerplate
Fx reduces the amount of code copy-pasted across your services.
It lets you define shared application setup in a single place,
and then reuse it across all your services.
Automatic plumbing
Fx automatically constructs your application's dependency graph.
A component added to the application can be used by any other component
without any additional configuration.
Learn more about the dependency container
Code reuse
Fx lets teams within your organization build loosely-coupled
and well-integrated shareable components referred to as modules.
Learn more about modules
Battle-tested
Fx is the backbone of nearly all Go services at Uber.
Get started
May 13, 2025
May 13, 2025