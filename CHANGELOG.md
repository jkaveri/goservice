# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2026-05-12

### Added

- Formatted error constructors in `errorcode` (`NewErrorf`, `NotFoundf`,
  `InvalidRequestf`, and the other `*f` variants) for `fmt`-style messages
  without tripping `go vet` on dynamic plain-text messages ([#1](https://github.com/jkaveri/goservice/pull/1)).

### Changed

- Mockery: generated mocks are organized under dedicated `mock/` subdirectories
  per package.
- `.gitignore` updates.

## [1.0.1] - 2026-04-16

### Changed

- Tag-only release on the same revision as v1.0.0 (module / release alignment).

## [1.0.0] - 2026-04-16

### Added

- Initial public `github.com/jkaveri/goservice` module: gRPC and gateway helpers,
  validation utilities, structured errors, logging integration, and repo
  tooling.

[Unreleased]: https://github.com/jkaveri/goservice/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/jkaveri/goservice/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/jkaveri/goservice/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/jkaveri/goservice/releases/tag/v1.0.0
