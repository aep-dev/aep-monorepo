---
rule:
  aep: 123
  name: [core, '0123', resource-annotation]
  summary: Resource messages should be annotated with `google.api.resource`.
permalink: /123/resource-annotation
redirect_from:
  - /0123/resource-annotation
---

# Resource annotation presence

This rule enforces that top-level messages that appear to represent resources
have a `google.api.resource` annotation, as described in [AEP-123][].

## Details

This rule scans all top-level messages, and assumes that messages with a
 `string name` field are resources unless the message name ends with `Request`.
For messages that are resources, it complains if the `google.api.resource`
annotation is missing.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  // A `google.api.resource` annotation should be here.
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
// (-- api-linter: core::0123::resource-annotation=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-123]: http://aep.dev/123
[aep.dev/not-precedent]: https://aep.dev/not-precedent
