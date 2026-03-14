import { useState } from 'react'
import { Bar, BarChart, CartesianGrid, Pie, PieChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'

import { PageShell } from '@/components/layout/PageShell'
import { useAccounts } from '@/hooks/useAccounts'
import { useStats } from '@/hooks/useStats'
import { getValidToken } from '@/lib/firebase'

function toMonthValue(date: Date): string {
    return date.toISOString().slice(0, 7)
}

function money(value: number): string {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR' }).format(value)
}

export function StatsPage() {
    const [period, setPeriod] = useState<'day' | 'week' | 'month' | 'year'>('month')
    const [value, setValue] = useState(toMonthValue(new Date()))
    const [accountId, setAccountId] = useState('')
    const [isExporting, setIsExporting] = useState(false)

    const { data: accounts = [] } = useAccounts()
    const { data: stats, isLoading } = useStats({ period, value, account_id: accountId || undefined })

    async function handleExportPdf() {
        try {
            setIsExporting(true)
            const token = await getValidToken()
            const params = new URLSearchParams({ period, value })
            if (accountId) {
                params.set('account_id', accountId)
            }

            const response = await fetch(`/api/v1/stats/export/pdf?${params.toString()}`, {
                headers: token ? { Authorization: `Bearer ${token}` } : {},
            })
            if (!response.ok) {
                throw new Error('Failed to export PDF')
            }

            const blob = await response.blob()
            const url = URL.createObjectURL(blob)
            const anchor = document.createElement('a')
            anchor.href = url
            anchor.download = `ledgerA-stats-${period}-${value}.pdf`
            anchor.click()
            URL.revokeObjectURL(url)
        } finally {
            setIsExporting(false)
        }
    }

    const chartData = stats?.timeseries ?? []
    const expensePie = (stats?.category_breakdown_expense ?? []).slice(0, 8).map((item) => ({
        name: item.subcategory || item.category,
        value: Math.abs(item.amount),
    }))

    return (
        <PageShell title="Stats">
            <div className="space-y-6">
                <section className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
                    <div className="grid gap-3 md:grid-cols-4">
                        <label className="space-y-1">
                            <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">Period</span>
                            <select
                                value={period}
                                onChange={(e) => setPeriod(e.target.value as 'day' | 'week' | 'month' | 'year')}
                                className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm"
                            >
                                <option value="day">Day</option>
                                <option value="week">Week</option>
                                <option value="month">Month</option>
                                <option value="year">Year</option>
                            </select>
                        </label>

                        <label className="space-y-1">
                            <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">Value</span>
                            <input
                                value={value}
                                onChange={(e) => setValue(e.target.value)}
                                className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm"
                                placeholder="YYYY-MM"
                            />
                        </label>

                        <label className="space-y-1">
                            <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">Account</span>
                            <select
                                value={accountId}
                                onChange={(e) => setAccountId(e.target.value)}
                                className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm"
                            >
                                <option value="">All accounts</option>
                                {accounts.map((account) => (
                                    <option key={account.id} value={account.id}>
                                        {account.name}
                                    </option>
                                ))}
                            </select>
                        </label>

                        <div className="flex items-end">
                            <button
                                type="button"
                                onClick={handleExportPdf}
                                disabled={isExporting}
                                className="w-full rounded-lg bg-slate-900 px-4 py-2 text-sm font-medium text-white disabled:opacity-60"
                            >
                                {isExporting ? 'Exporting...' : 'Export PDF'}
                            </button>
                        </div>
                    </div>
                </section>

                <section className="grid gap-4 md:grid-cols-3">
                    <StatCard label="Income" value={money(stats?.total_income ?? 0)} tone="income" />
                    <StatCard label="Expense" value={money(stats?.total_expense ?? 0)} tone="expense" />
                    <StatCard label="Net" value={money(stats?.net ?? 0)} tone="neutral" />
                </section>

                <section className="grid gap-4 lg:grid-cols-2">
                    <article className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
                        <h2 className="text-base font-semibold text-slate-900">Income vs Expense</h2>
                        <div className="mt-4 h-72">
                            {isLoading ? (
                                <p className="text-sm text-slate-500">Loading chart...</p>
                            ) : (
                                <ResponsiveContainer width="100%" height="100%">
                                    <BarChart data={chartData}>
                                        <CartesianGrid strokeDasharray="3 3" />
                                        <XAxis dataKey="label" />
                                        <YAxis />
                                        <Tooltip formatter={(v) => money(Number(v))} />
                                        <Bar dataKey="income" fill="#059669" radius={[8, 8, 0, 0]} />
                                        <Bar dataKey="expense" fill="#dc2626" radius={[8, 8, 0, 0]} />
                                    </BarChart>
                                </ResponsiveContainer>
                            )}
                        </div>
                    </article>

                    <article className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
                        <h2 className="text-base font-semibold text-slate-900">Expense Breakdown</h2>
                        <div className="mt-4 h-72">
                            {isLoading ? (
                                <p className="text-sm text-slate-500">Loading chart...</p>
                            ) : (
                                <ResponsiveContainer width="100%" height="100%">
                                    <PieChart>
                                        <Pie
                                            data={expensePie}
                                            dataKey="value"
                                            nameKey="name"
                                            innerRadius={55}
                                            outerRadius={100}
                                            paddingAngle={2}
                                            fill="#1d4ed8"
                                        />
                                        <Tooltip formatter={(v) => money(Number(v))} />
                                    </PieChart>
                                </ResponsiveContainer>
                            )}
                        </div>
                    </article>
                </section>
            </div>
        </PageShell>
    )
}

function StatCard({ label, value, tone }: { label: string; value: string; tone: 'income' | 'expense' | 'neutral' }) {
    const toneClass =
        tone === 'income'
            ? 'border-emerald-200 bg-emerald-50'
            : tone === 'expense'
                ? 'border-rose-200 bg-rose-50'
                : 'border-slate-200 bg-slate-50'

    return (
        <article className={`rounded-2xl border p-4 ${toneClass}`}>
            <p className="text-xs font-semibold uppercase tracking-wide text-slate-500">{label}</p>
            <p className="mt-2 text-2xl font-bold text-slate-900">{value}</p>
        </article>
    )
}
