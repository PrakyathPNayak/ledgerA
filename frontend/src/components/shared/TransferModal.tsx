import { useState } from 'react'
import type { FormEvent } from 'react'

import { useAccounts } from '@/hooks/useAccounts'
import { useCategories } from '@/hooks/useCategories'
import { useTransfer } from '@/hooks/useTransfer'

interface TransferModalProps {
    isOpen: boolean
    onClose: () => void
}

function formatDate(d: Date): string {
    return d.toISOString().slice(0, 10)
}

function extractApiErrorMessage(error: unknown, fallback: string): string {
    const maybeMessage =
        (error as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message
    if (typeof maybeMessage === 'string' && maybeMessage.trim() !== '') {
        return maybeMessage
    }
    return fallback
}

export function TransferModal({ isOpen, onClose }: TransferModalProps) {
    const { data: accounts = [] } = useAccounts()
    const { data: categories = [] } = useCategories()
    const transfer = useTransfer()

    const [fromAccountId, setFromAccountId] = useState('')
    const [toAccountId, setToAccountId] = useState('')
    const [categoryId, setCategoryId] = useState('')
    const [amount, setAmount] = useState('')
    const [name, setName] = useState('Transfer')
    const [transactionDate, setTransactionDate] = useState(formatDate(new Date()))
    const [notes, setNotes] = useState('')
    const [error, setError] = useState<string | null>(null)

    function resetForm() {
        setFromAccountId('')
        setToAccountId('')
        setCategoryId('')
        setAmount('')
        setName('Transfer')
        setTransactionDate(formatDate(new Date()))
        setNotes('')
        setError(null)
    }

    function handleClose() {
        resetForm()
        onClose()
    }

    if (!isOpen) return null

    async function handleSubmit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault()
        const parsedAmount = Number(amount)
        if (!fromAccountId || !toAccountId || !categoryId || Number.isNaN(parsedAmount) || parsedAmount <= 0) return
        if (fromAccountId === toAccountId) {
            setError('Source and destination accounts must be different.')
            return
        }

        try {
            setError(null)
            await transfer.mutateAsync({
                from_account_id: fromAccountId,
                to_account_id: toAccountId,
                category_id: categoryId,
                amount: parsedAmount,
                transaction_date: transactionDate,
                name: name.trim() || 'Transfer',
                notes: notes.trim() || undefined,
            })
            handleClose()
        } catch (err) {
            setError(extractApiErrorMessage(err, 'Unable to process transfer.'))
        }
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
            <div className="w-full max-w-lg rounded-2xl border border-border bg-surface shadow-xl">
                <div className="flex items-center justify-between border-b border-border px-6 py-4">
                    <h2 className="text-xl font-semibold text-foreground">Transfer Between Accounts</h2>
                    <button type="button" onClick={handleClose} className="rounded-md px-3 py-1.5 text-sm text-muted hover:bg-surface-hover">
                        Close
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="space-y-5 px-6 py-5">
                    <div className="grid gap-4 md:grid-cols-2">
                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">From Account</span>
                            <select value={fromAccountId} onChange={(e) => setFromAccountId(e.target.value)} className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground" required>
                                <option value="">Select source</option>
                                {accounts.map((item) => (
                                    <option key={item.id} value={item.id}>{item.name}</option>
                                ))}
                            </select>
                        </label>

                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">To Account</span>
                            <select value={toAccountId} onChange={(e) => setToAccountId(e.target.value)} className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground" required>
                                <option value="">Select destination</option>
                                {accounts.filter((a) => a.id !== fromAccountId).map((item) => (
                                    <option key={item.id} value={item.id}>{item.name}</option>
                                ))}
                            </select>
                        </label>
                    </div>

                    <div className="grid gap-4 md:grid-cols-2">
                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Category</span>
                            <select value={categoryId} onChange={(e) => setCategoryId(e.target.value)} className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground" required>
                                <option value="">Select category</option>
                                {categories.map((item) => (
                                    <option key={item.id} value={item.id}>{item.name}</option>
                                ))}
                            </select>
                        </label>

                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Name</span>
                            <input value={name} onChange={(e) => setName(e.target.value)} className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground" required />
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

                    <label className="space-y-1">
                        <span className="text-sm font-medium text-secondary">Notes</span>
                        <textarea value={notes} onChange={(e) => setNotes(e.target.value)} rows={2} className="w-full rounded-lg border border-border bg-surface px-3 py-2 text-sm text-foreground" />
                    </label>

                    {error ? (
                        <p className="rounded-md bg-negative-muted px-3 py-2 text-sm text-negative">{error}</p>
                    ) : null}

                    <div className="flex justify-end gap-2 border-t border-border pt-4">
                        <button type="button" onClick={handleClose} className="rounded-lg border border-border px-4 py-2 text-sm text-secondary">
                            Cancel
                        </button>
                        <button type="submit" disabled={transfer.isPending} className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white disabled:opacity-60">
                            {transfer.isPending ? 'Transferring...' : 'Transfer'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    )
}
