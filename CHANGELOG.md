# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- New error codes in `errorcode` for common operational failures:
  `too_many_requests`, `timeout`, `unavailable`, `unimplemented`, and
  `failed_precondition`. Each ships with a constructor, an `*f` variant,
  an `Is*` helper, and a `With*` wrapper, plus matching HTTP and gRPC
  mappings in `grpc/interceptors/wraperror`.

## [1.2.0] - 2026-05-12

### Changes
- Add errorcode.Wrap and Wrapf with unit tests. (25f4996)
- remove unused command (1b8253e)
- .cursor rules (a9c3325)
- feat(errorcode): add too_many_requests, timeout, unavailable, unimplemented, failed_precondition (00e068e)
- commit code (38d77d7)
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

[Unreleased]: https://github.com/jkaveri/goservice/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/jkaveri/goservice/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/jkaveri/goservice/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/jkaveri/goservice/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/jkaveri/goservice/releases/tag/v1.0.0
