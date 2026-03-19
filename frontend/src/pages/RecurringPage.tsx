import { useState } from 'react'
import type { FormEvent } from 'react'

import { PageShell } from '@/components/layout/PageShell'
import { useAccounts } from '@/hooks/useAccounts'
import { useCategories } from '@/hooks/useCategories'
import {
    useRecurringTransactions,
    useCreateRecurring,
    useDeleteRecurring,
    useUpdateRecurring,
} from '@/hooks/useRecurring'
import type { RecurringTransaction } from '@/types'

function money(value: number): string {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR' }).format(value)
}

export function RecurringPage() {
    const { data: result, isLoading } = useRecurringTransactions()
    const items: RecurringTransaction[] = Array.isArray(result) ? result : (result as any)?.data ?? []
    const { data: accounts = [] } = useAccounts()
    const { data: categories = [] } = useCategories()

    const createRecurring = useCreateRecurring()
    const deleteRecurring = useDeleteRecurring()
    const updateRecurring = useUpdateRecurring()

    const [showForm, setShowForm] = useState(false)
    const [name, setName] = useState('')
    const [amount, setAmount] = useState('')
    const [accountId, setAccountId] = useState('')
    const [categoryId, setCategoryId] = useState('')
    const [frequency, setFrequency] = useState('monthly')
    const [startDate, setStartDate] = useState(new Date().toISOString().slice(0, 10))
    const [endDate, setEndDate] = useState('')
    const [notes, setNotes] = useState('')

    async function handleCreate(e: FormEvent<HTMLFormElement>) {
        e.preventDefault()
        if (!name.trim() || !amount || !accountId || !categoryId) return

        await createRecurring.mutateAsync({
            name: name.trim(),
            amount: Number(amount),
            account_id: accountId,
            category_id: categoryId,
            frequency,
            start_date: startDate,
            end_date: endDate || undefined,
            notes: notes || undefined,
        })

        setShowForm(false)
        setName('')
        setAmount('')
        setAccountId('')
        setCategoryId('')
        setFrequency('monthly')
        setNotes('')
    }

    return (
        <PageShell title="Recurring Transactions">
            <section className="rounded-2xl border border-border bg-surface p-4 shadow-sm">
                <div className="mb-4 flex items-center justify-between">
                    <h2 className="text-lg font-semibold text-foreground">Scheduled Transactions</h2>
                    <div className="flex items-center gap-3">
                        <p className="text-sm text-muted">{items.length} recurring</p>
                        <button
                            onClick={() => setShowForm(!showForm)}
                            className="rounded-lg bg-accent px-3 py-1.5 text-sm font-medium text-white"
                        >
                            {showForm ? 'Cancel' : '+ New'}
                        </button>
                    </div>
                </div>

                {showForm && (
                    <form className="mb-4 grid gap-3 rounded-xl border border-border bg-elevated p-4 sm:grid-cols-2" onSubmit={handleCreate}>
                        <input
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="Name"
                            className="rounded-lg border border-border px-3 py-2 text-sm"
                            required
                        />
                        <input
                            value={amount}
                            onChange={(e) => setAmount(e.target.value)}
                            inputMode="decimal"
                            placeholder="Amount (negative = expense)"
                            className="rounded-lg border border-border px-3 py-2 text-sm"
                            required
                        />
                        <select value={accountId} onChange={(e) => setAccountId(e.target.value)} className="rounded-lg border border-border px-3 py-2 text-sm" required>
                            <option value="">Select Account</option>
                            {accounts.map((a) => <option key={a.id} value={a.id}>{a.name}</option>)}
                        </select>
                        <select value={categoryId} onChange={(e) => setCategoryId(e.target.value)} className="rounded-lg border border-border px-3 py-2 text-sm" required>
                            <option value="">Select Category</option>
                            {categories.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
                        </select>
                        <select value={frequency} onChange={(e) => setFrequency(e.target.value)} className="rounded-lg border border-border px-3 py-2 text-sm">
                            <option value="daily">Daily</option>
                            <option value="weekly">Weekly</option>
                            <option value="monthly">Monthly</option>
                            <option value="yearly">Yearly</option>
                        </select>
                        <input value={startDate} onChange={(e) => setStartDate(e.target.value)} type="date" className="rounded-lg border border-border px-3 py-2 text-sm" required />
                        <input value={endDate} onChange={(e) => setEndDate(e.target.value)} type="date" placeholder="End date (optional)" className="rounded-lg border border-border px-3 py-2 text-sm" />
                        <input value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Notes (optional)" className="rounded-lg border border-border px-3 py-2 text-sm" />
                        <button type="submit" disabled={createRecurring.isPending} className="col-span-full rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white disabled:opacity-60">
                            {createRecurring.isPending ? 'Creating...' : 'Create Recurring Transaction'}
                        </button>
                    </form>
                )}

                {isLoading ? (
                    <p className="text-sm text-muted">Loading...</p>
                ) : items.length === 0 ? (
                    <p className="text-sm text-muted">No recurring transactions yet.</p>
                ) : (
                    <div className="space-y-2">
                        {items.map((item) => (
                            <RecurringCard
                                key={item.id}
                                item={item}
                                accountName={accounts.find((a) => a.id === item.account_id)?.name ?? '—'}
                                categoryName={categories.find((c) => c.id === item.category_id)?.name ?? '—'}
                                onToggle={() => updateRecurring.mutate({ id: item.id, is_active: !item.is_active })}
                                onDelete={() => deleteRecurring.mutate(item.id)}
                            />
                        ))}
                    </div>
                )}
            </section>
        </PageShell>
    )
}

function RecurringCard({ item, accountName, categoryName, onToggle, onDelete }: {
    item: RecurringTransaction
    accountName: string
    categoryName: string
    onToggle: () => void
    onDelete: () => void
}) {
    const isExpense = item.amount < 0
    return (
        <article className="rounded-xl border border-border bg-elevated p-3">
            <div className="flex items-start justify-between gap-3">
                <div className="flex-1">
                    <div className="flex items-center gap-2">
                        <span className={`text-sm font-semibold ${isExpense ? 'text-negative' : 'text-positive'}`}>
                            {money(item.amount)}
                        </span>
                        <span className="text-sm font-medium text-foreground">{item.name}</span>
                        <span className={`rounded-full px-2 py-0.5 text-xs ${item.is_active ? 'bg-positive-muted text-positive' : 'bg-negative-muted text-negative'}`}>
                            {item.is_active ? 'Active' : 'Paused'}
                        </span>
                    </div>
                    <p className="mt-1 text-xs text-muted">
                        {accountName} • {categoryName} • {item.frequency} • Next: {item.next_due_date}
                    </p>
                    {item.end_date && <p className="text-xs text-muted">Ends: {item.end_date}</p>}
                </div>
                <div className="flex gap-2">
                    <button onClick={onToggle} className="rounded-lg border border-border px-2 py-1 text-xs text-secondary hover:bg-surface-hover">
                        {item.is_active ? 'Pause' : 'Resume'}
                    </button>
                    <button onClick={onDelete} className="rounded-lg border border-border px-2 py-1 text-xs text-negative hover:bg-negative-muted">
                        Delete
                    </button>
                </div>
            </div>
        </article>
    )
}
