# Agents Guidelines

- Before making changes to this repository, read the contribution and development guidelines in [CONTRIBUTING.md](CONTRIBUTING.md).
- Follow the repository conventions, Go workspace patterns, and dependency synchronization workflows outlined in [CONTRIBUTING.md](CONTRIBUTING.md).
- All changes must be made via a pull request rather than committing directly to the main branch.
- All changes must pass tests across all modules using `go test ./...`.
- If modifying shared libraries like [aep-lib-go](aep-lib-go), run [scripts/sync-deps.py](scripts/sync-deps.py) to ensure dependent modules are synchronized.

## References

- Project contribution and development guide: [CONTRIBUTING.md](CONTRIBUTING.md)
- Architecture and repository initialization design document: [DESIGN/2026-08-16-repository-initialization-and-grafting.md](DESIGN/2026-08-16-repository-initialization-and-grafting.md)
- Root Go workspace configuration: [go.work](go.work)
- Dependency synchronization helper script: [scripts/sync-deps.py](scripts/sync-deps.py)
