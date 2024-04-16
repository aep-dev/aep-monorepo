---
rule:
  aep: 158
  name: [core, '0158', response-plural-first-field]
  summary: First field of Paginated RPCs' response should be plural.
permalink: /158/response-plural-first-field
redirect_from:
  - /0158/response-plural-first-field
---

# Paginated methods: Page token field

This rule enforces that all `List` and `Search` methods have a plural name
repeatable field as a first field in the response message, as mandated in
[AEP-158][].

## Details

This rule looks at any message matching `List*Response` or `Search*Response`
that has `next_page_token` field and complains if the first field's name is not
plural.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect
message ListBooksResponse {
  // Field name should be `books`.
  repeated Book book = 1;
  string next_page_token = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message or
above the field. Remember to also include an [aep.dev/not-precedent][] comment
explaining why.

```proto
// (-- api-linter: core::0158::response-plural-first-field=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message ListBooksResponse {
  repeated Book book = 1;
  string next_page_token = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-158]: https://aep.dev/158
[aep.dev/not-precedent]: https://aep.dev/not-precedent
