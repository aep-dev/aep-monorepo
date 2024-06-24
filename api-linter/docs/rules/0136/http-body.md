---
rule:
  aep: 136
  name: [core, '0136', http-body]
  summary: Custom methods must have the HTTP body set to `*`.
permalink: /136/http-body
redirect_from:
  - /0136/http-body
---

# Custom methods: HTTP body

This rule enforces that all custom methods set the HTTP `body` to `*`, as
advised in [AEP-136][].

## Details

This rule looks at any method that is not a standard method, and complains if
the HTTP `body` field is not set to `"*"`. It also permits the path of the
resource.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/books}:checkout"
    // `body: "*"` should be included.
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/books}:checkout"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0136::http-body=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/books}:checkout"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-136]: https://aep.dev/136
[aep.dev/not-precedent]: https://aep.dev/not-precedent
