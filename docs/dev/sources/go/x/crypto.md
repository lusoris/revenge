# Go x/crypto Package

> Source: https://pkg.go.dev/golang.org/x/crypto
> Fetched: 2026-01-31
> Content-Hash: auto-generated
> Type: html

---

## Overview

**Module:** `golang.org/x/crypto`
**Version:** v0.47.0
**License:** BSD-3-Clause
**Repository:** https://cs.opensource.google/go/x/crypto

This repository holds supplementary Go cryptography packages providing advanced cryptographic implementations beyond the standard library.

---

## Major Sub-Packages

### Key Management & Certificates

| Package | Description |
|---------|-------------|
| `acme` | ACME spec implementation for Let's Encrypt |
| `acme/autocert` | Automatic certificate access from ACME-based CAs |
| `ocsp` | OCSP response parsing (RFC 2560) |

### Password Hashing & Key Derivation

| Package | Description |
|---------|-------------|
| `argon2` | Argon2 key derivation function |
| `bcrypt` | Provos and Mazieres's bcrypt adaptive hashing |
| `pbkdf2` | PBKDF2 key derivation (RFC 2898) |
| `scrypt` | scrypt key derivation function |
| `hkdf` | HMAC-based Extract-and-Expand Key Derivation (RFC 5869) |

### Symmetric Encryption

| Package | Description |
|---------|-------------|
| `chacha20` | ChaCha20 and XChaCha20 (RFC 8439) |
| `chacha20poly1305` | ChaCha20-Poly1305 AEAD and XChaCha20-Poly1305 |
| `blowfish` | Bruce Schneier's Blowfish encryption |
| `twofish` | Bruce Schneier's Twofish encryption |
| `cast5` | CAST5 (RFC 2144) |
| `tea` | TEA encryption algorithm |
| `xtea` | XTEA encryption |
| `xts` | XTS cipher mode (IEEE P1619/D16) |
| `salsa20` | Salsa20 stream cipher |

### Elliptic Curve Cryptography

| Package | Description |
|---------|-------------|
| `curve25519` | X25519 scalar multiplication (RFC 7748) |
| `ed25519` | Ed25519 signature algorithm |
| `bn256` | Bilinear group implementation |

### Hash Functions

| Package | Description |
|---------|-------------|
| `blake2b` | BLAKE2b hash (RFC 7693) |
| `blake2s` | BLAKE2s hash (RFC 7693) |
| `sha3` | SHA-3 hash and SHAKE (FIPS 202) |
| `md4` | MD4 hash (RFC 1320) |
| `ripemd160` | RIPEMD-160 hash |

### Message Authentication

| Package | Description |
|---------|-------------|
| `poly1305` | Poly1305 one-time message authentication code |
| `nacl/auth` | Secret key authentication |

### Public-Key Cryptography

| Package | Description |
|---------|-------------|
| `nacl/box` | Public-key authenticated encryption |
| `nacl/sign` | Public-key signatures |
| `nacl/secretbox` | Secret key authenticated encryption |
| `openpgp` | High-level OpenPGP message operations |

### SSH

| Package | Description |
|---------|-------------|
| `ssh` | SSH client and server implementation |
| `ssh/agent` | ssh-agent protocol client and server |
| `ssh/knownhosts` | known_hosts parser and writer |
| `ssh/terminal` | Terminal support functions |

### Other

| Package | Description |
|---------|-------------|
| `otr` | Off The Record protocol |
| `pkcs12` | PKCS#12 implementation |
| `cryptobyte` | Length-prefixed binary message parsing |

---

## Status & Best Practices

- Valid go.mod - Official Go module dependency management
- Redistributable License - BSD-3-Clause
- Tagged Versions - Predictable builds
- Pre-v1 - Not yet stable, may have breaking changes

---

## Resources

- **Report Issues:** Use https://go.dev/issues with "x/crypto:" prefix
- **Contribute:** See https://go.dev/doc/contribute
- **Repository:** https://go.googlesource.com/crypto
- **Security:** https://go.dev/security/policy
