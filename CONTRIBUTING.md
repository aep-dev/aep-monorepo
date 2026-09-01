# Contributing to AEP Monorepo

- [Contributing to AEP Monorepo](#contributing-to-aep-monorepo)
  - [Repository Layout](#repository-layout)
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
├── go.work                       # Go Workspace configuration linking local modules
└── go.work.sum
```

- [aep-e2e-validator](aep-e2e-validator): This project contains the end-to-end conformance testing suite for AEP APIs.
- [aep-lib-go](aep-lib-go): This project contains the common Go library with core types, OpenAPI parsing, and resource schema representations.
- [aepc](aepc): This project contains the compiler and code generator for AEP services.
- [aepcli](aepcli): This project contains the command line interface tool for interacting with AEP-compliant APIs.
- [api-linter](api-linter): This project contains the protobuf linter enforcing AEP style guidelines.

## Local Development & Go Workspace

The monorepo uses Go workspaces configured via [go.work](go.work):
- Cross-module local development functions without requiring local replace directives in go.mod files.
- Edits in shared libraries like [aep-lib-go](aep-lib-go) are immediately visible to downstream tools such as [api-linter](api-linter), [aepcli](aepcli), and [aepc](aepc).

Common development workflows:
- Run all unit tests across all modules from the repository root:
  ```bash
  just test
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

## Dependency Management & Synchronization

- Changes to shared libraries like [aep-lib-go](aep-lib-go) must be synchronized in the `go.mod` files of all dependent modules (`aepcli`, `aepc`, `aep-e2e-validator`).
- Synchronize all dependent module `go.mod` files automatically:
  ```bash
  just fix
  ```
- Verify whether dependent `go.mod` files are synchronized against the target branch:
  ```bash
  just check
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

- Verify that tests pass locally before opening a pull request:
  ```bash
  just test
  ```
- If modifying [aep-lib-go](aep-lib-go), run the dependency synchronization check before opening a PR:
  ```bash
  just check
  ```
- When merging, a squash and rebase strategy is used. Context and iteration are preserved in the pull request.

## References

- Design document detailing monorepo initialization and git history grafting: [DESIGN/2026-08-16-repository-initialization-and-grafting.md](DESIGN/2026-08-16-repository-initialization-and-grafting.md)
- Root Go workspace configuration linking all submodules: [go.work](go.work)
- Command runner recipe definitions: [justfile](justfile)
- GitHub Actions workflow configurations directory: [.github/workflows](.github/workflows)
- Internal dependency synchronization helper script: [scripts/sync-deps.py](scripts/sync-deps.py)
- Conventional Commits specification: https://www.conventionalcommits.org/
- AEP development guidelines and standards: https://aep.dev



