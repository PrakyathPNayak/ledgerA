import { useMemo, useState } from 'react'

import { PageShell } from '@/components/layout/PageShell'
import { AddTransactionModal } from '@/components/shared/AddTransactionModal'
import { TransactionRow } from '@/components/shared/TransactionRow'
import { SortableTable } from '@/components/shared/SortableTable'
import { useAccounts } from '@/hooks/useAccounts'
import { useCategories } from '@/hooks/useCategories'
import { useStats } from '@/hooks/useStats'
import { useTransactions } from '@/hooks/useTransactions'
import type { Transaction } from '@/types'

const tableColumns = [
    { key: 'transaction_date', label: 'Date', sortable: true },
    { key: 'name', label: 'Name', sortable: true },
    { key: 'account', label: 'Account' },
    { key: 'category', label: 'Category' },
    { key: 'amount', label: 'Amount', sortable: true, className: 'text-right' },
]

function money(value: number): string {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR' }).format(value)
}

export function DashboardPage() {
    const [isAddOpen, setIsAddOpen] = useState(false)
    const [sortBy, setSortBy] = useState<'transaction_date' | 'amount' | 'name'>('transaction_date')
    const [sortDir, setSortDir] = useState<'asc' | 'desc'>('desc')

    const { data: transactions = [], refetch } = useTransactions({ sort_by: sortBy, sort_dir: sortDir, per_page: 10, page: 1 })
    const { data: accounts = [] } = useAccounts()
    const { data: categories = [] } = useCategories()
    const { data: stats } = useStats({ period: 'month', value: new Date().toISOString().slice(0, 7) })

    const accountMap = useMemo(() => Object.fromEntries(accounts.map((item) => [item.id, item.name])), [accounts])
    const categoryMap = useMemo(() => Object.fromEntries(categories.map((item) => [item.id, item.name])), [categories])
    const subcategoryMap = useMemo(
        () =>
            Object.fromEntries(
                categories.flatMap((category) => (category.subcategories ?? []).map((sub) => [sub.id, sub.name])),
            ),
        [categories],
    )

    function toggleSort(column: string) {
        if (column !== 'transaction_date' && column !== 'amount' && column !== 'name') {
            return
        }
        if (sortBy === column) {
            setSortDir((prev) => (prev === 'asc' ? 'desc' : 'asc'))
            return
        }
        setSortBy(column)
        setSortDir('desc')
    }

    return (
        <PageShell title="Dashboard">
            <div className="space-y-6">
                <section className="grid gap-4 md:grid-cols-3">
                    <StatCard label="Income" value={money(stats?.total_income ?? 0)} tone="income" />
                    <StatCard label="Expense" value={money(stats?.total_expense ?? 0)} tone="expense" />
                    <StatCard label="Net" value={money(stats?.net ?? 0)} tone="neutral" />
                </section>

                <section className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
                    <div className="mb-4 flex items-center justify-between">
                        <h2 className="text-lg font-semibold text-slate-900">Recent Transactions</h2>
                        <button
                            type="button"
                            onClick={() => setIsAddOpen(true)}
                            className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-medium text-white"
                        >
                            Add Transaction
                        </button>
                    </div>

                    <SortableTable
                        columns={tableColumns}
                        sortBy={sortBy}
                        sortDir={sortDir}
                        onSort={toggleSort}
                        hasRows={transactions.length > 0}
                        emptyMessage="No transactions yet. Add one to begin tracking."
                    >
                        {transactions.map((transaction: Transaction) => (
                            <TransactionRow
                                key={transaction.id}
                                transaction={transaction}
                                accountName={accountMap[transaction.account_id]}
                                categoryName={categoryMap[transaction.category_id]}
                                subcategoryName={subcategoryMap[transaction.subcategory_id]}
                            />
                        ))}
                    </SortableTable>
                </section>
            </div>

            <AddTransactionModal
                isOpen={isAddOpen}
                onClose={() => setIsAddOpen(false)}
                onCreated={() => {
                    void refetch()
                }}
            />
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
