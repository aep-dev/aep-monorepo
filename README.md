# AEP Monorepo

- This repository contains the consolidated source code and tooling ecosystem for [API Enhancement Proposals (AEP)](https://aep.dev).
- It is built and tested using [Bazel](https://bazel.build/) with Bzlmod and Gazelle, and supports Go workspaces and TypeScript for native development.

## Projects

- [aep-e2e-validator](aep-e2e-validator): Runtime end-to-end conformance testing tool for AEP HTTP APIs.
- [aep-lib-go](aep-lib-go): Shared Go core libraries, common types, OpenAPI parsers, and resource schema representations.
- [aep-lib-ts](aep-lib-ts): Shared TypeScript core libraries, case conversion utilities, OpenAPI schemas, and API client helpers.
- [aepc](aepc): Resource-oriented compiler and code generator for AEP services, protocol buffers, and OpenAPI specifications.
- [aepcli](aepcli): Dynamic command-line interface tool generated from OpenAPI definitions for AEP APIs.
- [api-linter](api-linter): Protocol buffer linter enforcing AEP design rules and naming standards.

## Getting Started

- Clone the repository:
  ```bash
  git clone https://github.com/aep-dev/aep-monorepo.git
  cd aep-monorepo
  ```
- Run all tests across the monorepo using Bazel or Just:
  ```bash
  just test
  # or: bazel test //...
  ```
- Build all targets across the workspace:
  ```bash
  just build
  # or: bazel build //...
  ```
- Run tests for an individual module or target:
  ```bash
  bazel test //aepcli/...
  bazel test //aep-lib-go/...
  bazel test //aep-lib-ts/...
  ```
- Synchronize internal module dependencies and update BUILD files:
  ```bash
  just fix
  ```

## Documentation

- Guidelines for contributing, local development, module tagging, and pull requests: [CONTRIBUTING.md](CONTRIBUTING.md)
- Bazel architecture and migration design document: [DESIGN/2026-09-01-bazel-migration.md](DESIGN/2026-09-01-bazel-migration.md)
- TypeScript library graft design document: [DESIGN/2026-09-01-aep-lib-ts-grafting.md](DESIGN/2026-09-01-aep-lib-ts-grafting.md)

## References

- API Enhancement Proposals standard documentation and specifications: https://aep.dev
- Root Bazel module configuration: [MODULE.bazel](MODULE.bazel)
- Root Go workspace configuration linking all submodules: [go.work](go.work)
- Project contribution and development guide: [CONTRIBUTING.md](CONTRIBUTING.md)
- GitHub Actions CI/CD workflows directory: [.github/workflows](.github/workflows)
