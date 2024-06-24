---
rule:
  aep: 135
  name: [core, '0135', request-unknown-fields]
  summary: Delete RPCs should not have unexpected fields in the request.
permalink: /135/request-unknown-fields
redirect_from:
  - /0135/request-unknown-fields
---

# Delete methods: Unknown fields

This rule enforces that all `Delete` standard methods do not have unexpected
fields, as mandated in [AEP-135][].

## Details

This rule looks at any message matching `Delete*Request` and complains if it
comes across any fields other than:

- `string path` ([AEP-135][])
- `bool allow_missing` ([AEP-135][])
- `bool force` ([AEP-135][])
- `string etag` ([AEP-154][])
- `string request_id` ([AEP-155][])
- `bool validate_only` ([AEP-163][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRequest {
  string path = 1;
  string library_id = 2;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message DeleteBookRequest {
  string path = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message DeleteBookRequest {
  string path = 1;

  // (-- api-linter: core::0135::request-unknown-fields=disabled
  //     aep.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-135]: https://aep.dev/135
[aep-154]: https://aep.dev/154
[aep-155]: https://aep.dev/155
[aep-163]: https://aep.dev/163
[aep.dev/not-precedent]: https://aep.dev/not-precedent
