---
rule:
  aep: 234
  name: [core, '0234', request-unknown-fields]
  summary: Batch Update RPCs should not have unexpected fields in the request.
permalink: /234/request-unknown-fields
redirect_from:
  - /0234/request-unknown-fields
---

# Batch Update methods: Unknown request fields

This rule enforces that all `BatchUpdate` standard methods do not have unexpected
fields, as mandated in [AEP-234][].

## Details

This rule looks at any message matching `BatchUpdate*Request` and complains if it comes
across any fields other than:

- `bool allow_missing` ([AEP-134][])
- `string parent` ([AEP-234][])
- `string request_id` ([AEP-155][])
- `repeated Update*Request requests` ([AEP-234][])
- `google.protobuf.FieldMask update_mask` ([AEP-134][])
- `bool validate_only` ([AEP-163][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated UpdateBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];

  string library_id = 3;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated UpdateBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated UpdateBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];

  // (-- api-linter: core::0234::request-unknown-fields=disabled
  //     aep.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-134]: https://aep.dev/134
[aep-155]: https://aep.dev/155
[aep-163]: https://aep.dev/163
[aep-234]: https://aep.dev/234
[aep.dev/not-precedent]: https://aep.dev/not-precedent
