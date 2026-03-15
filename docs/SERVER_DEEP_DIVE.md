# Server Deep Dive and Operations Handbook

## 1. Purpose of This Document

This document explains the backend server in depth for maintainers and contributors.
It covers:

- How each server component works in this codebase
- Why the selected tech stack is useful for this system
- Design decisions and trade-offs
- End-to-end request flow and data flow
- How to run, validate, troubleshoot, and operate the server

This is intentionally detailed and implementation-oriented, not a high-level marketing overview.

---

## 2. System Overview

`ledgerA` backend is a layered Go application with the following structure:

- Entrypoint and dependency wiring: `cmd/server/main.go`
- Config loading: `internal/config`
- HTTP and routing: `internal/handler`
- Middleware: `internal/middleware`
- Business logic: `internal/service`
- Persistence and query logic: `internal/repository`
- Domain models: `internal/model`
- Infrastructure packages:
  - Firebase Admin integration: `pkg/firebase`
  - PDF generation support: `pkg/pdf`
- SQL schema and evolution: `migrations`

High-level data flow:

1. HTTP request enters Gin router.
2. Middleware runs (logging, panic recovery, auth where required).
3. Handler validates/parses request and resolves authenticated user.
4. Service executes business rules and orchestration.
5. Repository performs DB reads/writes via GORM/PostgreSQL.
6. Handler emits standardized API envelope (`data` or `error`).

---

## 3. Tech Stack and Why It Helps

### 3.1 Go

Benefits in this backend:

- Simple concurrency model when needed (`goroutines`, context propagation).
- Strong static typing across DTO/service/repository boundaries.
- Fast build/test cycle and low runtime overhead.
- Good operational ergonomics for single-binary deployment.

### 3.2 Gin

Why Gin is a good fit here:

- Lightweight, fast HTTP router.
- Explicit middleware chain keeps auth/logging/recovery separated.
- Easy binding for query/body structs (`ShouldBindJSON`, `ShouldBindQuery`).
- Clear route grouping (`/api/v1`, protected subgroup).

### 3.3 GORM + PostgreSQL

Why this combination works for ledgerA:

- PostgreSQL provides reliable transactions, constraints, indexes, and JSONB for audit diffs.
- GORM reduces repetitive CRUD while still allowing query-level control.
- Model tags align DB schema with code intent.
- Transactional service logic can use `db.Transaction` for atomic multi-step writes.

### 3.4 Firebase Admin SDK

Why this auth model was chosen:

- Client authenticates through Firebase; backend verifies Firebase ID tokens.
- Backend does not manage passwords, reducing credential-handling complexity.
- Verified token claims expose stable Firebase UID used to map local users.

### 3.5 Zerolog

Operational value:

- Structured logs with request metadata and request IDs.
- Minimal allocation overhead and good throughput.
- Easy ingestion in centralized log tooling.

---

## 4. Backend Components in Detail

## 4.1 Entrypoint: `cmd/server/main.go`

Responsibilities:

- Load runtime config (`config.Load`).
- Initialize logger and log level.
- Open PostgreSQL connection through GORM.
- Configure SQL pool limits.
- Initialize Firebase Admin auth client.
- Construct repositories, services, handlers.
- Build router with all dependencies.
- Start HTTP server with read/write/idle timeouts.
- Handle graceful shutdown on SIGINT/SIGTERM.

Important design notes:

- Dependency injection is manual and explicit in `main.go`, which keeps lifecycle and composition transparent.
- DB pooling is configured early with conservative defaults (`max open`, `max idle`, `conn max lifetime`).
- Firebase initialization failures are logged as warning instead of fatal, which allows startup but protected routes may fail until auth is configured.

## 4.2 Configuration: `internal/config/config.go`

Current behavior:

- Loads `.env` if present (`godotenv.Load()` best-effort).
- Requires `DATABASE_URL`.
- Defaults `PORT` to `8080` if missing.
- Reads `GIN_MODE`, CORS origins, and Firebase-related env values.

Operational caveat:

- Token verification path in `pkg/firebase` uses `FIREBASE_CREDENTIALS_PATH`.
- Keep `.env.example` and runtime env aligned around `FIREBASE_CREDENTIALS_PATH` for Firebase Admin initialization.

## 4.3 Router and Handlers: `internal/handler/*`

### Router

`internal/handler/router.go` wires:

- Public endpoints:
  - `GET /api/v1/health`
  - `POST /api/v1/auth/sync`
- Protected endpoints under auth middleware:
  - `/users/me`
  - `/accounts`
  - `/categories` and `/subcategories`
  - `/transactions`
  - `/quick-transactions`
  - `/stats/*`

### Handler pattern

Each handler typically:

1. Resolves authenticated user (via Firebase UID -> local user).
2. Binds and validates input.
3. Calls service method.
4. Converts model to DTO response envelope.

This keeps handlers thin and pushes business rules into services.

### API envelope contract

Responses use a standard contract from `internal/dto/response.go`:

- Success: `{ "data": ... }`
- Paginated: `{ "data": [...], "meta": {...} }`
- Error: `{ "error": { "code": "...", "message": "..." } }`

## 4.4 Middleware: `internal/middleware/*`

### Logger middleware

- Assigns/propagates `X-Request-ID`.
- Logs method, path, status, latency, client IP.
- Improves traceability for debugging and incident analysis.

### Recovery middleware

- Recovers panics and emits standardized 500 error.
- Prevents process crash from uncaught panics.

### Auth middleware

- Expects `Authorization: Bearer <token>`.
- Verifies token via Firebase Admin SDK.
- Injects `firebase_uid` and optional `email` into Gin context.
- Rejects malformed/missing/invalid tokens with 401.

## 4.5 Services: `internal/service/*`

Services implement business rules and cross-repository orchestration.

### User service

- Syncs Firebase-authenticated identity into local `users` table.
- Reads and updates profile fields.

### Account service

- CRUD for accounts.
- Persists audit entries for create/update/delete.

### Category and subcategory services

- CRUD + list operations.
- Duplicate-safe creation flow (returns existing where applicable instead of failing writes unnecessarily).
- Audit entries captured for write operations.

### Transaction service (core ledger logic)

This service is critical and intentionally transactional.

Key flows:

- **Create transaction**:
  - Validates account/category ownership.
  - Auto-resolves or creates a subcategory if absent.
  - In one DB transaction:
    - inserts transaction
    - updates account balance (`current_balance += amount`)
    - writes audit record

- **Update transaction**:
  - Calculates balance delta from old vs new amount.
  - In one DB transaction:
    - updates transaction
    - updates account balance by delta
    - writes audit record

- **Delete transaction**:
  - In one DB transaction:
    - soft-deletes transaction
    - reverses amount impact on account balance
    - writes audit record

Why this matters:

- Prevents drift between transaction ledger and account summary balances.
- Guarantees atomicity under partial failures.

### Quick transaction service

- Maintains reusable transaction templates.
- Can execute a template into a real transaction through transaction service.

### Stats service

- Computes summary totals.
- Aggregates timeline series by selected period and value.
- Aggregates category/subcategory contribution for donut charts.
- Supports account-level filtering.

## 4.6 Repository Layer: `internal/repository/*`

Repositories encapsulate persistence access patterns:

- Query filtering and sorting logic.
- Pagination and count retrieval.
- Ownership constraints (`user_id`) enforced at query level.

Notable repository behavior:

- Transaction list supports flexible filters (`account`, `category`, date range, search, type, sorting).
- Soft deletes are naturally respected by GORM model behavior.

## 4.7 Domain Models: `internal/model/*`

Models describe DB-backed entities:

- `User`, `Account`, `Category`, `Subcategory`, `Transaction`, `QuickTransaction`, `AuditLog`.

Design choices:

- UUID primary keys for distributed-safe identifiers.
- `numeric(20,4)` for money to avoid float representation issues at storage level.
- `AuditLog` is append-only and explicitly mapped to table `audit_log`.

## 4.8 Firebase Package: `pkg/firebase/firebase.go`

Key design:

- Singleton initialization via `sync.Once`.
- Reads service account file path from `FIREBASE_CREDENTIALS_PATH`.
- Exposes `Init` and `VerifyIDToken`.

Why this is useful:

- Prevents repeated expensive client initialization.
- Thread-safe by construction.
- Clear failure messages when auth is misconfigured.

---

## 5. Data and Schema Design Decisions

## 5.1 Migration-driven schema

- Schema defined in SQL migration files under `migrations/`.
- Versioned, repeatable DB evolution.
- Enables deterministic environment setup.

## 5.2 Soft-delete strategy

Most business entities use soft deletes, enabling:

- safer "delete" operations
- historical consistency for references
- possible future restore semantics

`audit_log` remains append-only and non-soft-deleted.

## 5.3 Denormalized account balance

`accounts.current_balance` is maintained transactionally for fast reads.

Trade-off:

- Pros: fast dashboard/account queries.
- Cons: must enforce atomic update logic rigorously (handled in transaction service).

---

## 6. Auth and Identity Flow

1. Frontend obtains Firebase ID token after sign-in.
2. Frontend calls `POST /api/v1/auth/sync` with token and profile info.
3. Backend verifies token through Firebase Admin.
4. Backend creates/updates local user by Firebase UID.
5. Frontend uses bearer token for protected routes.
6. Auth middleware verifies token on each protected request.

Rationale:

- Backend remains source of authorization and data ownership checks.
- Firebase handles credentialing and token issuance.

---

## 7. Request Lifecycle Example

Example: `POST /api/v1/transactions`

1. Request enters router.
2. `AuthMiddleware` verifies bearer token and sets `firebase_uid`.
3. Transaction handler resolves local user via `firebase_uid`.
4. Handler binds JSON DTO.
5. Transaction service validates account/category, prepares subcategory.
6. Service opens DB transaction:
   - create transaction row
   - adjust account balance
   - create audit row
7. Commit transaction.
8. Handler returns `201` with standardized DTO envelope.

Failure handling:

- Validation parse issues return 4xx envelopes.
- Runtime/repository errors return `ERR_INTERNAL` 500 envelopes.
- Panics are recovered by middleware and converted to safe 500 response.

---

## 8. How to Run the Server (Detailed)

This section is intentionally lengthy to cover common local setups and failure modes.

## 8.1 Prerequisites

- Go >= 1.24
- Node >= 22, npm >= 10 (for frontend if you run full stack)
- PostgreSQL 15+ (native or Docker)
- `migrate` CLI available in PATH
- `golangci-lint` available in PATH
- Firebase service account JSON for protected auth flows

## 8.2 Create environment files

From repo root:

```bash
cp .env.example .env
cp frontend/.env.example frontend/.env
```

Set backend `.env` values at minimum:

- `DATABASE_URL`
- `PORT` (optional, defaults to 8080)
- `CORS_ALLOWED_ORIGINS` (for local FE, usually `http://localhost:5173`)
- `FIREBASE_CREDENTIALS_PATH` (absolute path to service account JSON)
- `FIREBASE_PROJECT_ID`

Set frontend `frontend/.env` with `VITE_FIREBASE_*` values.

## 8.3 Start PostgreSQL (Docker path)

```bash
docker-compose up -d db
docker ps
```

Check connectivity quickly:

```bash
set -a && source .env && set +a
psql "$DATABASE_URL" -c '\dt'
```

## 8.4 Run migrations

```bash
set -a && source .env && set +a
migrate -path ./migrations -database "$DATABASE_URL" up
```

Or with make target:

```bash
set -a && source .env && set +a
make migrate-up
```

## 8.5 Start backend server

Development mode from source:

```bash
set -a && source .env && set +a
go run ./cmd/server
```

Binary mode:

```bash
go build -o bin/server ./cmd/server
set -a && source .env && set +a
./bin/server
```

Health check:

```bash
curl -iS http://localhost:8080/api/v1/health
```

Expected response includes:

- HTTP `200 OK`
- body `{"status":"ok"}`

## 8.6 (Optional) Start frontend for end-to-end testing

```bash
cd frontend
npm install
npm run dev -- --host 0.0.0.0 --port 5173
```

URLs:

- Frontend: `http://localhost:5173`
- API health: `http://localhost:8080/api/v1/health`

## 8.7 Build and quality gates

Backend gates:

```bash
go build ./...
go vet ./...
go test ./... -race -count=1
golangci-lint run
```

Frontend gates:

```bash
cd frontend
npx tsc --noEmit
npx eslint src/
npm run build
```

## 8.8 Go cache/workspace permission workaround

If your environment blocks writes to default Go paths, run with workspace-local cache variables:

```bash
mkdir -p .gopath .gocache .gomodcache
GOPATH=$PWD/.gopath GOCACHE=$PWD/.gocache GOMODCACHE=$PWD/.gomodcache go test ./...
```

Use the same pattern for `go build`, `go vet`, and `go run` if necessary.

## 8.9 Restarting safely after code changes

If you run a compiled binary (`./bin/server`), source changes do not apply until rebuild/restart.

Recommended restart sequence:

```bash
pkill -f '^./bin/server$|/bin/server$' || true
go build -o bin/server ./cmd/server
set -a && source .env && set +a
./bin/server
```

This avoids stale-process confusion during debugging.

---

## 9. Troubleshooting Guide

## 9.1 `401` on protected endpoints

Check:

- Authorization header format is `Bearer <id_token>`
- Firebase Admin initialized correctly (`FIREBASE_CREDENTIALS_PATH` valid)
- Token is current and issued for expected Firebase project

## 9.2 Firebase init warnings at startup

If startup logs warn about Firebase init skipped:

- verify `FIREBASE_CREDENTIALS_PATH` exists and is readable
- verify JSON is valid service account credentials
- verify project alignment between backend and frontend Firebase config

## 9.3 Migration errors

Typical causes:

- wrong `DATABASE_URL`
- DB not running
- migration version drift

Useful commands:

```bash
migrate -path ./migrations -database "$DATABASE_URL" version
migrate -path ./migrations -database "$DATABASE_URL" up
```

## 9.4 CORS issues in browser

Set `CORS_ALLOWED_ORIGINS` to include frontend origin(s), for example:

```dotenv
CORS_ALLOWED_ORIGINS=http://localhost:5173
```

Multiple origins are comma-separated.

## 9.5 Stats charts not showing useful data

Check:

- transactions exist in selected period/value
- `/api/v1/stats/summary` returns non-empty `timeseries` and breakdown arrays
- frontend selected filters (`period`, `value`, `account`) are valid

---

## 10. Design Decisions and Trade-Off Summary

1. Keep handlers thin, services rich
- Improves testability and separation of concerns.

2. Maintain `current_balance` denormalized
- Fast reads, with strict transactional updates to maintain correctness.

3. Append-only audit logging
- Better traceability for financial record mutations.

4. Firebase token verification server-side
- Strong ownership checks and trust boundary at backend.

5. Standardized API envelope
- Frontend can implement consistent success/error handling.

6. Migration-first schema management
- Reproducible environments and explicit schema history.

---

## 11. Suggested Ongoing Improvements

- Add richer domain-level error types and map them to more precise 4xx codes instead of broad 500s.
- Implement audit query endpoint (currently stubbed).
- Add reconciliation job/command to recompute account balances from ledger as integrity check.
- Expand structured logging with user/account identifiers on mutating endpoints.
- Add OpenAPI spec generation from route contracts for stronger API documentation automation.
