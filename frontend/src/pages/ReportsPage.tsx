import { useState } from 'react'
import {
    Bar,
    BarChart,
    CartesianGrid,
    Legend,
    Line,
    LineChart,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
} from 'recharts'

import { PageShell } from '@/components/layout/PageShell'
import { useAccounts } from '@/hooks/useAccounts'
import { useMonthlyReport } from '@/hooks/useStats'
import type { MonthlyDataPoint } from '@/types'

function formatMoney(value: number): string {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }).format(value)
}

type BarTooltipProps = {
    active?: boolean
    payload?: Array<{ dataKey?: string; value?: number; fill?: string }>
    label?: string
}

function BarTooltip({ active, payload, label }: BarTooltipProps) {
    if (!active || !payload || payload.length === 0) return null
    const income = payload.find((p) => p.dataKey === 'total_income')?.value ?? 0
    const expense = payload.find((p) => p.dataKey === 'total_expense')?.value ?? 0
    const net = income - expense
    return (
        <div className="rounded-lg border border-border bg-elevated p-3 text-sm shadow-lg">
            <p className="mb-1 font-semibold text-foreground">{label}</p>
            <p className="text-positive">Income: {formatMoney(income)}</p>
            <p className="text-negative">Expense: {formatMoney(expense)}</p>
            <p className={net >= 0 ? 'text-positive' : 'text-negative'}>Net: {formatMoney(net)}</p>
        </div>
    )
}

type NetTooltipProps = {
    active?: boolean
    payload?: Array<{ value?: number }>
    label?: string
}

function NetTooltip({ active, payload, label }: NetTooltipProps) {
    if (!active || !payload || payload.length === 0) return null
    const net = payload[0]?.value ?? 0
    return (
        <div className="rounded-lg border border-border bg-elevated p-3 text-sm shadow-lg">
            <p className="mb-1 font-semibold text-foreground">{label}</p>
            <p className={net >= 0 ? 'text-positive' : 'text-negative'}>Net: {formatMoney(net)}</p>
        </div>
    )
}

function SummaryTable({ months }: { months: MonthlyDataPoint[] }) {
    const totalIncome = months.reduce((s, m) => s + m.total_income, 0)
    const totalExpense = months.reduce((s, m) => s + m.total_expense, 0)
    const totalNet = totalIncome - totalExpense

    return (
        <div className="overflow-x-auto rounded-lg border border-border">
            <table className="w-full text-sm">
                <thead>
                    <tr className="border-b border-border bg-elevated text-left text-muted">
                        <th className="px-4 py-2">Month</th>
                        <th className="px-4 py-2 text-right">Income</th>
                        <th className="px-4 py-2 text-right">Expense</th>
                        <th className="px-4 py-2 text-right">Net</th>
                        <th className="px-4 py-2 text-right">Top Expense Category</th>
                    </tr>
                </thead>
                <tbody>
                    {[...months].reverse().map((m) => {
                        const net = m.total_income - m.total_expense
                        const topCat = m.top_expense[0]?.category ?? '—'
                        return (
                            <tr key={m.month} className="border-b border-border last:border-0 hover:bg-surface-hover">
                                <td className="px-4 py-2 font-medium text-foreground">{m.month}</td>
                                <td className="px-4 py-2 text-right text-positive">{formatMoney(m.total_income)}</td>
                                <td className="px-4 py-2 text-right text-negative">{formatMoney(m.total_expense)}</td>
                                <td className={`px-4 py-2 text-right font-semibold ${net >= 0 ? 'text-positive' : 'text-negative'}`}>
                                    {formatMoney(net)}
                                </td>
                                <td className="px-4 py-2 text-right text-secondary">{topCat}</td>
                            </tr>
                        )
                    })}
                </tbody>
                <tfoot>
                    <tr className="bg-elevated font-semibold text-foreground">
                        <td className="px-4 py-2">Total ({months.length} months)</td>
                        <td className="px-4 py-2 text-right text-positive">{formatMoney(totalIncome)}</td>
                        <td className="px-4 py-2 text-right text-negative">{formatMoney(totalExpense)}</td>
                        <td className={`px-4 py-2 text-right ${totalNet >= 0 ? 'text-positive' : 'text-negative'}`}>
                            {formatMoney(totalNet)}
                        </td>
                        <td />
                    </tr>
                </tfoot>
            </table>
        </div>
    )
}

export function ReportsPage() {
    const [monthsCount, setMonthsCount] = useState(12)
    const [accountId, setAccountId] = useState('')

    const { data: accountsData } = useAccounts()
    const accounts = accountsData?.data ?? []

    const filters = { months: monthsCount, ...(accountId ? { account_id: accountId } : {}) }
    const { data, isLoading, error } = useMonthlyReport(filters)

    const months = data?.months ?? []

    // Build top-5 categories across all months for trend chart
    const categoryTotals: Record<string, number> = {}
    months.forEach((m) => {
        m.top_expense.forEach((item) => {
            categoryTotals[item.category] = (categoryTotals[item.category] ?? 0) + item.amount
        })
    })
    const topCategories = Object.entries(categoryTotals)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 5)
        .map(([cat]) => cat)

    const trendData = months.map((m) => {
        const row: Record<string, string | number> = { month: m.month }
        topCategories.forEach((cat) => {
            const found = m.top_expense.find((item) => item.category === cat)
            row[cat] = found?.amount ?? 0
        })
        return row
    })

    const TREND_COLORS = ['#ef4444', '#f97316', '#eab308', '#10b981', '#3b82f6']

    return (
        <PageShell title="Monthly Reports">
            {/* Controls */}
            <div className="mb-6 flex flex-wrap items-center gap-3">
                <div className="flex items-center gap-2">
                    <label className="text-sm text-muted">Period:</label>
                    <select
                        value={monthsCount}
                        onChange={(e) => setMonthsCount(Number(e.target.value))}
                        className="rounded-lg border border-border bg-elevated px-3 py-1.5 text-sm text-foreground focus:ring-1 focus:ring-accent focus:outline-none"
                    >
                        <option value={3}>Last 3 months</option>
                        <option value={6}>Last 6 months</option>
                        <option value={12}>Last 12 months</option>
                        <option value={24}>Last 24 months</option>
                    </select>
                </div>
                <div className="flex items-center gap-2">
                    <label className="text-sm text-muted">Account:</label>
                    <select
                        value={accountId}
                        onChange={(e) => setAccountId(e.target.value)}
                        className="rounded-lg border border-border bg-elevated px-3 py-1.5 text-sm text-foreground focus:ring-1 focus:ring-accent focus:outline-none"
                    >
                        <option value="">All Accounts</option>
                        {accounts.map((a) => (
                            <option key={a.id} value={a.id}>{a.name}</option>
                        ))}
                    </select>
                </div>
            </div>

            {isLoading && (
                <div className="grid place-items-center py-20 text-sm text-muted">Loading monthly report…</div>
            )}
            {error && (
                <div className="rounded-lg border border-border bg-negative-muted px-4 py-3 text-sm text-negative">
                    Failed to load monthly report.
                </div>
            )}

            {!isLoading && !error && months.length > 0 && (
                <div className="space-y-8">
                    {/* Income vs Expense bar chart */}
                    <div className="rounded-xl border border-border bg-surface p-4">
                        <h2 className="mb-4 text-sm font-semibold text-secondary uppercase tracking-wide">Income vs Expense</h2>
                        <ResponsiveContainer width="100%" height={280}>
                            <BarChart data={months} margin={{ top: 4, right: 16, left: 0, bottom: 0 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" />
                                <XAxis dataKey="month" tick={{ fontSize: 11, fill: 'var(--color-text-muted)' }} />
                                <YAxis tick={{ fontSize: 11, fill: 'var(--color-text-muted)' }} tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`} />
                                <Tooltip content={<BarTooltip />} />
                                <Legend wrapperStyle={{ fontSize: 12, color: 'var(--color-text-secondary)' }} />
                                <Bar dataKey="total_income" name="Income" fill="var(--color-positive)" radius={[4, 4, 0, 0]} />
                                <Bar dataKey="total_expense" name="Expense" fill="var(--color-negative)" radius={[4, 4, 0, 0]} />
                            </BarChart>
                        </ResponsiveContainer>
                    </div>

                    {/* Net savings line chart */}
                    <div className="rounded-xl border border-border bg-surface p-4">
                        <h2 className="mb-4 text-sm font-semibold text-secondary uppercase tracking-wide">Net Savings Trend</h2>
                        <ResponsiveContainer width="100%" height={220}>
                            <LineChart data={months} margin={{ top: 4, right: 16, left: 0, bottom: 0 }}>
                                <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" />
                                <XAxis dataKey="month" tick={{ fontSize: 11, fill: 'var(--color-text-muted)' }} />
                                <YAxis tick={{ fontSize: 11, fill: 'var(--color-text-muted)' }} tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`} />
                                <Tooltip content={<NetTooltip />} />
                                <Line
                                    type="monotone"
                                    dataKey="net"
                                    name="Net"
                                    stroke="var(--color-accent)"
                                    strokeWidth={2}
                                    dot={{ r: 4 }}
                                    activeDot={{ r: 6 }}
                                />
                            </LineChart>
                        </ResponsiveContainer>
                    </div>

                    {/* Category spending trends */}
                    {topCategories.length > 0 && (
                        <div className="rounded-xl border border-border bg-surface p-4">
                            <h2 className="mb-4 text-sm font-semibold text-secondary uppercase tracking-wide">
                                Top Category Spend Trends
                            </h2>
                            <ResponsiveContainer width="100%" height={260}>
                                <LineChart data={trendData} margin={{ top: 4, right: 16, left: 0, bottom: 0 }}>
                                    <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" />
                                    <XAxis dataKey="month" tick={{ fontSize: 11, fill: 'var(--color-text-muted)' }} />
                                    <YAxis tick={{ fontSize: 11, fill: 'var(--color-text-muted)' }} tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`} />
                                    <Tooltip
                                        contentStyle={{
                                            background: 'var(--color-elevated)',
                                            border: '1px solid var(--color-border)',
                                            borderRadius: '8px',
                                            fontSize: '12px',
                                            color: 'var(--color-text-foreground)',
                                        }}
                                        formatter={(val: number) => formatMoney(val)}
                                    />
                                    <Legend wrapperStyle={{ fontSize: 12, color: 'var(--color-text-secondary)' }} />
                                    {topCategories.map((cat, idx) => (
                                        <Line
                                            key={cat}
                                            type="monotone"
                                            dataKey={cat}
                                            name={cat}
                                            stroke={TREND_COLORS[idx % TREND_COLORS.length]}
                                            strokeWidth={2}
                                            dot={{ r: 3 }}
                                        />
                                    ))}
                                </LineChart>
                            </ResponsiveContainer>
                        </div>
                    )}

                    {/* Summary table */}
                    <div>
                        <h2 className="mb-3 text-sm font-semibold text-secondary uppercase tracking-wide">Monthly Summary</h2>
                        <SummaryTable months={months} />
                    </div>
                </div>
            )}

            {!isLoading && !error && months.length === 0 && (
                <div className="grid place-items-center py-20 text-sm text-muted">No data available for this period.</div>
            )}
        </PageShell>
    )
}
