---
rule:
  aep: 231
  name: [core, '0231', request-unknown-fields]
  summary: Batch Get RPCs should not have unexpected fields in the request.
permalink: /231/request-unknown-fields
redirect_from:
  - /0231/request-unknown-fields
---

# Batch Get methods: Unknown request fields

This rule enforces that all `BatchGet` standard methods do not have unexpected
fields, as mandated in [AEP-231][].

## Details

This rule looks at any message matching `BatchGet*Request` and complains if it comes
across any fields other than:

- `repeated string names` ([AEP-231][])
- `string parent` ([AEP-231][])
- `google.protobuf.FieldMask read_mask` ([AEP-157][])
- `repeated Get*Request requests` ([AEP-231][])
- `View view` ([AEP-157][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchGetBooksRequest {
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
message BatchGetBooksRequest {
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
message BatchGetBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated string names = 2 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  // (-- api-linter: core::0231::request-unknown-fields=disabled
  //     aep.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-157]: https://aep.dev/157
[aep-231]: https://aep.dev/231
[aep.dev/not-precedent]: https://aep.dev/not-precedent
