# Architecture

## The shape

One Go program, `cmd/app`, serving HTTP on port 8080 with the standard
library's `net/http` and its `ServeMux`. It runs in a container built by the
two-stage `Dockerfile`; CI runs `go vet` and `go test -race` and publishes the
image only if they pass.

Today a request goes: `ServeMux` (built by `newMux` in `cmd/app/main.go`) ->
handler -> JSON response. Only the four common routes exist (`/health`,
`/version`, `/`, `POST /reset`) and they are covered by an `httptest` test.
`internal/version` reports the build, and `internal/money` holds the exact
money type, which nothing calls yet.

The intended shape for the rest, following the project description:
`csvimport` turns an export into valid rows plus rejections, `store` holds
transactions, rules and budgets and assigns categories, `report` builds the
monthly totals, and an `api` package wires them to the routes so `main.go`
stays small.

## The data

Everything lives in memory in the running process. Nothing survives a restart,
and `POST /reset` is meant to discard it all. Today there is nothing to discard.

An amount is an `int64` number of cents (`money.Amount`). On the wire it is a
JSON *string* with exactly two decimals, negative for money going out
(`"-12.50"`). It is never a JSON number and never passes through a float, so
`"0.10"` plus `"0.20"` is exactly `"0.30"`. Parsing is strict and reports
overflow instead of wrapping.

Planned, not built: a transaction is an id (integer from 1), a date, an
amount, a description and a category; rules are an ordered list of
`{match, category}`; budgets map a category to a monthly limit.

## Three decisions

1. **Integer cents instead of a decimal library or `float64`.** `float64`
   cannot hold 0.10 exactly and drifts; the project is judged on this. A
   third-party decimal package would work, but cents in an `int64` are enough
   for two-decimal money, need no dependency, and make addition exact. The cost
   is that I write parsing and formatting myself, and the range is limited to
   about 92 quadrillion units, which is why overflow is checked.
2. **Standard library only, with the template's `ServeMux`.** A router such as
   chi or gin was the alternative. Go's `ServeMux` already matches methods and
   paths (`"GET /health"`), the template uses it, and there are only about ten
   routes. Fewer dependencies means a smaller image and less to explain.
3. **In-memory state, no database.** SQLite or a file would survive restarts.
   I rejected that for now because `/reset` must wipe everything, tests need a
   known starting state, and the description does not ask for persistence. The
   price is that data is lost when the container stops.

## What is not built yet

Honestly, most of the project. Built: the four common routes, the build
stamping, and the `money` package with tests.

Not built:

- `POST /transactions` and `GET /transactions` with the filters `from`, `to`,
  `category` and `q`, newest first.
- CSV import (`POST /transactions/import`): the parser, per-line rejections,
  optional header, mixed date formats, decimal commas, quoted fields.
- Duplicate protection on import. The plan is a fingerprint of date, amount and
  normalised description, numbered within a file so two identical coffees both
  import, but I have not built it or tested it.
- Rules (`PUT /rules`): ordered, case-insensitive substring match, first match
  wins, `uncategorised` when none matches.
- Budgets (`PUT /budgets`) and the monthly report
  (`GET /reports/monthly`) with `over_budget`.
- The error shape `{"error","message"}` with 400/404/409/422 for the project
  routes; the helper exists in `main.go` but nothing uses it yet.
- `POST /reset` does not clear anything yet, because no state exists.
- The 12-month cap on date ranges.
- Tests for all of the above. Only `money` and the four common routes have
  real tests so far.
- The final report, `docs/report.docx`.
