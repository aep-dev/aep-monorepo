---
rule:
  aep: 134
  name: [core, '0134', response-message-name]
  summary: Update methods must return the resource.
permalink: /134/response-message-name
redirect_from:
  - /0134/response-message-name
---

# Update methods: Resource response message

This rule enforces that all `Update` RPCs have a response message of the
resource, as mandated in [AEP-134][].

## Details

This rule looks at any message matching beginning with `Update`, and complains
if the name of the corresponding output message does not match the name of the
RPC with the prefix `Update` removed.

It also permits a response of `google.longrunning.Operation`; in this case, it
checks the `response_type` in the `google.longrunning.operation_info`
annotation and ensures that _it_ corresponds to the name of the RPC with the
prefix `Update` removed.

## Examples

### Standard

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `Book`.
rpc UpdateBook(UpdateBookRequest) returns (UpdateBookResponse) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
  };
}
```

### Long-running operation

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc UpdateBook(UpdateBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
  };
  option (google.longrunning.operation_info) = {
    response_type: "UpdateBookResponse"  // Should be "Book".
    metadata_type: "UpdateBookMetadata"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc UpdateBook(UpdateBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
  };
  option (google.longrunning.operation_info) = {
    response_type: "Book"
    metadata_type: "UpdateBookMetadata"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0134::response-message-name=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc UpdateBook(UpdateBookRequest) returns (UpdateBookResponse) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-134]: https://aep.dev/134
[aep.dev/not-precedent]: https://aep.dev/not-precedent
