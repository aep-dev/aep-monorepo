# Justfile for AEP monorepo tasks

export USE_BAZEL_VERSION := "7.4.1"

# Default recipe: list available recipes
default:
    @just --list

# Run tests across all workspace targets using Bazel
test:
    bazel test //...

# Build all workspace targets using Bazel
build:
    bazel build //...

# Run Gazelle to generate/update Bazel BUILD files
gazelle:
    bazel run //:gazelle

# Synchronize internal monorepo dependencies and update BUILD files
fix:
    bazel run //:fix
    bazel run //:gazelle

# Verify internal monorepo dependencies are synchronized
check:
    bazel run //:check
