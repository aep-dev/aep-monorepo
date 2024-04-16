---
rule:
  aep: 134
  name: [core, '0134', http-body]
  summary: Update methods must have the HTTP body set to the resource.
permalink: /134/http-body
redirect_from:
  - /0134/http-body
---

# Update methods: HTTP body

This rule enforces that all `Update` RPCs set the HTTP `body` to the resource,
as mandated in [AEP-134][].

## Details

This rule looks at any method whose name begins with `Update`, and complains if
the HTTP `body` field is not set to the resource being created.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "*"  // This should be "book".
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

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0134::http-body=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-134]: https://aep.dev/134
[aep.dev/not-precedent]: https://aep.dev/not-precedent
