---
rule:
  aep: 123
  name: [core, '0123', resource-name-field]
  summary: Resource messages should have a `string name` field.
permalink: /123/resource-name-field
redirect_from:
  - /0123/resource-name-field
---

# Resource `name` field

This rule enforces that messages that appear to represent resources have a
`string name` field, as described in [AEP-123][].

## Details

This rule scans all messages that have a `google.api.resource` annotation, and
complains if the `name` field is missing or if it is any type other than
singular `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect: missing `string name` field.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };
}
```

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  // Should be `string`, not `bytes`.
  bytes name = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message, or
above the field if it is the wrong type.

```proto
// (-- api-linter: core::0123::resource-name-field=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-123]: http://aep.dev/123
[aep.dev/not-precedent]: https://aep.dev/not-precedent
