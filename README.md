# AEP Monorepo

- This repository contains the consolidated source code and tooling ecosystem for [API Enhancement Proposals (AEP)](https://aep.dev).
- It is managed as a multi-module Go workspace linking shared libraries, command-line tools, linters, code generators, and infrastructure providers.

## Projects

- [aep-e2e-validator](aep-e2e-validator): Runtime end-to-end conformance testing tool for AEP HTTP APIs.
- [aep-lib-go](aep-lib-go): Shared Go core libraries, common types, OpenAPI parsers, and resource schema representations.
- [aepc](aepc): Resource-oriented compiler and code generator for AEP services, protocol buffers, and OpenAPI specifications.
- [aepcli](aepcli): Dynamic command-line interface tool generated from OpenAPI definitions for AEP APIs.
- [api-linter](api-linter): Protocol buffer linter enforcing AEP design rules and naming standards.

## Getting Started

- Clone the repository:
  ```bash
  git clone https://github.com/aep-dev/aep-monorepo.git
  cd aep-monorepo
  ```
- The repository uses Go workspaces defined in [go.work](go.work).
- Run tests across all workspace modules:
  ```bash
  go test ./...
  ```
- Run tests for an individual module:
  ```bash
  go test ./aepcli/...
  ```

## Documentation

- Guidelines for contributing, local development, module tagging, and pull requests: [CONTRIBUTING.md](CONTRIBUTING.md)

## References

- API Enhancement Proposals standard documentation and specifications: https://aep.dev
- Root Go workspace configuration linking all submodules: [go.work](go.work)
- Project contribution and development guide: [CONTRIBUTING.md](CONTRIBUTING.md)
- GitHub Actions CI/CD workflows directory: [.github/workflows](.github/workflows)
