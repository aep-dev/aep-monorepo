---
rule:
  aep: 235
  name: [core, '0235', request-unknown-fields]
  summary: Batch Delete RPCs should not have unexpected fields in the request.
permalink: /235/request-unknown-fields
redirect_from:
  - /0235/request-unknown-fields
---

# Batch Delete methods: Unknown request fields

This rule enforces that all `BatchDelete` standard methods do not have unexpected
fields, as mandated in [AEP-235][].

## Details

This rule looks at any message matching `BatchDelete*Request` and complains if it comes
across any fields other than:

- `bool force` ([AEP-135][])
- `repeated string names` ([AEP-235][])
- `string parent` ([AEP-235][])
- `string request_id` ([AEP-155][])
- `repeated Delete*Request requests` ([AEP-235][])
- `bool validate_only` ([AEP-163][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchDeleteBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated string names = 2 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string library_id = 3;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchDeleteBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated string names = 2 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message BatchDeleteBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated string names = 2 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  // (-- api-linter: core::0235::request-unknown-fields=disabled
  //     aep.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-135]: https://aep.dev/135
[aep-155]: https://aep.dev/155
[aep-163]: https://aep.dev/163
[aep-235]: https://aep.dev/235
[aep.dev/not-precedent]: https://aep.dev/not-precedent
