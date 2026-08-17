# Justfile for AEP monorepo tasks

# Default recipe: list available recipes
default:
    @just --list

# Run tests across all workspace modules
test:
    go test ./aep-lib-go/... ./aepcli/... ./api-linter/... ./aep-e2e-validator/... ./aepc/example/service/... ./aepc/pkg/...

# Synchronize internal monorepo dependencies
fix:
    ./scripts/sync-deps.py

# Verify internal monorepo dependencies are synchronized
check:
    ./scripts/sync-deps.py --check
