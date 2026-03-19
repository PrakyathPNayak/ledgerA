import { useEffect, useMemo, useState } from 'react'

import { PageShell } from '@/components/layout/PageShell'
import { useAccounts } from '@/hooks/useAccounts'
import { useCompare } from '@/hooks/useStats'
import type { ComparePeriodData, Period } from '@/types'

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
        return Array.from({ length: 30 }, (_, i) => {
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
            return { value: isoWeek(d), label: isoWeek(d) }
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

function formatMoney(value: number): string {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR' }).format(value)
}

function formatPct(value: number): string {
    if (!Number.isFinite(value)) return 'N/A'
    const sign = value > 0 ? '+' : ''
    return `${sign}${value.toFixed(1)}%`
}

function PeriodCard({ data, label }: { data: ComparePeriodData; label: string }) {
    return (
        <div className="rounded-2xl border border-border bg-surface p-5 shadow-sm">
            <h3 className="text-sm font-semibold uppercase tracking-wide text-muted">{label}</h3>
            <p className="mt-1 text-xs text-muted">{data.label}</p>

            <div className="mt-4 grid grid-cols-3 gap-3">
                <div>
                    <p className="text-xs text-muted">Income</p>
                    <p className="text-lg font-bold text-positive">{formatMoney(data.total_income)}</p>
                </div>
                <div>
                    <p className="text-xs text-muted">Expense</p>
                    <p className="text-lg font-bold text-negative">{formatMoney(data.total_expense)}</p>
                </div>
                <div>
                    <p className="text-xs text-muted">Net</p>
                    <p className={`text-lg font-bold ${data.net >= 0 ? 'text-positive' : 'text-negative'}`}>
                        {formatMoney(data.net)}
                    </p>
                </div>
            </div>

            {data.top_expense.length > 0 && (
                <div className="mt-4">
                    <p className="text-xs font-semibold uppercase tracking-wide text-muted">Top Expenses</p>
                    <ul className="mt-2 space-y-1">
                        {data.top_expense.map((item, i) => (
                            <li key={i} className="flex items-center justify-between text-sm">
                                <span className="text-secondary">{item.subcategory || item.category}</span>
                                <span className="font-medium text-foreground">{formatMoney(Math.abs(item.amount))}</span>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    )
}

export function ComparePage() {
    const [period, setPeriod] = useState<Period>('month')
    const [value1, setValue1] = useState('')
    const [value2, setValue2] = useState('')
    const [accountId, setAccountId] = useState('')

    const { data: accounts = [] } = useAccounts()
    const valueOptions = useMemo(() => getValueOptions(period), [period])

    useEffect(() => {
        setValue1(valueOptions[0]?.value ?? '')
        setValue2(valueOptions[1]?.value ?? '')
    }, [period, valueOptions])

    const { data: compare, isLoading } = useCompare({
        period,
        value1,
        value2,
        account_id: accountId || undefined,
    })

    return (
        <PageShell title="Compare">
            <div className="space-y-6">
                <section className="rounded-2xl border border-border bg-surface p-4 shadow-sm">
                    <div className="grid gap-3 md:grid-cols-4">
                        <label className="space-y-1">
                            <span className="text-xs font-semibold uppercase tracking-wide text-muted">Period</span>
                            <select
                                value={period}
                                onChange={(e) => setPeriod(e.target.value as Period)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm border-border bg-elevated text-foreground"
                            >
                                <option value="day">Day</option>
                                <option value="week">Week</option>
                                <option value="month">Month</option>
                                <option value="year">Year</option>
                            </select>
                        </label>

                        <label className="space-y-1">
                            <span className="text-xs font-semibold uppercase tracking-wide text-muted">Period 1</span>
                            <select
                                value={value1}
                                onChange={(e) => setValue1(e.target.value)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm border-border bg-elevated text-foreground"
                            >
                                {valueOptions.map((opt) => (
                                    <option key={opt.value} value={opt.value}>{opt.label}</option>
                                ))}
                            </select>
                        </label>

                        <label className="space-y-1">
                            <span className="text-xs font-semibold uppercase tracking-wide text-muted">Period 2</span>
                            <select
                                value={value2}
                                onChange={(e) => setValue2(e.target.value)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm border-border bg-elevated text-foreground"
                            >
                                {valueOptions.map((opt) => (
                                    <option key={opt.value} value={opt.value}>{opt.label}</option>
                                ))}
                            </select>
                        </label>

                        <label className="space-y-1">
                            <span className="text-xs font-semibold uppercase tracking-wide text-muted">Account</span>
                            <select
                                value={accountId}
                                onChange={(e) => setAccountId(e.target.value)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm border-border bg-elevated text-foreground"
                            >
                                <option value="">All accounts</option>
                                {accounts.map((account) => (
                                    <option key={account.id} value={account.id}>{account.name}</option>
                                ))}
                            </select>
                        </label>
                    </div>
                </section>

                {isLoading && <p className="text-sm text-muted">Loading comparison...</p>}

                {compare && (
                    <>
                        <section className="grid gap-4 md:grid-cols-3">
                            <article className={`rounded-2xl border p-4 ${compare.income_change_pct >= 0 ? 'border-positive/20 bg-positive-muted' : 'border-negative/20 bg-negative-muted'}`}>
                                <p className="text-xs font-semibold uppercase tracking-wide text-muted">Income Change</p>
                                <p className="mt-2 text-2xl font-bold text-foreground">{formatPct(compare.income_change_pct)}</p>
                            </article>
                            <article className={`rounded-2xl border p-4 ${compare.expense_change_pct <= 0 ? 'border-positive/20 bg-positive-muted' : 'border-negative/20 bg-negative-muted'}`}>
                                <p className="text-xs font-semibold uppercase tracking-wide text-muted">Expense Change</p>
                                <p className="mt-2 text-2xl font-bold text-foreground">{formatPct(compare.expense_change_pct)}</p>
                            </article>
                            <article className={`rounded-2xl border p-4 ${compare.net_change >= 0 ? 'border-positive/20 bg-positive-muted' : 'border-negative/20 bg-negative-muted'}`}>
                                <p className="text-xs font-semibold uppercase tracking-wide text-muted">Net Change</p>
                                <p className="mt-2 text-2xl font-bold text-foreground">{formatMoney(compare.net_change)}</p>
                            </article>
                        </section>

                        <section className="grid gap-4 lg:grid-cols-2">
                            <PeriodCard data={compare.period1} label="Period 1" />
                            <PeriodCard data={compare.period2} label="Period 2" />
                        </section>
                    </>
                )}
            </div>
        </PageShell>
    )
}
