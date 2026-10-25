# Triage Labels

The canonical triage state roles use these strings:

- `needs-triage`
- `needs-info`
- `ready-for-agent`
- `ready-for-human`
- `wontfix`

For the local Markdown tracker, record `bug` or `enhancement` in `Category:` and the applicable state role in `Triage:`. Reserve `Status:` for the independent work lifecycle (`open`, `claimed`, `resolved`, or a non-fulfilling closure). Ordinary triage transitions never change `Status:`; closing a ticket is the exception: done work sets `Status: resolved`, and a `wontfix` closure sets `Status: closed (wontfix)`, `closed (duplicate)`, or `closed (out-of-scope)`, which never unblocks a dependent ticket.
