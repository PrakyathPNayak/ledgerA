# ledgerA

ledgerA is a full-stack expenditure tracker PWA with Firebase-authenticated users, account-level balance tracking, category/subcategory classification, quick transaction templates, and period-based analytics with PDF export.

## Tech Stack
- Backend: Go, Gin, GORM, PostgreSQL, Firebase Admin SDK, Zerolog
- Frontend: React, TypeScript, Vite, Tailwind, React Query, Zustand, Recharts, dnd-kit, vite-plugin-pwa
- Tooling: golangci-lint, gofmt, migrate, Docker Compose

## Prerequisites
- Go 1.24+
- Node.js 22+
- npm 10+
- PostgreSQL 15+ (or Docker)
- `golangci-lint` available in PATH
- `migrate` CLI (for SQL migrations)

## Quick Start
1. Start PostgreSQL:
```bash
docker-compose up -d db
```
2. Configure environment:
```bash
cp .env.example .env
```
3. Set required values in `.env` (see table below).
4. Apply migrations:
```bash
migrate -path ./migrations -database "$DATABASE_URL" up
```
5. Start backend:
```bash
go run ./cmd/server
```
6. Start frontend (new terminal):
```bash
cd frontend
npm install
npm run dev -- --host 0.0.0.0 --port 5173
```

Frontend: `http://localhost:5173`
Backend health: `http://localhost:8080/api/v1/health`

## Environment Variables
| Name | Required | Example | Description |
|---|---|---|---|
| `DATABASE_URL` | Yes | `postgres://postgres:password@localhost:5432/exptracker?sslmode=disable` | PostgreSQL DSN for backend |
| `PORT` | No | `8080` | Backend HTTP port |
| `GIN_MODE` | No | `release` | Gin mode (`debug`, `release`, `test`) |
| `CORS_ALLOWED_ORIGINS` | No | `http://localhost:5173` | Comma-separated allowed origins |
| `FIREBASE_PROJECT_ID` | No* | `my-project-id` | Firebase project identifier |
| `FIREBASE_CREDENTIALS` | No* | `/path/service-account.json` | Firebase credentials file path |
| `FIREBASE_CREDENTIALS_JSON` | No* | `{...}` | Raw service account JSON |
| `VITE_FIREBASE_API_KEY` | Yes (frontend) | `AIza...` | Firebase JS SDK API key |
| `VITE_FIREBASE_AUTH_DOMAIN` | Yes (frontend) | `my-project.firebaseapp.com` | Firebase auth domain |
| `VITE_FIREBASE_PROJECT_ID` | Yes (frontend) | `my-project-id` | Firebase project for web SDK |
| `VITE_FIREBASE_APP_ID` | Yes (frontend) | `1:...:web:...` | Firebase app ID |

\* Firebase admin credentials are required for protected backend auth flows.

## Make Targets
| Target | Command | Purpose |
|---|---|---|
| `dev` | `docker-compose up -d` | Start local infra services |
| `build` | `go build -o bin/server ./cmd/server` | Build backend binary |
| `test` | `go test ./... -race -count=1` | Run backend tests with race detector |
| `lint` | `golangci-lint run` | Run backend linters |
| `fmt` | `gofmt -s -w .` | Format Go code |
| `migrate-up` | `migrate -path ./migrations -database "$DATABASE_URL" up` | Apply migrations |
| `migrate-down` | `migrate -path ./migrations -database "$DATABASE_URL" down` | Roll back migrations |

## Development Gates
Frontend:
```bash
cd frontend
npx tsc --noEmit
npx eslint src/ --max-warnings 0
npx vite build
```

Backend:
```bash
go build ./...
go vet ./...
go test ./... -race -count=1
golangci-lint run
```

## Project Structure
- `cmd/server` backend entrypoint
- `internal/config` env loading and runtime config
- `internal/model` domain models
- `internal/repository` database access layer
- `internal/service` business logic and tests
- `internal/handler` HTTP handlers/router
- `internal/middleware` auth/logger/recovery middleware
- `pkg/firebase` Firebase admin integration
- `pkg/pdf` PDF report generation
- `migrations` SQL schema migrations
- `frontend` React PWA

## Deep Documentation
- Server internals and operations handbook: `docs/SERVER_DEEP_DIVE.md`
