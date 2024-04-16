---
rule:
  aep: 131
  name: [core, '0131', request-name-field]
  summary: Get RPCs must have a `string name` field in the request.
permalink: /131/request-name-field
redirect_from:
  - /0131/request-name-field
---

# Get methods: Name field

This rule enforces that all `Get` standard methods have a `string name` field
in the request message, as mandated in [AEP-131][].

## Details

This rule looks at any message matching `Get*Request` and complains if
the `name` field type is not `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  bytes name = 1;  // Field type should be `string`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto

message GetBookRequest {
  // (-- api-linter: core::0131::request-name-field=disabled
  //     aep.dev/not-precedent: This uses `bytes` for historical reasons. --)
  bytes name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-131]: https://aep.dev/131
[aep.dev/not-precedent]: https://aep.dev/not-precedent
