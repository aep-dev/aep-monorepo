---
rule:
  aep: 203
  name: [core, '0203', resource-name-identifier]
  summary: Resource name field must have field behavior IDENTIFIER.
permalink: /203/resource-name-identifier
redirect_from:
  - /0203/resource-name-identifier
---

# Resource name: IDENTIFIER

This rule enforces that the field on a resource representing the resource's name
is annotated with with `(google.api.field_behavior) = IDENTIFIER`, as mandated
by [AEP-203][].

## Details

This rule looks at every resource message and complains if the resource name
field is not annotated with  `(google.api.field_behavior) = IDENTIFER`.

## Examples

**Incorrect** code for this rule:

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
  // Missing IDENTIFIER field behavior.
  string name = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
  string name = 1 [(google.api.field_behavior) = IDENTIFIER];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
  // (-- api-linter: core::0203::resource-name-identifier=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-203]: https://aep.dev/203
[aep.dev/not-precedent]: https://aep.dev/not-precedent
