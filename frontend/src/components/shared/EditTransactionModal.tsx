import { useEffect, useMemo, useState } from 'react'
import type { FormEvent } from 'react'

import { useAccounts } from '@/hooks/useAccounts'
import { useCategories } from '@/hooks/useCategories'
import { useUpdateTransaction } from '@/hooks/useTransactions'
import type { Transaction } from '@/types'

const QUICK_AMOUNTS = [50, 100, 250, 500, 1000]

type TransactionKind = 'expense' | 'income'

interface EditTransactionModalProps {
    transaction: Transaction | null
    onClose: () => void
}

function extractApiErrorMessage(error: unknown, fallback: string): string {
    const maybeMessage =
        (error as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message
    if (typeof maybeMessage === 'string' && maybeMessage.trim() !== '') {
        return maybeMessage
    }
    return fallback
}

export function EditTransactionModal({ transaction, onClose }: EditTransactionModalProps) {
    const { data: accounts = [] } = useAccounts()
    const { data: categories = [] } = useCategories()
    const updateTransaction = useUpdateTransaction()

    const [kind, setKind] = useState<TransactionKind>('expense')
    const [accountId, setAccountId] = useState('')
    const [categoryId, setCategoryId] = useState('')
    const [subcategoryId, setSubcategoryId] = useState('')
    const [name, setName] = useState('')
    const [amount, setAmount] = useState('')
    const [transactionDate, setTransactionDate] = useState('')
    const [notes, setNotes] = useState('')
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        if (transaction) {
            setKind(transaction.amount >= 0 ? 'income' : 'expense')
            setAccountId(transaction.account_id)
            setCategoryId(transaction.category_id)
            setSubcategoryId(transaction.subcategory_id ?? '')
            setName(transaction.name)
            setAmount(String(Math.abs(transaction.amount)))
            setTransactionDate(transaction.transaction_date)
            setNotes(transaction.notes ?? '')
            setError(null)
        }
    }, [transaction])

    const selectedCategory = useMemo(
        () => categories.find((cat) => cat.id === categoryId),
        [categories, categoryId],
    )

    if (!transaction) return null

    async function handleSubmit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault()
        if (!transaction) return

        const parsedAmount = Number(amount)
        if (Number.isNaN(parsedAmount) || parsedAmount <= 0) return

        const signedAmount = kind === 'expense' ? -Math.abs(parsedAmount) : Math.abs(parsedAmount)

        try {
            setError(null)
            await updateTransaction.mutateAsync({
                id: transaction.id,
                account_id: accountId || undefined,
                category_id: categoryId || undefined,
                subcategory_id: subcategoryId || undefined,
                name: name.trim() || undefined,
                amount: signedAmount,
                transaction_date: transactionDate || undefined,
                notes: notes.trim() || undefined,
            })
            onClose()
        } catch (err) {
            setError(extractApiErrorMessage(err, 'Unable to update transaction.'))
        }
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
            <div className="w-full max-w-2xl rounded-2xl border border-border bg-surface shadow-xl">
                <div className="flex items-center justify-between border-b border-border px-6 py-4">
                    <h2 className="text-xl font-semibold text-foreground">Edit Transaction</h2>
                    <button type="button" onClick={onClose} className="rounded-md px-3 py-1.5 text-sm text-muted hover:bg-surface-hover">
                        Close
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="space-y-5 px-6 py-5">
                    <div className="flex gap-2">
                        <button
                            type="button"
                            onClick={() => setKind('expense')}
                            className={`rounded-full px-4 py-2 text-sm font-medium ${kind === 'expense' ? 'bg-rose-100 text-negative' : 'bg-elevated text-secondary'}`}
                        >
                            Expense
                        </button>
                        <button
                            type="button"
                            onClick={() => setKind('income')}
                            className={`rounded-full px-4 py-2 text-sm font-medium ${kind === 'income' ? 'bg-positive-muted text-positive' : 'bg-elevated text-secondary'}`}
                        >
                            Income
                        </button>
                    </div>

                    <div className="grid gap-4 md:grid-cols-2">
                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Account</span>
                            <select value={accountId} onChange={(e) => setAccountId(e.target.value)} className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground" required>
                                <option value="">Select account</option>
                                {accounts.map((item) => (
                                    <option key={item.id} value={item.id}>{item.name}</option>
                                ))}
                            </select>
                        </label>

                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Category</span>
                            <select
                                value={categoryId}
                                onChange={(e) => { setCategoryId(e.target.value); setSubcategoryId('') }}
                                className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground" required
                            >
                                <option value="">Select category</option>
                                {categories.map((item) => (
                                    <option key={item.id} value={item.id}>{item.name}</option>
                                ))}
                            </select>
                        </label>
                    </div>

                    <div className="grid gap-4 md:grid-cols-2">
                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Subcategory</span>
                            <select value={subcategoryId} onChange={(e) => setSubcategoryId(e.target.value)} className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground">
                                <option value="">Auto / none</option>
                                {(selectedCategory?.subcategories ?? []).map((item) => (
                                    <option key={item.id} value={item.id}>{item.name}</option>
                                ))}
                            </select>
                        </label>

                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Name</span>
                            <input value={name} onChange={(e) => setName(e.target.value)} className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground" placeholder="e.g. Groceries" required />
                        </label>
                    </div>

                    <div className="grid gap-4 md:grid-cols-2">
                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Amount</span>
                            <input value={amount} onChange={(e) => setAmount(e.target.value)} inputMode="decimal" className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground" placeholder="0.00" required />
                        </label>

                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Date</span>
                            <input type="date" value={transactionDate} onChange={(e) => setTransactionDate(e.target.value)} className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground" required />
                        </label>
                    </div>

                    <div className="space-y-2">
                        <p className="text-xs font-semibold uppercase tracking-wide text-muted">Quick Fill Amount</p>
                        <div className="flex flex-wrap gap-2">
                            {QUICK_AMOUNTS.map((value) => (
                                <button key={value} type="button" onClick={() => setAmount(String(value))} className="rounded-full border border-border px-3 py-1 text-sm text-secondary">
                                    {value}
                                </button>
                            ))}
                        </div>
                    </div>

                    <label className="space-y-1">
                        <span className="text-sm font-medium text-secondary">Notes</span>
                        <textarea value={notes} onChange={(e) => setNotes(e.target.value)} rows={3} className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground" />
                    </label>

                    {error ? (
                        <p className="rounded-md bg-negative-muted px-3 py-2 text-sm text-negative">{error}</p>
                    ) : null}

                    <div className="flex justify-end gap-2 border-t border-border pt-4">
                        <button type="button" onClick={onClose} className="rounded-lg border border-border px-4 py-2 text-sm text-secondary">
                            Cancel
                        </button>
                        <button type="submit" disabled={updateTransaction.isPending} className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white disabled:opacity-60">
                            {updateTransaction.isPending ? 'Saving...' : 'Save Changes'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    )
}
