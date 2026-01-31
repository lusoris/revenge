# Kubebuilder Book

> Source: https://book.kubebuilder.io/
> Fetched: 2026-01-30T23:56:24.679469+00:00
> Content-Hash: b90bdd3fe6b4eae3
> Type: html

---

Tip

Impatient readers may head straight to

Quick Start

.

Important

Using previous version of Kubebuilder? Check the legacy documentation for

v1

,

v2

or

v3

.

Who is this for

Users of Kubernetes

Users of Kubernetes will develop a deeper understanding of Kubernetes through learning
the fundamental concepts behind how APIs are designed and implemented.  This book
will teach readers how to develop their own Kubernetes APIs and the
principles from which the core Kubernetes APIs are designed.

Including:

The structure of Kubernetes APIs and Resources

API versioning semantics

Self-healing

Garbage Collection and Finalizers

Declarative vs Imperative APIs

Level-Based vs Edge-Base APIs

Resources vs Subresources

Kubernetes API extension developers

API extension developers will learn the principles and concepts behind implementing canonical
Kubernetes APIs, as well as simple tools and libraries for rapid execution.  This
book covers pitfalls and misconceptions that extension developers commonly encounter.

Including:

How to batch multiple events into a single reconciliation call

How to configure periodic reconciliation

Forthcoming

When to use the lister cache vs live lookups

Garbage Collection vs Finalizers

How to use Declarative vs Webhook Validation

How to implement API versioning

Why Kubernetes APIs

Kubernetes APIs provide consistent and well defined endpoints for
objects adhering to a consistent and rich structure.

This approach has fostered a rich ecosystem of tools and libraries for working
with Kubernetes APIs.

Users work with the APIs through declaring objects as

yaml

or

json

config, and using
common tooling to manage the objects.

Building services as Kubernetes APIs provides many advantages to plain old REST, including:

Hosted API endpoints, storage, and validation.

Rich tooling and CLIs such as

kubectl

and

kustomize

.

Support for AuthN and granular AuthZ.

Support for API evolution through API versioning and conversion.

Facilitation of adaptive / self-healing APIs that continuously respond to changes
in the system state without user intervention.

Kubernetes as a hosting environment

Developers may build and publish their own Kubernetes APIs for installation into
running Kubernetes clusters.

Contribution

If you like to contribute to either this book or the code, please be so kind
to read our

Contribution

guidelines first.

Resources

Repository:

sigs.k8s.io/kubebuilder

Slack channel:

#kubebuilder

Google Group:

kubebuilder@googlegroups.com