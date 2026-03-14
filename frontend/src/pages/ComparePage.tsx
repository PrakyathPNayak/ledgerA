import { PageShell } from '@/components/layout/PageShell'

export function ComparePage() {
    return (
        <PageShell title="Compare">
            <section className="relative overflow-hidden rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">
                <div className="absolute -right-8 -top-8 h-28 w-28 rounded-full bg-sky-100" />
                <div className="absolute -bottom-10 -left-10 h-36 w-36 rounded-full bg-amber-100" />

                <div className="relative z-10 max-w-2xl">
                    <p className="inline-flex rounded-full bg-slate-900 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-white">
                        Coming Soon
                    </p>
                    <h2 className="mt-4 text-3xl font-bold tracking-tight text-slate-900">Account vs Account comparison</h2>
                    <p className="mt-3 text-sm leading-6 text-slate-600">
                        This workspace is prepared for side-by-side trend comparisons with aligned date ranges,
                        variance highlights, and category-level deltas. The data contract is already available on the
                        backend via `/stats/compare`.
                    </p>
                    <ul className="mt-6 grid gap-2 text-sm text-slate-700 md:grid-cols-2">
                        <li className="rounded-lg bg-slate-50 px-3 py-2">Period-over-period variance</li>
                        <li className="rounded-lg bg-slate-50 px-3 py-2">Income vs expense drift</li>
                        <li className="rounded-lg bg-slate-50 px-3 py-2">Top category contributors</li>
                        <li className="rounded-lg bg-slate-50 px-3 py-2">Export and share snapshots</li>
                    </ul>
                </div>
            </section>
        </PageShell>
    )
}
