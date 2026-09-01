# History Grafting and Integration of `aep-lib-ts`

* **Date**: 2026-09-01
* **Status**: Completed

## 1. Overview

This document records the grafting of [`aep-lib-ts`](https://github.com/aep-dev/aep-lib-ts) into `aep-monorepo` while preserving all 17 historical Git commits and line-by-line attribution.

---

## 2. Source Repository Metadata

* **Source Directory**: `/Users/TZTWH7/workspace/aep-lib-ts`
* **Target Monorepo Subdirectory**: `aep-lib-ts/`
* **Historical Commits Imported**: 17 commits (`328cb24` through `7586d78`)
* **Package Name**: `@aep_dev/aep-lib-ts`
* **Language/Framework**: TypeScript 5.3+, Node.js 20+, Jest / ts-jest

---

## 3. Graft Methodology

The graft was executed using `git-filter-repo` to rewrite commit tree paths without altering commit timestamps, authors, or messages:

1. Clone source repository:
   ```bash
   git clone /Users/TZTWH7/workspace/aep-lib-ts /tmp/aep-lib-ts-graft
   cd /tmp/aep-lib-ts-graft
   ```
2. Rewrite subdirectory paths:
   ```bash
   git-filter-repo --to-subdirectory-filter aep-lib-ts
   ```
3. Add temporary remote and merge:
   ```bash
   git remote add tmp-graft-ts /tmp/aep-lib-ts-graft
   git fetch tmp-graft-ts
   git merge tmp-graft-ts/main --allow-unrelated-histories -m "graft: import aep-lib-ts with rewritten path history" --no-edit
   git remote remove tmp-graft-ts
   ```

---

## 4. Bazel Target Architecture

`aep-lib-ts` is integrated into Bazel:
* `//aep-lib-ts:test`: Runs the Jest unit test suite via `sh_test` (35 unit tests across 4 test suites).
* `//aep-lib-ts:build`: Runs the TypeScript compiler (`tsc`).
* `//aep-lib-ts:srcs`: Exports source and configuration filegroups.

---

## 5. CI/CD Workflows

* **[`.github/workflows/aep-lib-ts-ci.yml`](.github/workflows/aep-lib-ts-ci.yml)**: Runs `bazel test //aep-lib-ts/...` on PRs modifying TypeScript code.
* **[`.github/workflows/aep-lib-ts-release.yml`](.github/workflows/aep-lib-ts-release.yml)**: Builds and publishes `@aep_dev/aep-lib-ts` to npm upon pushing `aep-lib-ts/v*` tags.
