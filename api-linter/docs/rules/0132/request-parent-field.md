---
rule:
  aep: 132
  name: [core, '0132', request-parent-field]
  summary: List RPCs must have a `parent` field in the request.
permalink: /132/request-parent-field
redirect_from:
  - /0132/request-parent-field
---

# List methods: Parent field

This rule enforces that all `List` standard methods have a `string parent`
field in the request message, as mandated in [AEP-132][].

## Details

This rule looks at any message matching `List*Request` and complains if the 
`parent` field has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  // Field type should be `string`.
  bytes parent = 1;
  int32 page_size = 2;
  string page_token = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message ListBooksRequest {
  // (-- api-linter: core::0132::request-parent-field=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  bytes parent = 1;
  int32 page_size = 2;
  string page_token = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-132]: https://aep.dev/132
[aep.dev/not-precedent]: https://aep.dev/not-precedent
