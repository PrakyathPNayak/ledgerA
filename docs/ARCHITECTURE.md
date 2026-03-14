# ARCHITECTURE

## Pre-Build Review

[DA: Auth Strategy]
Decision: Verify Firebase tokens on backend every protected request.
Challenge: High latency due to network verification.
Evidence: Firebase Admin SDK leverages cached public keys; verification happens locally w/o external network call per request.
Rebuttal: It operates offline using standard JWT signature verification. No latency issue.
Verdict: PROCEED

[DA: Balance Maintenance]
Decision: App logic using DB transactions and pessimistic locking (SELECT FOR UPDATE).
Challenge: Race conditions resulting in corrupted balances.
Evidence: Concurrent requests inserting transactions for the same account.
Rebuttal: `FOR UPDATE` locks the account row until transaction commit, guaranteeing serial consistency.
Verdict: REVISED (App logic with explicit row-locking)

[DA: PDF Library Choice]
Decision: go-pdf/fpdf over gofpdf.
Challenge: gofpdf is the original and more well-known.
Evidence: gofpdf repository is archived and explicitly states no further updates.
Rebuttal: go-pdf/fpdf is the modern, maintained fork.
Verdict: PROCEED

[DA: State Management Choice]
Decision: TanStack Query (server state) + Zustand (UI state).
Challenge: Too many libraries. Why not use Context API?
Evidence: Context triggers re-renders for all consumers indiscriminately on any change.
Rebuttal: Zustand allows slice-level subscription. TanStack Query handles caching logic that Context would require thousands of lines to replicate.
Verdict: PROCEED

[DA: PWA Caching Strategy]
Decision: NetworkFirst for `/api/*`, CacheFirst for static assets.
Challenge: Offline mode will fail on POST/PATCH easily.
Evidence: Workbox requires BackgroundSync plugin.
Rebuttal: For now, we will notify the user they are offline and disable writes if offline. A read-only offline state is acceptable for v1.
Verdict: PROCEED

[ToT: current_balance storage]
Branch A: Trigger based (DB maintains balance). Score: 7/10. Hard to test.
Branch B: Computed view (SUM every time). Score: 6/10. Bad perf as data grows.
Branch C: App logic with Transaction + locking. Score: 9/10. Testable, performant, predictable.
Chosen: Branch C

[ToT: running balance calculation strategy]
Branch A: Calculate entirely in DB using Window function. Score: 9/10. Exact.
Branch B: Calculate in frontend based on rows. Score: 4/10. Fails with pagination.
Chosen: Branch A. We will use `SUM(amount) OVER (PARTITION BY account_id ORDER BY transaction_date, created_at)` in query or just backend sequential summation depending on route.
