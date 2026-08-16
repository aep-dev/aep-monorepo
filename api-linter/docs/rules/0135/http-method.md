---
rule:
  aep: 135
  name: [core, '0135', http-method]
  summary: Delete methods must use the DELETE HTTP verb.
permalink: /135/http-method
redirect_from:
  - /0135/http-method
---

# Delete methods: DELETE HTTP verb

This rule enforces that all `Delete` RPCs use the `DELETE` HTTP verb, as
mandated in [AEP-135][].

## Details

This rule looks at any message matching beginning with `Delete`, and complains
if the HTTP verb is anything other than `DELETE`. It _does_ check additional
bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc DeleteBook(DeleteBookRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/books/*}"  // Should be `delete:`.
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc DeleteBook(DeleteBookRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{path=publishers/*/books/*}"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0135::http-method=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc DeleteBook(DeleteBookRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/books/*}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-135]: https://aep.dev/135
[aep.dev/not-precedent]: https://aep.dev/not-precedent
