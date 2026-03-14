# ADR-0001: Balance Maintenance
Date: 2026-03-14
Status: Accepted

## Context
Accounts need a `current_balance`. When transactions are added/edited/deleted, it must stay in sync.

## Decision
We will use application logic wrapped in GORM transactions with pessimistic locking (`SELECT * FROM accounts FOR UPDATE`).

## Rationale
- Easy to unit test in Go.
- High performance (reads are fast, writes are safely serialized).

## Consequences
- Requires strict adherence to using the transaction wrapper in the Service layer.

## Alternatives Considered
- DB triggers (rejected: hidden magic, hard to test).
- Computed views (rejected: performance concerns).
