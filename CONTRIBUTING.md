# Contributing to ledgerA

Thank you for your interest in contributing! Please read this guide before opening a pull request.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Development Setup](#development-setup)
3. [Branch Naming](#branch-naming)
4. [Commit Style](#commit-style)
5. [Code Style](#code-style)
6. [Running Tests](#running-tests)
7. [Adding a Database Migration](#adding-a-database-migration)
8. [Pull Request Process](#pull-request-process)

---

## Prerequisites

| Tool | Version |
|------|---------|
| Go | 1.22+ |
| Node.js | 20+ |
| PostgreSQL | 15+ |
| Docker + Docker Compose | any recent |
| golangci-lint | 1.58+ |

---

## Development Setup

```bash
# 1. Clone
git clone https://github.com/<your-org>/ledgerA.git
cd ledgerA

# 2. Copy env and edit
cp .env.example .env
# Set DATABASE_URL, FIREBASE_CREDENTIALS_PATH, CORS_ALLOWED_ORIGINS

# 3. Start Postgres
docker compose up -d db

# 4. Apply migrations
make migrate-up

# 5. Run backend
make dev            # hot-reload via 'go run ./cmd/server'

# 6. Run frontend (separate terminal)
cd frontend && npm install && npm run dev
```

---

## Branch Naming

| Type | Pattern | Example |
|------|---------|---------|
| Feature | `feat/<short-description>` | `feat/audit-export` |
| Bug fix | `fix/<short-description>` | `fix/stats-nil-pointer` |
| Docs | `docs/<short-description>` | `docs/deployment-guide` |
| Chore | `chore/<short-description>` | `chore/upgrade-gin` |

Branch off `main`. Keep branches short-lived; rebase before merging.

---

## Commit Style

ledgerA follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/):

```
<type>(<scope>): <short summary>

[optional body]

[optional footer: BREAKING CHANGE or Closes #issue]
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`

**Scope examples:** `handler`, `service`, `repo`, `frontend`, `migration`, `config`

Examples:

```
feat(handler): add audit log list endpoint
fix(service): handle nil subcategory on transaction create
docs(readme): update local dev steps
chore(deps): upgrade firebase-admin to v6
```

---

## Code Style

### Go

- Run `gofmt -s -w ./cmd ./internal ./pkg` before committing.
- Linting: `golangci-lint run ./cmd/... ./internal/... ./pkg/...` must produce zero warnings.
- Enabled linters (see `.golangci.yml`): `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`, `gofmt`, `revive`, `gocyclo`, `gosec`.
- Always handle returned errors — do not use `_` to discard errors from functions that can fail.
- Keep cyclomatic complexity below 15 per function (enforced by `gocyclo`).
- New packages must be placed under `internal/` (not exported).

### TypeScript / React

- Strict TypeScript — `tsc --noEmit` must pass with zero errors.
- ESLint: `npx eslint src/` must produce zero errors.
- Prefer named exports over default exports.
- Components live in `src/components/`, pages in `src/pages/`, API hooks in `src/hooks/`.
- Use `@/` path alias (mapped to `src/` in `tsconfig.json`).
- Do not commit any `console.log` statements.

---

## Running Tests

```bash
# Backend unit tests (with race detector)
make test
# Equivalent: go test ./... -race -count=1

# Frontend type check + lint
cd frontend
npx tsc --noEmit
npx eslint src/

# Full production build check
npm run build
```

There are currently no integration tests. When adding them, place them under `internal/<package>/<package>_integration_test.go` and gate with a `//go:build integration` build tag.

---

## Adding a Database Migration

```bash
# CREATE the migration file pair
make migrate-create NAME=add_tags_to_transactions
# Creates: migrations/<timestamp>_add_tags_to_transactions.up.sql
#          migrations/<timestamp>_add_tags_to_transactions.down.sql
```

Rules:

1. **Never edit a migration that has already been committed to `main`.** Create a new migration instead.
2. Every `.up.sql` must have a corresponding `.down.sql` that fully reverts the change.
3. Use `numeric(20,4)` for monetary columns, `uuid` for IDs, `timestamptz` for timestamps.
4. Add indexes for any column used in `WHERE` or `ORDER BY` clauses.
5. Test both directions locally: `make migrate-up` then `make migrate-down`.

---

## Pull Request Process

1. Ensure all gates pass locally:
   ```bash
   make build test lint    # backend
   cd frontend && npx tsc --noEmit && npx eslint src/ && npm run build
   ```
2. Fill out the PR template (description, linked issue, testing steps).
3. Squash-merge is preferred. Keep the commit message in Conventional Commits format.
4. At least one approval is required before merging.
5. Delete the branch after merge.

---

## Questions?

Open a GitHub Discussion or drop a note in the PR — we're happy to help.
