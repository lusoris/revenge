# testcontainers-go

> Auto-fetched from [https://golang.testcontainers.org/](https://golang.testcontainers.org/)
> Last Updated: 2026-01-28T21:47:10.725736+00:00

---

Content
testcontainers-go
Home
Home
Table of contents
About Testcontainers for Go
GoDoc
Who is using Testcontainers Go?
License
Copyright
Quickstart
Features
Features
How to create a container
Wait Strategies
Wait Strategies
Introduction
Exec
Exit
File
Health
HostPort
HTTP
Log
Multi
SQL
TLS
Walk
Copying data into a container
Following Container Logs
Garbage Collector
Build from Dockerfile
Executing commands
Networking and communicating with containers
Custom configuration
Image name substitution
Test Session Semantics
Authentication with Docker
Using Docker Compose
TLS certificates
Modules
Modules
Testcontainers for Go modules
Aerospike
ArangoDB
Apache ActiveMQ Artemis
Azure
Azurite
Cassandra
Chroma
ClickHouse
CockroachDB
Consul
Couchbase
Databend
DinD (Docker in Docker)
Docker MCP Gateway
Docker Model Runner
Dolt
DynamoDB
Elasticsearch
etcd
GCloud
Grafana LGTM
Inbucket
InfluxDB
K3s
K6
Kafka (KRaft)
LocalStack
MariaDB
Meilisearch
Memcached
Milvus
Minio
MockServer
MongoDB Atlas Local
MongoDB
MS SQL Server
MySQL
NATS
NebulaGraph
Neo4j
Ollama
OpenFGA
OpenLDAP
OpenSearch
Pinecone
Postgres
Apache Pulsar
Qdrant
RabbitMQ
Redis
Redpanda
Registry
ScyllaDB
Socat
Solace Pubsub+
SurrealDB
Toxiproxy
Valkey
Vault
Vearch
Weaviate
YugabyteDB
Examples
Examples
Code examples
Nginx
System Requirements
System Requirements
Go version
General Docker requirements
Continuous Integration
Continuous Integration
AWS CodeBuild
Bitbucket Pipelines
CircleCI
Concourse CI
Patterns for running tests inside a Docker container
Drone CI
GitLab CI
Tekton
Travis
Using Colima with Docker
Using Podman instead of Docker
Using Rancher Desktop
Usage Metrics
Dependabot
Contributing
Getting help
Join the community
Table of contents
About Testcontainers for Go
GoDoc
Who is using Testcontainers Go?
License
Copyright
Welcome to Testcontainers for Go!
¶
Not using Go? Here are other supported languages!
Java
Go
.NET
Node.js
Python
Rust
Haskell
Ruby
About Testcontainers for Go
¶
Testcontainers for Go
is a Go package that makes it simple to create and clean up container-based dependencies for
automated integration/smoke tests. The clean, easy-to-use API enables developers to programmatically define containers
that should be run as part of a test and clean up those resources when the test is done.
To start using
Testcontainers for Go
please read our
quickstart guide
.
GoDoc
¶
Inline documentation and docs where the code live is crucial for us. Go has nice support for them and we provide
examples as well. Check it out at
pkg.go.dev/github.com/testcontainers/testcontainers-go
.
Who is using Testcontainers Go?
¶
Elastic
- Testing of the APM Server, and E2E testing for Beats
Telegraf
- Integration testing the plugin-driven server agent for collecting & reporting metrics
Intel
- Reference implementation design E2E testing for microservice-based solutions
OpenTelemetry
- Integration testing of the OpenTelemetry Collector receivers
License
¶
This project is opensource and you can have a look at the code on
GitHub
. See
LICENSE
.
Copyright
¶
Copyright (c) 2018-present Gianluca Arbezzano and other authors. Check out our
lovely contributors
.