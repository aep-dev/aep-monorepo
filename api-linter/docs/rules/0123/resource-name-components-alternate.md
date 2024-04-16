---
rule:
  aep: 123
  name: [core, '0123', resource-name-components-alternate]
  summary: |
    Resource name components should alternate between collection and
    identifiers.
permalink: /123/resource-name-components-alternate
redirect_from:
  - /0123/resource-name-components-alternate
---

# Resource name components alternate

This rule enforces that messages that have a `google.api.resource` annotation
have `pattern` annotations that alternate between collection and identifier, as
described in [AEP-123][].

## Details

This rule scans messages with a `google.api.resource` annotation, and validates
that each `pattern` alternated between collection and identifiers.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    // two collections next to each other.
    pattern: "publishers/books/{book}"
  };
  string name = 1;
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

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::0123::resource-name-components-alternate=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/books/{book}"
  };
  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-123]: http://aep.dev/123
[aep.dev/not-precedent]: https://aep.dev/not-precedent
