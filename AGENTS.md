# Agents Guidelines

- Before making changes to this repository, read the contribution and development guidelines in [CONTRIBUTING.md](CONTRIBUTING.md).
- Follow the repository conventions, Bazel build patterns, Go workspace patterns, and dependency synchronization workflows outlined in [CONTRIBUTING.md](CONTRIBUTING.md).
- All changes must be made via a pull request rather than committing directly to the main branch.
- Always run all tests locally using `bazel test //...` (or `just test`) and verify checks pass with `bazel run //:check` (or `just check`) before pushing changes or creating/updating a pull request.
- All changes must pass tests across all modules using `bazel test //...` (or `just test`).
- If modifying shared libraries like [aep-lib-go](aep-lib-go), run `bazel run //:fix` (or `just fix`) to ensure dependent modules and Bazel files are synchronized, and verify with `bazel run //:check` (or `just check`).
- If adding or changing Go files/packages, run `bazel run //:gazelle` (or `just gazelle`) to regenerate/update Bazel BUILD definitions.

## References

- Project contribution and development guide: [CONTRIBUTING.md](CONTRIBUTING.md)
- Architecture and repository initialization design document: [DESIGN/2026-08-16-repository-initialization-and-grafting.md](DESIGN/2026-08-16-repository-initialization-and-grafting.md)
- Bazel build migration design document: [DESIGN/2026-09-01-bazel-migration.md](DESIGN/2026-09-01-bazel-migration.md)
- Root Bazel module configuration: [MODULE.bazel](MODULE.bazel)
- Command runner recipe definitions: [justfile](justfile)
- Dependency synchronization helper script: [scripts/sync-deps.py](scripts/sync-deps.py)
