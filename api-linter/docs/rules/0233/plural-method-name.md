---
rule:
  aep: 233
  name: [core, '0233', plural-method-name]
  summary: Batch Create methods must have plural names.
permalink: /233/plural-method-name
redirect_from:
  - /0233/plural-method-name
---

# Batch Create methods: Plural method name

This rule enforces that all `BatchCreate` RPCs have a plural resource in the
method name, as mandated in [AEP-233][].

## Details

This rule looks at any method whose name begins with `BatchCreate`, and complains
if the name of the resource in the method name is singular.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Method name should be `BatchCreateBooks`
rpc BatchCreateBook(BatchCreateBookRequest) returns (BatchCreateBookResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchCreate"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BatchCreateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchCreate"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0233::plural-method-name=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc BatchCreateBook(BatchCreateBookRequest) returns (BatchCreateBookResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchCreate"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-233]: https://aep.dev/233
[aep.dev/not-precedent]: https://aep.dev/not-precedent
