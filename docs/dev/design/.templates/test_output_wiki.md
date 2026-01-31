# Example Feature


<!-- SOURCES: fx, pgx, sqlc -->

<!-- DESIGN: 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES -->

> Example feature for testing documentation templates

---

## Status

| Dimension | Status |  |
|-----------|--------||
| Design | 🟡 |  |
| Sources | ✅ |  |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |


---

## Overview

This is an example feature used to demonstrate how the Jinja2 template system
generates both Claude-optimized and Wiki documentation from a single source.


**Key Features**:
- Template-based documentation
- Dual output (Claude + Wiki)
- SOURCE_OF_TRUTH integration


**Use Cases**:
- Test template rendering
- Validate documentation pipeline
- Example for new features

---


## Getting Started

This example feature is for documentation testing only.
In a real feature, this would contain step-by-step setup instructions.


### Prerequisites

- Revenge server installed
- Database configured

### Installation

```bash
# No installation needed for this example
# In real features, this would contain actual setup steps

```

### Configuration

No configuration needed for this example feature.


### Usage Examples

#### Create Example

How to create a new example

```
curl -X POST http://localhost:8096/api/v1/examples \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "My Example"}'

```

---

## Screenshots


---

## Troubleshooting

### Cannot create example

**Symptoms**: 403 Forbidden error

**Solution**: Ensure you have example:write scope in your token

---

## Community

- **GitHub Issues**: https://github.com/lusoris/revenge/issues
- **Discussions**: https://github.com/lusoris/revenge/discussions
- **Documentation**: https://github.com/lusoris/revenge/wiki

---

