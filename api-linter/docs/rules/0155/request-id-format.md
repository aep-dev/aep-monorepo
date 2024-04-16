---
rule:
  aep: 155
  name: [core, '0155', request-id-format]
  summary: Annotate request_id with UUID4 format.
permalink: /155/request-id-format
redirect_from:
  - /0155/request-id-format
---

# `request_id` format annotation

This rule encourages the use of the `UUID4` format annotation on the
`request_id` field, as mandated in [AEP-155][].

## Details

This rule looks on for fields named `request_id` and complains if it does not
have the `(google.api.field_info).format = UUID4` annotation or has a format
other than `UUID4`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBookRequest {
  string parent = 1;

  Book book = 2;

  string request_id = 3; // missing (google.api.field_info).format = UUID4
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBookRequest {
  string parent = 1;

  Book book = 2;

  string request_id = 3 [(google.api.field_info).format = UUID4];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field or its
enclosing message. Remember to also include an [aep.dev/not-precedent][]
comment explaining why.

```proto
message CreateBookRequest {
  string parent = 1;

  Book book = 2;

  // (-- api-linter: core::0155::request-id-format=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string request_id = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-155]: https://aep.dev/155
[aep.dev/not-precedent]: https://aep.dev/not-precedent