# Contributing to AEP Monorepo

Welcome to the `aep-monorepo` codebase! This repository consolidates multiple AEP tools, libraries, linters, and plugins into a single, cohesive codebase.

---

## 1. Repository Layout

This monorepo utilizes a **flat, root-level module structure**. Each top-level directory corresponds to a distinct project/module:

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

---

## 2. Local Development & Go Workspace

The monorepo uses Go workspaces (`go.work`). This enables seamless cross-module development: changes made to `aep-lib-go` are immediately recognized by `api-linter`, `aepcli`, and `aepc` without needing module tag releases or local `replace` directives.

### Common Commands

* **Run all unit tests across all modules**:
  ```bash
  go test ./...
  ```
* **Work on / test a specific module**:
  ```bash
  cd aepcli
  go test ./...
  # or from repository root:
  go test ./aepcli/...
  ```
* **Run linting (api-linter)**:
  ```bash
  cd api-linter
  make lint
  ```

---

## 3. CI/CD & GitHub Actions Pipelines

Workflows live in the root `.github/workflows/` directory. Each project has **independent, isolated CI/CD pipelines** that trigger selectively based on path filtering:

| Workflow File | Target Project | Trigger Conditions |
| :--- | :--- | :--- |
| `api-linter-ci.yml` | `api-linter/` | Pull requests & pushes modifying `api-linter/**` |
| `api-linter-release.yml` | `api-linter/` | Tag push matching `api-linter/v*` |
| `aepc-ci.yml` | `aepc/` | Pull requests & pushes modifying `aepc/**` |
| `aepcli-ci.yml` | `aepcli/` | Pull requests & pushes modifying `aepcli/**` |
| `aepcli-release.yml` | `aepcli/` | Tag push matching `aepcli/v*` |
| `aep-lib-go-ci.yml` | `aep-lib-go/` | Pull requests & pushes modifying `aep-lib-go/**` |
| `aep-e2e-validator-ci.yml` | `aep-e2e-validator/` | Pull requests & pushes modifying `aep-e2e-validator/**` |
| `terraform-provider-aep-ci.yml` | `terraform-provider-aep/` | Pull requests & pushes modifying `terraform-provider-aep/**` |
| `terraform-provider-aep-integration.yml` | `terraform-provider-aep/` | Pull requests & pushes modifying `terraform-provider-aep/**` |
| `terraform-provider-aep-release.yml` | `terraform-provider-aep/` | Tag push matching `terraform-provider-aep/v*` |

---

## 4. Tagging and Releases

Because this repository houses multiple Go modules, releases and Go module versions use **prefixed tags**:

* `api-linter/v1.2.0`
* `aepcli/v0.4.1`
* `aep-lib-go/v0.3.0`
* `terraform-provider-aep/v0.5.0`

Pushing a prefixed tag will automatically trigger the corresponding release workflow and publish binary assets and release notes.

---

## 5. Pull Request Guidelines

1. **Keep PRs focused**: Wherever possible, keep changes localized to a single module or group of tightly coupled modules.
2. **Ensure tests pass**: Verify that your changes compile and pass tests (`go test ./<module>/...`) before opening a pull request.
3. **Preserve audit history**: Never rebase or squash commits in a way that destroys historical attribution.
