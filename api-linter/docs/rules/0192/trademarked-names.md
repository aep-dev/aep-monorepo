---
rule:
  aep: 192
  name: [core, '0192', trademarked-names]
  summary: Trademarked names should be used correctly.
permalink: /192/trademarked-names
redirect_from:
  - /0192/trademarked-names
---

# Trademarked names

This rule enforces trademarked names in public comments are not abbreviated and
follow the trademark owner's branding style, as mandated in [AEP-192][].

## Details

This rule looks at each descriptor in each proto file (exempting oneofs and the
file itself) and complains if it catches a common trademark mistake. It only
checks against a known problem list, and does not make guesses beyond that.

It currently catches common mistaken variants for:

- App Engine
- BigQuery
- Bigtable
- Bitbucket
- Cloud Storage
- Compute Engine
- Dataflow
- Dataprep
- Dialogflow
- Directory Sync
- GitHub
- GitLab
- G Suite
- Pub/Sub
- Service Mesh
- Stack Overflow

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;

  // A repository containing Markdown files for each chapter of
  // the book on Github.
  // (--           ^ Should be GitHub. --)
  string github_repo = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;

  // A repository containing Markdown files for each chapter of
  // the book on GitHub.
  string github_repo = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the descriptor.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;

  // (-- api-linter: core::0192::trademarked-names=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  // A repository containing Markdown files for each chapter of
  // the book on Github.
  string github_repo = 2;
}
```

[aep-192]: https://aep.dev/192
[aep.dev/not-precedent]: https://aep.dev/not-precedent
