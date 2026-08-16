---
rule:
  aep: 134
  name: [core, '0134', request-unknown-fields]
  summary: Update RPCs should not have unexpected fields in the request.
permalink: /134/request-unknown-fields
redirect_from:
  - /0134/request-unknown-fields
---

# Update methods: Unknown fields

This rule enforces that all `Update` standard methods do not have unexpected
fields, as mandated in [AEP-134][].

## Details

This rule looks at any message matching `Update*Request` and complains if it
comes across any fields other than:

- `{Resource} {resource}` ([AEP-134][])
- `bool allow_missing` ([AEP-134][])
- `google.protobuf.FieldMask update_mask` ([AEP-134][])
- `string request_id` ([AEP-155][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message UpdateBookRequest {
  Book book = 1;
  google.protobuf.FieldMask update_mask = 2;
  string library_id = 3;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message UpdateBookRequest {
  Book book = 1;
  google.protobuf.FieldMask update_mask = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message UpdateBookRequest {
  Book book = 1;
  google.protobuf.FieldMask update_mask = 2;
  // (-- api-linter: core::0134::request-unknown-fields=disabled
  //     aep.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 3;  // Non-standard field.
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-134]: https://aep.dev/134
[aep-155]: https://aep.dev/155
[aep.dev/not-precedent]: https://aep.dev/not-precedent
