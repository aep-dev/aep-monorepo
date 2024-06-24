---
rule:
  aep: 133
  name: [core, '0133', request-unknown-fields]
  summary: Create RPCs should not have unexpected fields in the request.
permalink: /133/request-unknown-fields
redirect_from:
  - /0133/request-unknown-fields
---

# Create methods: Unknown fields

This rule enforces that all `Create` standard methods do not have unexpected
fields, as mandated in [AEP-133][].

## Details

This rule looks at any message matching `Create*Request` and complains if it
comes across any fields other than:

- `string parent` ([AEP-133][])
- `{Resource} {resource}` ([AEP-133][])
- `string {resource}_id` ([AEP-133][])
- `string request_id` ([AEP-155][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string book_id = 3;
  string library_id = 4;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string id = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string id = 3;

  // (-- api-linter: core::0133::request-unknown-fields=disabled
  //     aep.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 4;  // Non-standard field.
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-133]: https://aep.dev/133
[aep-155]: https://aep.dev/155
[aep.dev/not-precedent]: https://aep.dev/not-precedent
