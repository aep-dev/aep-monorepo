# Contributing to AEP Monorepo

- [Contributing to AEP Monorepo](#contributing-to-aep-monorepo)
  - [Repository Layout](#repository-layout)
  - [Build System \& Bazel](#build-system--bazel)
  - [Local Development \& Go Workspace](#local-development--go-workspace)
  - [Dependency Management \& Synchronization](#dependency-management--synchronization)
  - [CI/CD \& GitHub Actions Pipelines](#cicd--github-actions-pipelines)
  - [Tagging and Releases](#tagging-and-releases)
  - [Conventional Commits](#conventional-commits)
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
├── scripts/                      # Shared developer scripts and automation tools
├── MODULE.bazel                  # Bazel Bzlmod configuration
├── BUILD.bazel                   # Root Bazel build definitions and Gazelle config
├── go.work                       # Go Workspace configuration linking local modules
└── go.work.sum
```

- [aep-e2e-validator](aep-e2e-validator): Runtime end-to-end conformance testing suite for AEP APIs.
- [aep-lib-go](aep-lib-go): Common Go library with core types, OpenAPI parsing, and resource schema representations.
- [aepc](aepc): Compiler and code generator for AEP services.
- [aepcli](aepcli): Command-line interface tool for interacting with AEP-compliant APIs.
- [api-linter](api-linter): Protobuf linter enforcing AEP design rules and naming standards.

## Build System & Bazel

The repository uses [Bazel](https://bazel.build/) (v7.4.1) as the unified polyglot build and test system, configured with Bzlmod ([MODULE.bazel](MODULE.bazel)) and Gazelle.

Common Bazel workflows (managed via [justfile](justfile)):
- Run all tests across all packages:
  ```bash
  just test
  # or: bazel test //...
  ```
- Build all binaries and packages:
  ```bash
  just build
  # or: bazel build //...
  ```
- Run tests for a specific module or package:
  ```bash
  bazel test //aepcli/...
  bazel test //aep-lib-go/...
  bazel test //api-linter/...
  ```
- Build a specific binary target:
  ```bash
  bazel build //aepcli/cmd/aepcli:aepcli
  bazel build //api-linter/cmd/api-linter:api-linter
  bazel build //aepc:aepc
  ```
- Regenerate or update Bazel `BUILD.bazel` files with Gazelle:
  ```bash
  just gazelle
  # or: bazel run //:gazelle
  ```

## Local Development & Go Workspace

For IDE integration (`gopls`, VSCode, GoLand), the repository maintains Go workspaces configured via [go.work](go.work):
- Cross-module local development functions without requiring local replace directives in go.mod files.
- Edits in shared libraries like [aep-lib-go](aep-lib-go) are immediately visible to downstream tools such as [api-linter](api-linter), [aepcli](aepcli), and [aepc](aepc).

## Dependency Management & Synchronization

- Changes to shared libraries like [aep-lib-go](aep-lib-go) must be synchronized across dependent modules (`aepcli`, `aepc`, `aep-e2e-validator`).
- Synchronize all dependent module `go.mod` files and Bazel definitions automatically:
  ```bash
  just fix
  # or: bazel run //:fix && bazel run //:gazelle
  ```
- Verify whether dependent `go.mod` files are synchronized against the target branch:
  ```bash
  just check
  # or: bazel run //:check
  ```

## CI/CD & GitHub Actions Pipelines

Workflows live in the root [.github/workflows](.github/workflows) directory:
- [bazel-ci.yml](.github/workflows/bazel-ci.yml) validates all test and build targets on pull requests and pushes to `main`.
- Module-specific pipelines provide path-based triggers and binary release packaging.

## Tagging and Releases

Releases and Go module version tags use module-prefixed tags:
- Tag format: `{module}/v{version}`
  - Example: `api-linter/v1.2.0`
  - Example: `aepcli/v0.4.1`
  - Example: `aep-lib-go/v0.3.0`
- Pushing a prefixed tag triggers the corresponding release workflow to build artifacts and create a GitHub release.

## Conventional Commits

Commit messages and pull request titles should follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:
- Structure commit messages as `<type>(<scope>): <description>`.
- Common types:
  - `feat`: A new feature or capability.
  - `fix`: A bug fix.
  - `docs`: Documentation changes.
  - `refactor`: Code changes that neither fix a bug nor add a feature.
  - `test`: Adding or correcting tests.
  - `chore`: Maintenance tasks, dependency updates, and tooling improvements.
- Common scopes include module names or subsystem areas (e.g. `feat(aepcli): add support for custom headers`, `fix(api-linter): correct resource name check`).

## Pull Request Guidelines

- Verify that tests and dependency checks pass locally before opening a pull request:
  ```bash
  just check
  just test
  ```
- When merging, a squash and rebase strategy is used. Context and iteration are preserved in the pull request.

## References

- Design document detailing monorepo initialization and git history grafting: [DESIGN/2026-08-16-repository-initialization-and-grafting.md](DESIGN/2026-08-16-repository-initialization-and-grafting.md)
- Design document detailing Bazel build system adoption: [DESIGN/2026-09-01-bazel-migration.md](DESIGN/2026-09-01-bazel-migration.md)
- Root Bazel module configuration: [MODULE.bazel](MODULE.bazel)
- Root Go workspace configuration linking all submodules: [go.work](go.work)
- Command runner recipe definitions: [justfile](justfile)
- GitHub Actions workflow configurations directory: [.github/workflows](.github/workflows)
- Internal dependency synchronization helper script: [scripts/sync-deps.py](scripts/sync-deps.py)
- Conventional Commits specification: https://www.conventionalcommits.org/
- AEP development guidelines and standards: https://aep.dev
