# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added

- Initial draft of CHANGELOG.md
- Planned features and improvements

---

## [1.0.0] - 2025-12-14

### Added

- First stable release of `jcscli`.
- Canonical JSON encoding (RFC 8785 compliant).
- Pretty-print option (`--pretty` / `-p`).
- Quiet mode (`--quiet` / `-q`) and verbose mode (`--verbose` / `-v`).
- Interactive mode (`--interactive` / `-i`) for typing JSON directly.
- Safe overwrite handling (`--overwrite` / `-w`).
- Version flag (`--version` / `-V`).
- Help flag (`--help` / `-h`) with usage examples.
- Exit codes: `0` success, `1` runtime error, `2` usage error.

---

## [0.1.0] - 2025-11-30

### Added

- Prototype CLI with basic canonicalization.
- File input (`--file` / `-f`) and output (`--output` / `-o`) support.

---

## [0.0.1] - 2025-11-15

### Added

- Initial proof-of-concept for canonical JSON library.
