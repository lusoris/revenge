# testify

> Source: https://pkg.go.dev/github.com/stretchr/testify
> Fetched: 2026-01-31T11:02:07.863520+00:00
> Content-Hash: 606a7045ca93b0b1
> Type: html

---

### Overview ¶

Module testify is a set of packages that provide many tools for testifying that your code will behave as you intend. 

Testify contains the following packages: 

The [github.com/stretchr/testify/assert](/github.com/stretchr/testify@v1.11.1/assert) package provides a comprehensive set of assertion functions that tie in to [the Go testing system](https://go.dev/doc/code#Testing). The [github.com/stretchr/testify/require](/github.com/stretchr/testify@v1.11.1/require) package provides the same assertions but as fatal checks. 

The [github.com/stretchr/testify/mock](/github.com/stretchr/testify@v1.11.1/mock) package provides a system by which it is possible to mock your objects and verify calls are happening as expected. 

The [github.com/stretchr/testify/suite](/github.com/stretchr/testify@v1.11.1/suite) package provides a basic structure for using structs as testing suites, and methods on those structs as tests. It includes setup/teardown functionality in the way of interfaces. 

A [golangci-lint](https://golangci-lint.run/) compatible linter for testify is available called [testifylint](https://github.com/Antonboom/testifylint). 
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
