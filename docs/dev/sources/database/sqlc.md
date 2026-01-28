# sqlc

> Auto-fetched from [https://docs.sqlc.dev/en/stable/](https://docs.sqlc.dev/en/stable/)
> Last Updated: 2026-01-28T21:44:12.494533+00:00

---

Overview
Installing sqlc
Tutorials
Getting started with MySQL
Getting started with PostgreSQL
Getting started with SQLite
Commands
generate
- Generating code
push
- Uploading projects
verify
- Verifying schema changes
vet
- Linting queries
How-to Guides
Retrieving rows
Counting rows
Inserting rows
Updating rows
Deleting rows
Preparing queries
Using transactions
Naming parameters
Modifying the database schema
Configuring generated structs
Embedding structs
Overriding types
Renaming fields
sqlc Cloud
Managed databases
Reference
Changelog
CLI
Configuration
Datatypes
Environment variables
Database and language support
Macros
Query annotations
Conceptual Guides
Using sqlc in CI/CD
Using Go and pgx
Using plugins
Developing sqlc
Privacy and data collection
Sponsored By
sqlc
sqlc Documentation
View page source
sqlc Documentation

And lo, the Great One looked down upon the people and proclaimed:
“SQL is actually pretty great”
sqlc generates
fully type-safe idiomatic Go code
from SQL. Here’s how it
works:
You write SQL queries
You run sqlc to generate Go code that presents type-safe interfaces to those
queries
You write application code that calls the methods sqlc generated
Seriously, it’s that easy. You don’t have to write any boilerplate SQL querying
code ever again.
Next
© Copyright 2024, Riza, Inc..
Built with
Sphinx
using a
theme
provided by
Read the Docs
.