---
rule:
  aep: 132
  name: [core, '0132', request-message-name]
  summary: List methods must have standardized request message names.
permalink: /132/request-message-name
redirect_from:
  - /0132/request-message-name
---

# List methods: Request message

This rule enforces that all `List` RPCs have a request message name of
`List*Request`, as mandated in [AEP-132][].

## Details

This rule looks at any message matching beginning with `List`, and complains if
the name of the corresponding input message does not match the name of the RPC
with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `ListBooksRequest`.
rpc ListBooks(ListBooksReq) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0132::request-message-name=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc ListBooks(ListBooksReq) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-132]: https://aep.dev/132
[aep.dev/not-precedent]: https://aep.dev/not-precedent
