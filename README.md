# Personal Finance Tracker

A small web service for tracking what you earn and spend. Transactions come in
from a bank CSV export or are added by hand, get sorted into categories by rules
the user writes, and come back out as a monthly report. It is meant for one
person keeping an eye on their own money.

- **Student:** Nistor Carina, group DMWD
- **Project:** 2 - Personal Finance Tracker
- **Language:** Go

## What it does

The service has no idea what a transaction means, so categorisation is an
ordered list of rules the user states and can change: a rule is a piece of text
to look for in the description and the category to assign. The first matching
rule wins, and the service never reorders them.

Money is never a floating-point number. Every amount is held as a whole number
of cents, so `0.10` plus `0.20` is exactly `0.30`.

**State (milestone 1):** the service starts, answers the four common routes
(`/health`, `/version`, `/`, `POST /reset`) and the exact-money type is written
and tested. The project routes (transactions, CSV import, rules, budgets,
reports) are **not built yet**. `docs/architecture.md` lists exactly what is
missing.

## How to run it

```
docker run --rm -p 8080:8080 ghcr.io/carinanistor/personal-finance-tracker:latest
```

To build and run without Docker:

```
go run ./cmd/app
go run ./cmd/app --version
```

To run the tests:

```
go vet ./...
go test ./... -race
```

## What you should see

```
$ curl -i http://localhost:8080/health
HTTP/1.1 200 OK
{"status":"ok","uptime_seconds":3}
```

`curl http://localhost:8080/version` returns `app`, `version`, `commit` and
`built_at`, and `POST /reset` returns `204`.
