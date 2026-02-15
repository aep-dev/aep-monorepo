# aep-e2e-validator

An AEP validator that ensures compatibility of AEP HTTP APIs end-to-end.

## Goals

aep-e2e-validator covers the gap of validation of runtime functionality defined in [aep.dev](https://aep.dev).

Other tools, such as [aep-openapi-linter](https://github.com/aep-dev/aep-openapi-linter) and [aep-linter](https://github.com/aep-dev/api-linter), perform validation of static interface definitions (e.g. openapi or protobuf files) that define the theoretical behavior of services. However, some behavior cannot be validated based on interface definitions alone (e.g. if a proper status code is used). This is where one can use `aep-e2e-validator`.

## Caution: please run against staging / development APIs

End-to-end validation requires the creation, deletion, list, and so on of the APIs that is being tested. As such, _it is not recommended_ to run this tool against a production API. Instead, it is recommended to run this against a staging API instead, possibly as an automated test in a CI/CD pipeline.

## User Guide

```
go run main.go validate --config "http://localhost:8000/openapi.json" --collection shelves
```
