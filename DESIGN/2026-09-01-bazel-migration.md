# Bazel Build System Migration & Tooling Architecture

* **Date**: 2026-09-01
* **Status**: Completed

## 1. Overview

This document describes the design and implementation of the [Bazel](https://bazel.build/) build system for `aep-monorepo`. Adopting Bazel provides hermetic, reproducible, fast incremental builds and uniform test execution across all projects in the monorepo.

---

## 2. Architecture & Design Principles

### Bzlmod-First Dependency Management
The monorepo uses Bazel 7 with Bzlmod ([MODULE.bazel](MODULE.bazel)) rather than legacy `WORKSPACE` files:
- **`rules_go`**: Compiles Go packages, tests, and binaries hermetically.
- **`gazelle`**: Automatically scans Go code and generates idiomatic `BUILD.bazel` files.
- **`go_sdk`**: Pinned to Go 1.24.0 toolchain to ensure identical builds across local workstations and CI environments.

### Gazelle Namespace & Import Mapping
Gazelle is configured at the monorepo root:
```starlark
# gazelle:prefix github.com/aep-dev/aep-monorepo
# gazelle:proto disable_global
# gazelle:exclude **/testdata
```
This enables Gazelle to map import paths across modules (e.g. `github.com/aep-dev/aep-monorepo/aep-lib-go/pkg/api`) to Bazel package labels (`//aep-lib-go/pkg/api:api`).

### Dependency Conflict Resolution
To prevent conflicts between `wellknownimports` in `protocompile` and `google.golang.org/protobuf`, Gazelle overrides are configured in `MODULE.bazel`:
```starlark
go_deps.gazelle_override(
    path = "github.com/bufbuild/protocompile",
    directives = [
        "gazelle:resolve go google.golang.org/protobuf/types/pluginpb @org_golang_google_protobuf//types/pluginpb",
        "gazelle:resolve go google.golang.org/protobuf/types/descriptorpb @org_golang_google_protobuf//types/descriptorpb",
        "gazelle:resolve go google.golang.org/protobuf/reflect/protoreflect @org_golang_google_protobuf//reflect/protoreflect",
        "gazelle:resolve go google.golang.org/protobuf/proto @org_golang_google_protobuf//proto",
    ],
)
```

---

## 3. Target Hierarchy

The monorepo provides the following standard target structure:

| Target Pattern            | Description                                                            |
| :------------------------ | :--------------------------------------------------------------------- |
| `//...`                   | All packages, binaries, and tests in the monorepo                      |
| `//:fix`                  | Alias to `//scripts:fix` (runs dependency synchronization)             |
| `//:check`                | Alias to `//scripts:check` (verifies dependency synchronization)       |
| `//:gazelle`              | Runs Gazelle to generate and update all `BUILD.bazel` files            |
| `//aep-lib-go/...`        | Core library packages and tests                                        |
| `//aepcli/...`            | CLI library, test suite, and binary `//aepcli/cmd/aepcli:aepcli`       |
| `//api-linter/...`        | Linter rules, engine tests, and binaries `//api-linter/cmd/api-linter` |
| `//aepc/...`              | Compiler, CEL-to-SQL packages, and binary `//aepc:aepc`                |
| `//aep-e2e-validator/...` | Conformance suite, tests, and binary `//aep-e2e-validator/cmd:cmd`     |

---

## 4. Developer Workflows & `justfile` Integration

A [justfile](justfile) proxies standard tasks directly to Bazel targets:

```just
# Run all tests across the monorepo
test:
    bazel test //...

# Build all targets
build:
    bazel build //...

# Update Bazel BUILD files with Gazelle
gazelle:
    bazel run //:gazelle

# Synchronize internal dependencies and update BUILD files
fix:
    bazel run //:fix
    bazel run //:gazelle

# Verify internal dependencies are synchronized
check:
    bazel run //:check
```

---

## 5. Verification Commands

```bash
# Verify all targets build cleanly
bazel build //...

# Run all unit tests
bazel test //...

# Run dependency synchronization check
bazel run //:check
```

