# Contributing to AEP Monorepo

- [Contributing to AEP Monorepo](#contributing-to-aep-monorepo)
  - [Repository Layout](#repository-layout)
  - [Local Development \& Go Workspace](#local-development--go-workspace)
  - [CI/CD \& GitHub Actions Pipelines](#cicd--github-actions-pipelines)
  - [Tagging and Releases](#tagging-and-releases)
  - [Pull Request Guidelines](#pull-request-guidelines)
  - [References](#references)

## Repository Layout

This monorepo utilizes a flat, root-level module structure where each top-level directory corresponds to a distinct project:

```text
aep-monorepo/
├── DESIGN/                       # Architecture and design documents
├── aep-e2e-validator/            # End-to-end validator for AEP HTTP APIs
├── aep-lib-go/                   # Shared Go core libraries, types, and OpenAPI utilities
├── aepc/                         # Compiler/generator for AEP-compliant resource models
├── aepcli/                       # Command-line interface for interacting with AEP APIs
├── api-linter/                   # Protocol buffer linter for AEP rules
├── terraform-provider-aep/       # Dynamic Terraform provider for AEP APIs
├── go.work                       # Go Workspace configuration linking local modules
└── go.work.sum
```

- [aep-e2e-validator](aep-e2e-validator): This project contains the end-to-end conformance testing suite for AEP APIs.
- [aep-lib-go](aep-lib-go): This project contains the common Go library with core types, OpenAPI parsing, and resource schema representations.
- [aepc](aepc): This project contains the compiler and code generator for AEP services.
- [aepcli](aepcli): This project contains the command line interface tool for interacting with AEP-compliant APIs.
- [api-linter](api-linter): This project contains the protobuf linter enforcing AEP style guidelines.
- [terraform-provider-aep](terraform-provider-aep): This project contains the dynamic Terraform provider generated from AEP OpenAPI schemas.

## Local Development & Go Workspace

The monorepo uses Go workspaces configured via [go.work](go.work):
- Cross-module local development functions without requiring local replace directives in go.mod files.
- Edits in shared libraries like [aep-lib-go](aep-lib-go) are immediately visible to downstream tools such as [api-linter](api-linter), [aepcli](aepcli), and [aepc](aepc).

Common development workflows:
- Run all unit tests across all modules from the repository root:
  ```bash
  go test ./...
  ```
- Run tests for a specific module:
  ```bash
  go test ./aepcli/...
  ```
  Alternatively, run tests from within the module directory:
  ```bash
  cd aepcli && go test ./...
  ```
- Run linting for a specific module:
  ```bash
  cd api-linter && make lint
  ```

## CI/CD & GitHub Actions Pipelines

Workflows live in the root [.github/workflows](.github/workflows) directory:
- Pipelines are isolated per module using path-based triggers.
- Pull requests only execute workflows corresponding to modified modules.

## Tagging and Releases

Releases and Go module version tags use module-prefixed tags:
- Tag format: `{module}/v{version}`
  - Example: `api-linter/v1.2.0`
  - Example: `aepcli/v0.4.1`
  - Example: `aep-lib-go/v0.3.0`
  - Example: `terraform-provider-aep/v0.5.0`
- Pushing a prefixed tag triggers the corresponding release workflow to build artifacts and create a GitHub release.

## Pull Request Guidelines

- Verify that tests pass locally before opening a pull request:
  ```bash
  go test ./...
  ```
- When merging, a squash and rebase strategy is used. Context and iteration are preserved in the pull request.

## References

- Design document detailing monorepo initialization and git history grafting: [DESIGN/2026-08-16-repository-initialization-and-grafting.md](DESIGN/2026-08-16-repository-initialization-and-grafting.md)
- Root Go workspace configuration linking all submodules: [go.work](go.work)
- GitHub Actions workflow configurations directory: [.github/workflows](.github/workflows)
- AEP development guidelines and standards: https://aep.dev


