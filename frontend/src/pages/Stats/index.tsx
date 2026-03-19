import { useEffect, useMemo, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import {
    Bar,
    CartesianGrid,
    Cell,
    ComposedChart,
    Pie,
    PieChart,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
} from 'recharts'

import { PageShell } from '@/components/layout/PageShell'
import { useAccounts } from '@/hooks/useAccounts'
import { useStats } from '@/hooks/useStats'
import { getValidToken } from '@/lib/firebase'
import type { Period } from '@/types'

type DonutDatum = {
    name: string
    amount: number
    percentage: number
}

type BarsTooltipProps = {
    active?: boolean
    label?: string
    payload?: Array<{ dataKey?: string; value?: number }>
}

type DonutTooltipProps = {
    active?: boolean
    payload?: Array<{ payload?: DonutDatum }>
}

const EXPENSE_COLORS = ['#ef4444', '#f97316', '#f59e0b', '#eab308', '#84cc16', '#10b981', '#06b6d4', '#3b82f6']
const INCOME_COLORS = ['#22c55e', '#14b8a6', '#06b6d4', '#0ea5e9', '#3b82f6', '#6366f1', '#8b5cf6', '#a855f7']

function formatMoney(value: number): string {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR' }).format(value)
}

function pad(n: number): string {
    return String(n).padStart(2, '0')
}

function isoWeek(date: Date): string {
    const temp = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()))
    const day = temp.getUTCDay() || 7
    temp.setUTCDate(temp.getUTCDate() + 4 - day)
    const yearStart = new Date(Date.UTC(temp.getUTCFullYear(), 0, 1))
    const week = Math.ceil(((temp.getTime() - yearStart.getTime()) / 86400000 + 1) / 7)
    return `${temp.getUTCFullYear()}-W${pad(week)}`
}

function getValueOptions(period: Period): Array<{ value: string; label: string }> {
    const now = new Date()

    if (period === 'day') {
        return Array.from({ length: 14 }, (_, i) => {
            const d = new Date(now)
            d.setDate(now.getDate() - i)
            const value = `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
            return { value, label: value }
        })
    }

    if (period === 'week') {
        return Array.from({ length: 12 }, (_, i) => {
            const d = new Date(now)
            d.setDate(now.getDate() - i * 7)
            const value = isoWeek(d)
            return { value, label: value }
        })
    }

    if (period === 'month') {
        return Array.from({ length: 12 }, (_, i) => {
            const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
            const value = `${d.getFullYear()}-${pad(d.getMonth() + 1)}`
            return { value, label: value }
        })
    }

    return Array.from({ length: 6 }, (_, i) => {
        const year = String(now.getFullYear() - i)
        return { value: year, label: year }
    })
}

function BarsTooltip({ active, payload, label }: BarsTooltipProps) {
    if (!active || !payload || payload.length === 0) {
        return null
    }

    const income = Number(payload.find((item) => item.dataKey === 'income')?.value ?? 0)
    const expense = Number(payload.find((item) => item.dataKey === 'expense')?.value ?? 0)
    const net = income + expense

    return (
        <div className="rounded-lg border border-border bg-surface p-3 text-xs shadow-lg border-border bg-elevated">
            <p className="font-semibold text-foreground">{label}</p>
            <p className="mt-1 text-positive">Income: {formatMoney(income)}</p>
            <p className="text-negative">Expense: {formatMoney(expense)}</p>
            <p className="mt-1 font-semibold text-foreground">Net: {formatMoney(net)}</p>
        </div>
    )
}

function DonutTooltip({ active, payload }: DonutTooltipProps) {
    if (!active || !payload || payload.length === 0 || !payload[0]?.payload) {
        return null
    }

    const item = payload[0].payload

    return (
        <div className="rounded-lg border border-border bg-surface p-3 text-xs shadow-lg border-border bg-elevated">
            <p className="font-semibold text-foreground">{item.name}</p>
            <p className="mt-1 text-secondary">Amount: {formatMoney(item.amount)}</p>
            <p className="text-secondary text-muted">Share: {item.percentage.toFixed(2)}%</p>
        </div>
    )
}

export function StatsPage() {
    const [searchParams] = useSearchParams()
    const [period, setPeriod] = useState<Period>('month')
    const [value, setValue] = useState('')
    const [accountId, setAccountId] = useState('')
    const [isExporting, setIsExporting] = useState(false)

    const { data: accounts = [] } = useAccounts()

    const valueOptions = useMemo(() => getValueOptions(period), [period])

    useEffect(() => {
        if (!valueOptions.find((option) => option.value === value)) {
            setValue(valueOptions[0]?.value ?? '')
        }
    }, [period, value, valueOptions])

    useEffect(() => {
        const fromUrl = searchParams.get('account_id')
        if (fromUrl) {
            setAccountId(fromUrl)
        }
    }, [searchParams])

    const { data: stats, isLoading } = useStats({
        period,
        value,
        account_id: accountId || undefined,
    })

    const barData = stats?.timeseries ?? []
    const expenseData: DonutDatum[] = (stats?.category_breakdown_expense ?? []).slice(0, 8).map((item) => ({
        name: item.subcategory || item.category,
        amount: Math.abs(item.amount),
        percentage: item.percentage,
    }))
    const incomeData: DonutDatum[] = (stats?.category_breakdown_income ?? []).slice(0, 8).map((item) => ({
        name: item.subcategory || item.category,
        amount: Math.abs(item.amount),
        percentage: item.percentage,
    }))

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
            const a = document.createElement('a')
            a.href = url
            a.download = `ledgerA-stats-${period}-${value}.pdf`
            a.click()
            URL.revokeObjectURL(url)
        } finally {
            setIsExporting(false)
        }
    }

    return (
        <PageShell title="Stats">
            <div className="space-y-6">
                <section className="rounded-2xl border border-border bg-surface p-4 shadow-sm border-border bg-surface">
                    <div className="grid gap-3 md:grid-cols-4">
                        <label className="space-y-1">
                            <span className="text-xs font-semibold uppercase tracking-wide text-muted">Period</span>
                            <select
                                value={period}
                                onChange={(event) => setPeriod(event.target.value as Period)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                            >
                                <option value="day">Day</option>
                                <option value="week">Week</option>
                                <option value="month">Month</option>
                                <option value="year">Year</option>
                            </select>
                        </label>

                        <label className="space-y-1">
                            <span className="text-xs font-semibold uppercase tracking-wide text-muted">Value</span>
                            <select
                                value={value}
                                onChange={(event) => setValue(event.target.value)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                            >
                                {valueOptions.map((option) => (
                                    <option key={option.value} value={option.value}>
                                        {option.label}
                                    </option>
                                ))}
                            </select>
                        </label>

                        <label className="space-y-1">
                            <span className="text-xs font-semibold uppercase tracking-wide text-muted">Account</span>
                            <select
                                value={accountId}
                                onChange={(event) => setAccountId(event.target.value)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm"
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
                                className="w-full rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white disabled:opacity-60"
                            >
                                {isExporting ? 'Exporting...' : 'Export PDF'}
                            </button>
                        </div>
                    </div>
                </section>

                <section className="grid gap-4 md:grid-cols-3">
                    <StatCard label="Income" value={formatMoney(stats?.total_income ?? 0)} tone="income" />
                    <StatCard label="Expense" value={formatMoney(stats?.total_expense ?? 0)} tone="expense" />
                    <StatCard label="Net" value={formatMoney(stats?.net ?? 0)} tone="neutral" />
                </section>

                <section className="rounded-2xl border border-border bg-surface p-4 shadow-sm border-border bg-surface">
                    <h2 className="text-base font-semibold text-foreground">Income and Expense Timeline</h2>
                    <div className="mt-4" style={{ height: 320 }}>
                        {isLoading ? (
                            <p className="text-sm text-muted">Loading chart...</p>
                        ) : (
                            <ResponsiveContainer width="100%" height={320}>
                                <ComposedChart data={barData} margin={{ top: 10, right: 8, left: 0, bottom: 0 }}>
                                    <CartesianGrid strokeDasharray="3 3" />
                                    <XAxis dataKey="label" />
                                    <YAxis />
                                    <Tooltip content={<BarsTooltip />} />
                                    <Bar dataKey="income" fill="#16a34a" radius={[8, 8, 0, 0]} />
                                    <Bar dataKey="expense" fill="#dc2626" radius={[8, 8, 0, 0]} />
                                </ComposedChart>
                            </ResponsiveContainer>
                        )}
                    </div>
                </section>

                <section className="grid gap-4 lg:grid-cols-2">
                    <article className="rounded-2xl border border-border bg-surface p-4 shadow-sm border-border bg-surface">
                        <h2 className="text-base font-semibold text-foreground">Expense by Category</h2>
                        <div className="mt-4 h-72">
                            <ResponsiveContainer width="100%" height="100%">
                                <PieChart>
                                    <Pie
                                        data={expenseData}
                                        dataKey="amount"
                                        nameKey="name"
                                        innerRadius={60}
                                        outerRadius={100}
                                        paddingAngle={2}
                                    >
                                        {expenseData.map((entry, index) => (
                                            <Cell key={`${entry.name}-${index}`} fill={EXPENSE_COLORS[index % EXPENSE_COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip content={<DonutTooltip />} />
                                </PieChart>
                            </ResponsiveContainer>
                        </div>
                    </article>

                    <article className="rounded-2xl border border-border bg-surface p-4 shadow-sm border-border bg-surface">
                        <h2 className="text-base font-semibold text-foreground">Income by Category</h2>
                        <div className="mt-4 h-72">
                            <ResponsiveContainer width="100%" height="100%">
                                <PieChart>
                                    <Pie
                                        data={incomeData}
                                        dataKey="amount"
                                        nameKey="name"
                                        innerRadius={60}
                                        outerRadius={100}
                                        paddingAngle={2}
                                    >
                                        {incomeData.map((entry, index) => (
                                            <Cell key={`${entry.name}-${index}`} fill={INCOME_COLORS[index % INCOME_COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip content={<DonutTooltip />} />
                                </PieChart>
                            </ResponsiveContainer>
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
            ? 'border-positive/20 bg-positive-muted border-positive/20 bg-positive-muted'
            : tone === 'expense'
                ? 'border-negative/20 bg-negative-muted border-negative/20 bg-negative-muted'
                : 'border-border bg-elevated border-border bg-elevated'

    return (
        <article className={`rounded-2xl border p-4 ${toneClass}`}>
            <p className="text-xs font-semibold uppercase tracking-wide text-muted">{label}</p>
            <p className="mt-2 text-2xl font-bold text-foreground">{value}</p>
        </article>
    )
}
