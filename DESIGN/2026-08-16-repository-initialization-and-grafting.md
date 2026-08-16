# Repository Initialization and History Grafting

* **Date**: 2026-08-16
* **Status**: Completed

## 1. Overview

This document outlines the architecture, directory structure, and Git history migration strategy used to initialize `aep-monorepo` by consolidating multiple standalone AEP Go repositories into a unified multi-module monorepo while fully preserving individual commit histories, author metadata, diffs, and audit paths.

---

## 2. Monorepo Layout & Rationale

### Directory Structure

The repository uses a flat, root-level module structure:

```text
aep-monorepo/
├── DESIGN/
│   └── 2026-08-16-repository-initialization-and-grafting.md
├── aep-e2e-validator/
│   ├── go.mod
│   └── ...
├── aep-lib-go/
│   ├── go.mod
│   └── ...
├── aepc/
│   ├── go.mod
│   └── ...
├── aepcli/
│   ├── go.mod
│   └── ...
├── api-linter/
│   ├── go.mod
│   └── ...
├── terraform-provider-aep/
│   ├── go.mod
│   └── ...
├── go.work
└── go.work.sum
```

### Why Root-Level Modules?

1. **Go Multi-Module Releases**: Go submodule tagging uses the pattern `<dir>/vX.Y.Z` (e.g. `api-linter/v1.0.0`, `aepcli/v0.4.0`), keeping tag names aligned 1:1 with previous project names.
2. **Predictable Import Paths**: External Go consumers can import modules directly without intermediate taxonomy layers (e.g. `github.com/aep-dev/aep-monorepo/aep-lib-go`).
3. **Future Bazel Compatibility**: Flat top-level directories naturally map to root Bazel packages (`//aepcli:...`, `//api-linter:...`) and make it straightforward to adopt polyglot tooling (TypeScript, Rust, Python) alongside Go.

---

## 3. Git History Grafting Methodology

To ensure that running `git log <directory>`, `git log <file_path>`, or `git blame <file_path>` shows the complete chronological history of every file back to its inception, the migration applied **path-rewritten history grafting** rather than a standard `git subtree add`.

### The Problem with Standard Subtree Merges
Standard subtree merges keep the original tree structure (where files lived at `/`). Consequently, Git only knows about the `<directory>/` prefix in the single graft merge commit, causing `git log <directory>` to return only that merge commit.

### The Path-Rewriting Procedure

For each repository to migrate:

1. **Clone to an isolated temporary workspace**:
   ```bash
   git clone --branch <default-branch> --single-branch <repo-path> <temp-path>
   ```

2. **Rewrite commit trees with `git-filter-repo`**:
   Shift all files into a top-level subdirectory matching the repository name:
   ```bash
   git-filter-repo --to-subdirectory-filter <repo-name>
   ```
   * All historical commits now modify `<repo-name>/...` instead of root paths.
   * Commit messages, author names, email addresses, timestamps, and diffs are strictly preserved.

3. **Add temporary remote in monorepo and fetch**:
   ```bash
   git remote add tmp-remote-<repo-name> <temp-path>
   git fetch tmp-remote-<repo-name>
   ```

4. **Merge unrelated histories**:
   ```bash
   git merge tmp-remote-<repo-name>/<default-branch> \
       --allow-unrelated-histories \
       -m "graft: import <repo-name> with rewritten path history" \
       --no-edit
   ```

5. **Clean up temporary remotes and workspaces**:
   ```bash
   git remote remove tmp-remote-<repo-name>
   ```

---

## 4. Migrated Repositories Summary

| Repository | Source Branch | Destination Directory | Historical Commits |
| :--- | :--- | :--- | :--- |
| `api-linter` | `main` | `api-linter/` | ~1,181 commits |
| `aepc` | `main` | `aepc/` | ~71 commits |
| `aep-lib-go` | `main` | `aep-lib-go/` | ~81 commits |
| `aepcli` | `main` | `aepcli/` | ~50 commits |
| `aep-e2e-validator` | `main` | `aep-e2e-validator/` | ~15 commits |
| `terraform-provider-aep` | `master` | `terraform-provider-aep/` | ~318 commits |

---

## 5. Workspace Integration (`go.work`)

A root `go.work` file was created to link all Go modules for local development:

```go
go 1.26.3

use (
	./aep-e2e-validator
	./aep-lib-go
	./aepc
	./aepcli
	./api-linter
	./terraform-provider-aep
)
```

---

## 6. Verification Commands

To verify that the audit trail and Go toolchain are intact:

- **Check Directory Log**:
  ```bash
  git log aep-lib-go/
  ```
- **Check Specific File Blame / Log**:
  ```bash
  git log aep-lib-go/pkg/api/openapi.go
  git blame aep-lib-go/pkg/api/openapi.go
  ```
- **Run Module Tests**:
  ```bash
  go test ./api-linter/...
  go test ./aepcli/...
  go test ./aep-lib-go/...
  ```
