---
rule:
  aep: 154
  name: [core, '0154', field-type]
  summary: Etag fields must be strings.
permalink: /154/field-type
redirect_from:
  - /0154/field-type
---

# Etag field type

This rule enforces that `etag` fields are strings, as mandated in [AEP-154][].

## Details

This rule looks at any field named `etag` and complains if the type is anything
other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRequest {
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
  bytes etag = 2;  // Should be string.
}
```

**Correct** code for this rule:

```proto
// Correct.
message DeleteBookRequest {
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
  string etag = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message DeleteBookRequest {
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
  // (-- api-linter: core::0154::field-type=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  bytes etag = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-154]: https://aep.dev/154
[aep.dev/not-precedent]: https://aep.dev/not-precedent
