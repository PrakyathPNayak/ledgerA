import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'

import { PageShell } from '@/components/layout/PageShell'
import { AddTransactionModal } from '@/components/shared/AddTransactionModal'
import { EditTransactionModal } from '@/components/shared/EditTransactionModal'
import { TransactionDetailModal } from '@/components/shared/TransactionDetailModal'
import { TransactionRow } from '@/components/shared/TransactionRow'
import { SortableTable } from '@/components/shared/SortableTable'
import { useAccounts } from '@/hooks/useAccounts'
import { useCategories } from '@/hooks/useCategories'
import { useStats } from '@/hooks/useStats'
import { useTransactions, useDeleteTransaction } from '@/hooks/useTransactions'
import { useBudgetProgress } from '@/hooks/useBudgets'
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
    const [selectedTx, setSelectedTx] = useState<Transaction | null>(null)
    const [editingTx, setEditingTx] = useState<Transaction | null>(null)
    const [sortBy, setSortBy] = useState<'transaction_date' | 'amount' | 'name'>('transaction_date')
    const [sortDir, setSortDir] = useState<'asc' | 'desc'>('desc')

    const { data: transactions = [], refetch } = useTransactions({ sort_by: sortBy, sort_dir: sortDir, per_page: 10, page: 1 })
    const { data: accounts = [] } = useAccounts()
    const { data: categories = [] } = useCategories()
    const { data: stats } = useStats({ period: 'month', value: new Date().toISOString().slice(0, 7) })
    const { data: budgetProgressRaw = [] } = useBudgetProgress()
    const budgetProgress = Array.isArray(budgetProgressRaw) ? budgetProgressRaw : []
    const deleteMutation = useDeleteTransaction()

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

                {budgetProgress.length > 0 && (
                    <section className="rounded-2xl border border-border bg-surface p-4 shadow-sm">
                        <div className="mb-3 flex items-center justify-between">
                            <h2 className="text-base font-semibold text-foreground">Monthly Budgets</h2>
                            <Link to="/budgets" className="text-xs text-accent hover:underline">View all →</Link>
                        </div>
                        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                            {budgetProgress.map((b) => {
                                const cat = categories.find((c) => c.id === b.category_id)?.name ?? '—'
                                const over = b.percent > 100
                                return (
                                    <div key={b.id} className="rounded-xl border border-border bg-elevated p-3">
                                        <div className="flex items-center justify-between mb-1">
                                            <span className="text-xs font-medium text-foreground">{cat}</span>
                                            <span className={`text-xs font-semibold ${over ? 'text-negative' : 'text-muted'}`}>
                                                {b.percent.toFixed(0)}%
                                            </span>
                                        </div>
                                        <div className="h-1.5 w-full rounded-full bg-surface">
                                            <div
                                                className={`h-1.5 rounded-full ${over ? 'bg-[var(--negative)]' : 'bg-[var(--accent)]'}`}
                                                style={{ width: `${Math.min(b.percent, 100)}%` }}
                                            />
                                        </div>
                                        <p className="mt-1 text-xs text-muted">{money(b.spent)} of {money(b.amount)}</p>
                                    </div>
                                )
                            })}
                        </div>
                    </section>
                )}

                <section className="rounded-2xl border border-border bg-surface p-4 shadow-sm">
                    <div className="mb-4 flex items-center justify-between">
                        <h2 className="text-lg font-semibold text-foreground">Recent Transactions</h2>
                        <button
                            type="button"
                            onClick={() => setIsAddOpen(true)}
                            className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white"
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
                                onClick={setSelectedTx}
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

            <TransactionDetailModal
                transaction={selectedTx}
                accountName={selectedTx ? accountMap[selectedTx.account_id] : undefined}
                categoryName={selectedTx ? categoryMap[selectedTx.category_id] : undefined}
                subcategoryName={selectedTx ? subcategoryMap[selectedTx.subcategory_id] : undefined}
                onClose={() => setSelectedTx(null)}
                onEdit={(tx) => { setSelectedTx(null); setEditingTx(tx) }}
                onDelete={(tx) => { setSelectedTx(null); deleteMutation.mutate(tx.id) }}
            />

            <EditTransactionModal
                transaction={editingTx}
                onClose={() => setEditingTx(null)}
            />
        </PageShell>
    )
}

function StatCard({ label, value, tone }: { label: string; value: string; tone: 'income' | 'expense' | 'neutral' }) {
    const toneClass =
        tone === 'income'
            ? 'border-positive/20 bg-positive-muted'
            : tone === 'expense'
                ? 'border-negative/20 bg-negative-muted'
                : 'border-border bg-elevated'

    return (
        <article className={`rounded-2xl border p-4 ${toneClass}`}>
            <p className="text-xs font-semibold uppercase tracking-wide text-muted">{label}</p>
            <p className="mt-2 text-2xl font-bold text-foreground">{value}</p>
        </article>
    )
}
