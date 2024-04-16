---
rule:
  aep: 165
  name: [core, '0165', http-body]
  summary: Purge methods should use `*` as the HTTP body.
permalink: /165/http-body
redirect_from:
  - /0165/http-body
---

# Purge methods: HTTP body

This rule enforces that all `Purge` RPCs use `*` as the HTTP `body`, as mandated in
[AEP-165][].

## Details

This rule looks at any message beginning with `Purge`, and complains
if the HTTP `body` field is anything other than `*`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc PurgeBooks(PurgeBooksRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:purge"
    // body: "*" should be set.
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
// (-- api-linter: core::0165::http-body=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc PurgeBooks(PurgeBooksRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:purge"
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
