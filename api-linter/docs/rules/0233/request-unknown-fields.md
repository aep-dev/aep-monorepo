---
rule:
  aep: 233
  name: [core, '0233', request-unknown-fields]
  summary: Batch Create RPCs should not have unexpected fields in the request.
permalink: /233/request-unknown-fields
redirect_from:
  - /0233/request-unknown-fields
---

# Batch Create methods: Unknown request fields

This rule enforces that all `BatchCreate` standard methods do not have unexpected
fields, as mandated in [AEP-233][].

## Details

This rule looks at any message matching `BatchCreate*Request` and complains if it comes
across any fields other than:

- `string parent` ([AEP-233][])
- `string request_id` ([AEP-155][])
- `repeated Create*Request requests` ([AEP-233][])
- `bool validate_only` ([AEP-163][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchCreateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated CreateBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];

  string library_id = 3;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchCreateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated CreateBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message BatchCreateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated CreateBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];

  // (-- api-linter: core::0233::request-unknown-fields=disabled
  //     aep.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-155]: https://aep.dev/155
[aep-163]: https://aep.dev/163
[aep-233]: https://aep.dev/233
[aep.dev/not-precedent]: https://aep.dev/not-precedent
