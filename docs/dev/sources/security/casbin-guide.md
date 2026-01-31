# Casbin Documentation

> Source: https://casbin.org/docs/overview
> Fetched: 2026-01-30T23:54:40.660852+00:00
> Content-Hash: 1d59bc2f66c75a10
> Type: html

---

Basics

Overview

Casbin Overview

H

hsluoyz

·

5 min

read

Copy Markdown

Open

Feedback

Casbin is an efficient, open-source access control library designed to enforce authorization through support for multiple

access control models

.

Feedback

Implementing rule-based access control is straightforward: define subjects, objects, and permitted actions in a

policy

file using any format that suits your requirements. This pattern remains consistent throughout all Casbin implementations. Through the

model

file, developers and administrators maintain full authority over authorization logic, including layout, execution flow, and conditional requirements. To validate incoming requests against your defined policy and model files, Casbin provides an

Enforcer

component.

Languages Supported by Casbin

Feedback

Casbin offers native support across multiple programming languages, enabling seamless integration into diverse projects and workflows:

Casbin

jCasbin

node-Casbin

PHP-Casbin

Production-ready

Production-ready

Production-ready

Production-ready

PyCasbin

Casbin.NET

Casbin-CPP

Casbin-RS

Production-ready

Production-ready

Production-ready

Production-ready

Feature Set for Different Languages

Feedback

Our goal is feature parity across all language implementations, though we haven't achieved complete uniformity yet.

Feature

Go

Java

Node.js

PHP

Python

C#

Delphi

Rust

C++

Lua

Dart

Elixir

Enforcement

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

RBAC

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

ABAC

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

Scaling ABAC (

eval()

)

✅

✅

✅

✅

✅

✅

❌

✅

✅

✅

✅

✅

Adapter

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

❌

Management API

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

RBAC API

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

Batch API

✅

✅

✅

✅

✅

✅

❌

✅

✅

✅

❌

❌

Filtered Adapter

✅

✅

✅

✅

✅

✅

❌

✅

✅

✅

❌

❌

Watcher

✅

✅

✅

✅

✅

✅

✅

✅

✅

✅

❌

❌

Role Manager

✅

✅

✅

✅

✅

✅

❌

✅

✅

✅

✅

❌

Multi-Threading

✅

✅

✅

❌

✅

❌

❌

✅

❌

❌

❌

❌

'in' of matcher

✅

✅

✅

✅

✅

❌

✅

❌

❌

❌

✅

✅

Feedback

Note

- A checkmark (✅) for Watcher or Role Manager indicates that the interface exists in the core library, not necessarily that an implementation is available.

What is Casbin?

Feedback

Casbin serves as an authorization library for scenarios requiring controlled access to resources. In typical usage, a

subject

(user or service) requests access to an

object

(resource or entity) to perform an

action

(such as

read

,

write

, or

delete

). Developers define these actions according to their application needs. This represents the "standard" or classic

{ subject, object, action }

authorization flow that Casbin handles most commonly.

Feedback

Beyond this standard model, Casbin accommodates complex authorization scenarios by supporting

roles (RBAC)

,

attributes (ABAC)

, and other advanced patterns.

What Casbin Does

Feedback

Applies policy enforcement in the classic

{ subject, object, action }

format or any custom format you define, supporting both allow and deny authorizations.

Manages storage for the access control model and associated policies.

Handles user-role and role-role relationships (the role hierarchy concept in RBAC).

Recognizes built-in superusers such as

root

or

administrator

who have unrestricted access without requiring explicit permission rules.

Supplies various built-in operators for pattern matching in rules—for instance,

keyMatch

matches resource key

/foo/bar

to pattern

/foo*

.

What Casbin Does

NOT

Do

Feedback

User authentication (validating

username

and

password

credentials during login)

User or role list management

Feedback

Most applications already manage their own user accounts, roles, and credentials. Casbin wasn't designed as a password storage system—it focuses solely on authorization. That said, Casbin does maintain user-role associations when operating in RBAC mode.

How is this guide?

Good

Bad