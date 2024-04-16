---
rule:
  aep: 203
  name: [core, '0203', unordered-list-repeated]
  summary: Only repeated fields may be annotated with `UNORDERED_LIST`.
permalink: /203/unordered-list-repeated
redirect_from:
  - /0203/unordered-list-repeated
---

# Repeated fields: Unordered

This rule enforces that only repeated fields, not singular ones, are annotated
with `(google.api.field_behavior) = UNORDERED_LIST`, as mandated by [AEP-203][].

## Details

This rule looks at any field with a `(google.api.field_behavior) =
UNORDERED_LIST` annotation and complains if it is not a repeated field.

## Examples

**Incorrect** code for this rule:

```proto
message Book {
  // UNORDERED_LIST only applies to repeated fields.
  string name = 1 [(google.api.field_behavior) = UNORDERED_LIST];

  repeated string genres = 2 [(google.api.field_behavior) = UNORDERED_LIST];
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;

  repeated string genres = 2 [(google.api.field_behavior) = UNORDERED_LIST];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message Book {
  // (-- api-linter: core::0203::unordered-list-repeated=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [(google.api.field_behavior) = UNORDERED_LIST];

  repeated string genres = 2 [(google.api.field_behavior) = UNORDERED_LIST];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-203]: https://aep.dev/203
[aep.dev/not-precedent]: https://aep.dev/not-precedent
