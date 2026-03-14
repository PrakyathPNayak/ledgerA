# DEBUG REPORT

Generated: 2026-03-15
Project: ledgerA

## DIAG-001: Backend reachability
Status: PASS

Command:
```bash
curl -sv http://localhost:8080/api/v1/health 2>&1
```
Output:
```text
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /api/v1/health HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.5.0
> Accept: */*
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< X-Request-Id: 3370ee55-40a7-4f00-8939-d60efc627e57
< Date: Sat, 14 Mar 2026 20:09:39 GMT
< Content-Length: 15
{"status":"ok"}
```

Command:
```bash
curl -sv http://localhost:8080/api/v1/accounts -H "Authorization: Bearer test" 2>&1 | head -30
```
Output:
```text
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /api/v1/accounts HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.5.0
> Accept: */*
> Authorization: Bearer test
< HTTP/1.1 401 Unauthorized
< Content-Type: application/json; charset=utf-8
{"error":{"code":"ERR_UNAUTHORIZED","message":"token verification failed"}}
```

## DIAG-002: Database reachable and migrated
Status: PASS (after remediation)

Command:
```bash
pg_isready -h localhost -p 5432
```
Output:
```text
localhost:5432 - accepting connections
```

Initial command:
```bash
psql $DATABASE_URL -c "\dt" 2>&1
```
Initial output:
```text
Did not find any relations.
```

Initial command:
```bash
migrate -path ./migrations -database "$DATABASE_URL" version 2>&1
```
Initial output:
```text
Command 'migrate' not found
```

Remediation command:
```bash
./migrate -path ./migrations -database "$DATABASE_URL" up 2>&1
```
Output:
```text
1/u create_users
2/u create_accounts
3/u create_categories
4/u create_transactions
5/u create_quick_transactions
6/u create_audit_log
```

Re-check command:
```bash
psql $DATABASE_URL -c "\dt" 2>&1
```
Output:
```text
List of relations:
accounts
audit_log
categories
quick_transactions
schema_migrations
subcategories
transactions
users
```

Re-check command:
```bash
./migrate -path ./migrations -database "$DATABASE_URL" version 2>&1
```
Output:
```text
6
```

## DIAG-003: Backend auth write-path behavior and logs
Status: PASS (endpoint behavior), UNKNOWN (service logs)

Command:
```bash
curl -sv -X POST http://localhost:8080/api/v1/auth/sync \
  -H "Content-Type: application/json" \
  -d '{"firebase_token":"fake","email":"test@test.com","display_name":"Test"}' \
  2>&1 | tail -20
```
Output:
```text
> POST /api/v1/auth/sync HTTP/1.1
< HTTP/1.1 401 Unauthorized
{"error":{"code":"ERR_UNAUTHORIZED","message":"token verification failed"}}
```

Command:
```bash
journalctl -u exptracker --no-pager -n 50 2>/dev/null || tail -50 /tmp/server.log 2>/dev/null || echo "Check terminal where server is running for logs"
```
Output:
```text
-- No entries --
```

## DIAG-004: Frontend proxy wiring
Status: PASS

Command:
```bash
cat frontend/vite.config.ts | grep -A 10 "proxy\|server"
```
Output:
```text
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
},
```

## DIAG-005: Direct endpoint checks for broken features
Status: PASS (routes exist, middleware reached)

Command:
```bash
curl -sv -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer fake_token_to_test_middleware" \
  -d '{"name":"Test Account","opening_balance":1000}' 2>&1 | tail -10
```
Output:
```text
< HTTP/1.1 401 Unauthorized
{"error":{"code":"ERR_UNAUTHORIZED","message":"token verification failed"}}
```

Command:
```bash
curl -sv -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer fake_token_to_test_middleware" \
  -d '{"name":"Food"}' 2>&1 | tail -10
```
Output:
```text
< HTTP/1.1 401 Unauthorized
{"error":{"code":"ERR_UNAUTHORIZED","message":"token verification failed"}}
```

## DIAG-006: Dark mode implementation
Status: PARTIAL

Command:
```bash
grep -r "dark\|theme\|document.documentElement" frontend/src/ --include="*.ts" --include="*.tsx" | grep -v node_modules | grep -v ".d.ts"
```
Output highlights:
```text
frontend/src/store/uiStore.ts: theme state exists
frontend/src/components/layout/PageShell.tsx: document.documentElement class manipulation exists
frontend/src/components/layout/SettingsSheet.tsx: theme + setTheme usage exists
```

Command:
```bash
grep "darkMode" frontend/tailwind.config.ts
```
Output:
```text
darkMode: ['class'],
```

Command:
```bash
grep -r "classList\|dark\|documentElement" frontend/src/store/ --include="*.ts" --include="*.tsx"
```
Output:
```text
frontend/src/store/uiStore.ts:type ThemeMode = 'light' | 'dark' | 'system'
```

Finding:
- Theme is currently applied in `PageShell.tsx`, not in the store itself.
- This can still work, but differs from desired robust pattern where store action applies DOM class immediately.

## DIAG-007: Auth flow end-to-end wiring
Status: PASS (env present and code wired)

Command:
```bash
grep "VITE_FIREBASE" frontend/.env 2>/dev/null || grep "VITE_FIREBASE" .env 2>/dev/null || echo "NO FIREBASE ENV VARS FOUND IN FRONTEND"
```
Output:
```text
VITE_FIREBASE_API_KEY=AIza...
VITE_FIREBASE_AUTH_DOMAIN=ledgera-a6a52.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=ledgera-a6a52
VITE_FIREBASE_APP_ID=1:1045712278229:web:4b01b5a610dd46cbcf485c
```

Command:
```bash
cat frontend/src/lib/firebase.ts
cat frontend/src/main.tsx
cat frontend/src/App.tsx | head -60
```
Finding:
- `firebase.ts` correctly consumes `import.meta.env.VITE_FIREBASE_*` values.
- Auth store is initialized in `main.tsx` before router render.
- Protected routes use `RequireAuth` and redirect to `/login` when unauthenticated.

## Additional high-signal finding (outside DIAG list)
Status: FAIL (external/project-side)

Independent API probe confirms the recurring login error source:
- Firebase Identity Toolkit project config endpoint returns `404 CONFIGURATION_NOT_FOUND`.
- Enabling API via service usage from this service account returns `403 PERMISSION_DENIED`.

Implication:
- Local code and env are now wired correctly.
- Firebase project-side Authentication/Identity Platform configuration and/or IAM permissions remain incomplete.
