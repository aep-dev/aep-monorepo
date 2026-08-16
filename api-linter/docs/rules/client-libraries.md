---
permalink: /rules/client-libraries/
---

# Client Library rules

Client library rules are based on [client library-specific AEPs][]. They are
enabled by default, and **should** be enabled for APIs that ship public client
libraries. It is recommended that APIs which do not want to follow certain AEPs
within the client libraries section disable those rules individually, rather
than disabling the client library rules as a set.

{% include linter-group-listing.html start=4200 end=4299 %}

[client library-specific aeps]: https://aep.dev/client-libraries
