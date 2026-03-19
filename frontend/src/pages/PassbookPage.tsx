import { useState } from 'react'
import { Link, useParams } from 'react-router-dom'

import { PageShell } from '@/components/layout/PageShell'
import { EditTransactionModal } from '@/components/shared/EditTransactionModal'
import { TransactionDetailModal } from '@/components/shared/TransactionDetailModal'
import { SortableTable } from '@/components/shared/SortableTable'
import { TransactionRow } from '@/components/shared/TransactionRow'
import { useAccounts } from '@/hooks/useAccounts'
import { useCategories } from '@/hooks/useCategories'
import { useTransactions, useDeleteTransaction } from '@/hooks/useTransactions'
import type { Transaction } from '@/types'

const columns = [
    { key: 'transaction_date', label: 'Date', sortable: true },
    { key: 'name', label: 'Name', sortable: true },
    { key: 'category', label: 'Category' },
    { key: 'amount', label: 'Amount', sortable: true, className: 'text-right' },
    { key: 'balance', label: 'Balance', className: 'text-right' },
]

function money(value: number): string {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR' }).format(value)
}

export function PassbookPage() {
    const { id = '' } = useParams()
    const { data: accounts = [] } = useAccounts()
    const { data: categories = [] } = useCategories()
    const { data: transactions = [] } = useTransactions({ account_id: id, passbook_mode: true, per_page: 200, page: 1 })
    const deleteMutation = useDeleteTransaction()

    const [selectedTx, setSelectedTx] = useState<Transaction | null>(null)
    const [editingTx, setEditingTx] = useState<Transaction | null>(null)

    const account = accounts.find((item) => item.id === id)

    const categoryMap = Object.fromEntries(categories.map((item) => [item.id, item.name]))
    const subcategoryMap = Object.fromEntries(
        categories.flatMap((category) => (category.subcategories ?? []).map((sub) => [sub.id, sub.name])),
    )

    // Compute running balance starting from opening balance
    const openingBalance = account?.opening_balance ?? 0
    const runningBalances: number[] = []
    let balance = openingBalance
    for (const tx of transactions) {
        balance += tx.amount
        runningBalances.push(balance)
    }

    return (
        <PageShell title="Passbook">
            <section className="space-y-4 rounded-2xl border border-border bg-surface p-4 shadow-sm">
                <div className="flex items-center justify-between">
                    <div>
                        <p className="text-sm text-muted">Account</p>
                        <h2 className="text-2xl font-bold text-foreground">{account?.name ?? 'Unknown account'}</h2>
                    </div>
                    <div className="text-right">
                        <p className="text-xs uppercase tracking-wide text-muted">Current Balance</p>
                        <p className="text-2xl font-bold text-foreground">{money(account?.current_balance ?? 0)}</p>
                    </div>
                </div>

                <div className="overflow-hidden rounded-2xl border border-border bg-surface shadow-sm">
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-border">
                            <thead className="bg-elevated">
                                <tr>
                                    {columns.map((col) => (
                                        <th key={col.key} className={`px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-muted ${col.className ?? ''}`}>
                                            {col.label}
                                        </th>
                                    ))}
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-border-subtle bg-surface">
                                {transactions.length > 0 ? (
                                    <>
                                        <tr className="bg-elevated/50">
                                            <td className="px-4 py-2 text-sm text-muted" colSpan={3}>Opening Balance</td>
                                            <td className="px-4 py-2 text-right text-sm font-semibold text-foreground" colSpan={2}>{money(openingBalance)}</td>
                                        </tr>
                                        {transactions.map((tx: Transaction, idx: number) => (
                                            <tr key={tx.id} onClick={() => setSelectedTx(tx)} className="cursor-pointer transition-colors hover:bg-elevated">
                                                <td className="px-4 py-3 text-sm text-foreground">{tx.transaction_date}</td>
                                                <td className="px-4 py-3 text-sm font-medium text-foreground">{tx.name}</td>
                                                <td className="px-4 py-3 text-sm text-secondary">
                                                    {categoryMap[tx.category_id] ?? 'Unknown'}
                                                    {subcategoryMap[tx.subcategory_id] ? <span className="text-muted"> / {subcategoryMap[tx.subcategory_id]}</span> : null}
                                                </td>
                                                <td className={`px-4 py-3 text-right text-sm font-semibold ${tx.amount >= 0 ? 'text-positive' : 'text-negative'}`}>
                                                    {money(tx.amount)}
                                                </td>
                                                <td className="px-4 py-3 text-right text-sm font-semibold text-foreground">
                                                    {money(runningBalances[idx])}
                                                </td>
                                            </tr>
                                        ))}
                                    </>
                                ) : null}
                            </tbody>
                        </table>
                    </div>
                    {transactions.length === 0 ? <p className="px-4 py-8 text-center text-sm text-muted">No passbook entries yet for this account.</p> : null}
                </div>

                <Link to="/accounts" className="inline-flex text-sm font-medium text-secondary underline-offset-2 hover:underline">
                    Back to accounts
                </Link>
            </section>

            <TransactionDetailModal
                transaction={selectedTx}
                accountName={account?.name}
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
