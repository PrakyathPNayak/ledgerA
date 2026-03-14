import { PageShell } from '@/components/layout/PageShell'

const sections = [
    {
        title: 'Getting Started',
        body: 'Create at least one account, then add categories and your first transaction from Dashboard. Use Quick Transactions for repeated spends like rent, fuel, or subscriptions.',
    },
    {
        title: 'Understanding Income and Expense',
        body: 'Income entries should be positive amounts. Expense entries are stored as negative amounts. The net card on Dashboard and Stats reflects income minus expense.',
    },
    {
        title: 'Passbook Mode',
        body: 'Open any account from the Accounts page to see passbook-style chronology. This view helps reconcile account balance movement line by line.',
    },
    {
        title: 'Search and Filters',
        body: 'Use Search to combine full text, account, category, and transaction type filters. Load More fetches the next page without losing existing results.',
    },
    {
        title: 'Stats and PDF Export',
        body: 'Stats supports day/week/month/year windows. Export PDF downloads a shareable report snapshot for the selected filter set.',
    },
    {
        title: 'Offline and PWA',
        body: 'Install ledgerA from your browser menu. Cached assets allow quick start while offline. Mutations still require connectivity and are retried when connection resumes.',
    },
]

export function HelpPage() {
    return (
        <PageShell title="Help">
            <div className="space-y-4">
                <section className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
                    <h2 className="text-2xl font-bold tracking-tight text-slate-900">ledgerA Help Center</h2>
                    <p className="mt-2 text-sm leading-6 text-slate-600">
                        Find quick answers about accounts, transactions, reports, and productivity workflows.
                    </p>
                </section>

                <section className="rounded-2xl border border-slate-200 bg-white p-2 shadow-sm">
                    {sections.map((section) => (
                        <details key={section.title} className="group border-b border-slate-100 p-4 last:border-b-0">
                            <summary className="cursor-pointer list-none text-sm font-semibold text-slate-900">
                                <span className="inline-flex items-center gap-2">
                                    <span className="text-slate-400 transition-transform group-open:rotate-90">▶</span>
                                    {section.title}
                                </span>
                            </summary>
                            <p className="mt-3 text-sm leading-6 text-slate-600">{section.body}</p>
                        </details>
                    ))}
                </section>
            </div>
        </PageShell>
    )
}
