---
rule:
  aep: 165
  name: [core, '0165', http-uri-suffix]
  summary: Purge methods must have the correct URI suffix
permalink: /165/http-uri-suffix
redirect_from:
  - /0165/http-uri-suffix
---

# Purge methods: URI suffix

This rule enforces that `Purge` methods include the `:purge` suffix
in the REST URI, as mandated in [AEP-165][].

## Details

This rule looks at any method whose name starts with `Purge`, and
complains if the HTTP URI does not end with `:purge`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc PurgeBooks(PurgeBooksRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:delete"  // Should end with `:purge`
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "PurgeBooksResponse"
    metadata_type: "PurgeBooksMetadata"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc PurgeBooks(PurgeBooksRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:purge"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "PurgeBooksResponse"
    metadata_type: "PurgeBooksMetadata"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0165::http-uri-suffix=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc PurgeBooks(PurgeBooksRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:delete"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "PurgeBooksResponse"
    metadata_type: "PurgeBooksMetadata"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-165]: https://aep.dev/165
[aep.dev/not-precedent]: https://aep.dev/not-precedent
