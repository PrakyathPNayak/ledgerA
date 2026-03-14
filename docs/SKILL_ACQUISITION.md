# Skill Acquisition Log

[ReAct Loop Executions for Section 1]

Thought: I need to verify standard Go conventions before writing backend code.
Action: FETCH https://go.dev/doc/effective_go
Observation:
- Use clear, short, concise names for variables and packages.
- Always handle errors.
- Capitalize for exported symbols, lowercase for unexported.
- Use gofmt for formatting.

Thought: What are the community agreed upon code review comments?
Action: FETCH https://github.com/golang/go/wiki/CodeReviewComments
Observation:
- Interfaces belong in the package that uses them, not the one that implements them.
- Avoid panic.
- Don't use `this` or `self` as receiver names.

Thought: What is the standard folder structure for Go?
Action: FETCH https://github.com/golang-standards/project-layout/blob/master/README.md
Observation:
- `cmd/`: main applications.
- `internal/`: private application and library code.
- `pkg/`: library code that's ok to use by external applications.

Thought: I need to review Gin routing and middleware.
Action: FETCH https://gin-gonic.com/docs/quickstart/
Observation:
- Handlers take `*gin.Context`.
- Context handles payload bindings (`c.ShouldBindJSON`).
- Middleware uses `c.Next()` to pass control, `c.Abort()` to stop it.

Thought: How does GORM handle associations and soft deletes?
Action: FETCH https://gorm.io/docs/models.html
Observation:
- `gorm.DeletedAt` enables soft deletes natively.
- Use Preload to eager load associations.
- Transactions are easily wrapped with `db.Transaction(func(tx *gorm.DB) error)`.

Thought: How to configure golang-migrate?
Action: FETCH https://github.com/golang-migrate/migrate/tree/master/database/postgres
Observation:
- Format: `{version}_{title}.up.sql` and `{version}_{title}.down.sql`.
- Needs `postgres://user:pass@host/db?sslmode=disable`.

Thought: Firebase Admin SDK auth verification in Go.
Action: FETCH https://firebase.google.com/docs/auth/admin/verify-id-tokens
Observation:
- Initialize Firebase App using default credentials or specific SA JSON.
- `client.VerifyIDToken(ctx, idToken)` returns `*auth.Token`.
- UID is available at `token.UID`.

Thought: Fetching React 19 new features.
Action: FETCH https://react.dev/blog/2024/04/25/react-19
Observation:
- Includes new hooks like `useActionState` and `useFormStatus` for forms.
- Replaces forwardRef with ref as a prop.
- Optimized rendering under the hood.

Thought: Reviewing TanStack Query v5 patterns.
Action: FETCH https://tanstack.com/query/latest/docs/framework/react/overview
Observation:
- `useQuery` object syntax only now.
- Provides `isLoading`, `isError`, `data`.
- Mutations use `useMutation` with `onSuccess` for query invalidation.

Thought: How to configure Shadcn/ui in Vite?
Action: FETCH https://ui.shadcn.com/docs/installation/vite
Observation:
- Requires tailwindcss setup.
- Initialization via `npx shadcn@latest init`.
- Components are injected directly into `src/components/ui/`.

Thought: Postgres numeric and JSONB fields.
Action: FETCH https://www.postgresql.org/docs/15/datatype-numeric.html
Observation:
- `NUMERIC(20,4)` is ideal for large monetary values to avoid float precision loss.
- UUID works nicely with `gen_random_uuid()` in PG14+.
- JSONB is excellent for schemaless data like audit log diffs.

Thought: AWS EC2 t2.micro limits.
Action: FETCH https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html
Observation:
- Instance has 1 vCPU and 1GB RAM.
- Keep Postgres + Go footprints reasonably constrained to avoid OOM killer.
- Need a proper swap file just in case.

[End ReAct Loop Log]
