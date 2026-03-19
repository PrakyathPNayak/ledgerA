import { useState } from 'react'
import type { FormEvent } from 'react'

import { PageShell } from '@/components/layout/PageShell'
import { useCategories } from '@/hooks/useCategories'
import {
    useBudgetProgress,
    useCreateBudget,
    useDeleteBudget,
} from '@/hooks/useBudgets'
import type { BudgetProgress } from '@/types'

function money(value: number): string {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR' }).format(value)
}

export function BudgetsPage() {
    const { data: progressItems = [], isLoading } = useBudgetProgress()
    const budgets: BudgetProgress[] = Array.isArray(progressItems) ? progressItems : []
    const { data: categories = [] } = useCategories()
    const createBudget = useCreateBudget()
    const deleteBudget = useDeleteBudget()

    const [showForm, setShowForm] = useState(false)
    const [categoryId, setCategoryId] = useState('')
    const [amount, setAmount] = useState('')
    const [period, setPeriod] = useState('monthly')

    async function handleCreate(e: FormEvent<HTMLFormElement>) {
        e.preventDefault()
        if (!categoryId || !amount) return
        await createBudget.mutateAsync({ category_id: categoryId, amount: Number(amount), period })
        setShowForm(false)
        setCategoryId('')
        setAmount('')
    }

    return (
        <PageShell title="Budgets">
            <section className="rounded-2xl border border-border bg-surface p-4 shadow-sm">
                <div className="mb-4 flex items-center justify-between">
                    <h2 className="text-lg font-semibold text-foreground">Category Budgets</h2>
                    <button
                        onClick={() => setShowForm(!showForm)}
                        className="rounded-lg bg-accent px-3 py-1.5 text-sm font-medium text-white"
                    >
                        {showForm ? 'Cancel' : '+ Set Budget'}
                    </button>
                </div>

                {showForm && (
                    <form className="mb-4 grid gap-3 rounded-xl border border-border bg-elevated p-4 sm:grid-cols-3" onSubmit={handleCreate}>
                        <select value={categoryId} onChange={(e) => setCategoryId(e.target.value)} className="rounded-lg border border-border px-3 py-2 text-sm" required>
                            <option value="">Select Category</option>
                            {categories.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
                        </select>
                        <input
                            value={amount}
                            onChange={(e) => setAmount(e.target.value)}
                            inputMode="decimal"
                            placeholder="Budget Amount"
                            className="rounded-lg border border-border px-3 py-2 text-sm"
                            required
                        />
                        <select value={period} onChange={(e) => setPeriod(e.target.value)} className="rounded-lg border border-border px-3 py-2 text-sm">
                            <option value="monthly">Monthly</option>
                            <option value="yearly">Yearly</option>
                        </select>
                        <button type="submit" disabled={createBudget.isPending} className="col-span-full rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white disabled:opacity-60">
                            {createBudget.isPending ? 'Saving...' : 'Create Budget'}
                        </button>
                    </form>
                )}

                {isLoading ? (
                    <p className="text-sm text-muted">Loading...</p>
                ) : budgets.length === 0 ? (
                    <p className="text-sm text-muted">No budgets set. Create one to track spending limits.</p>
                ) : (
                    <div className="space-y-3">
                        {budgets.map((b) => {
                            const categoryName = categories.find((c) => c.id === b.category_id)?.name ?? '—'
                            const overBudget = b.percent > 100
                            return (
                                <article key={b.id} className="rounded-xl border border-border bg-elevated p-4">
                                    <div className="flex items-center justify-between mb-2">
                                        <div>
                                            <span className="text-sm font-semibold text-foreground">{categoryName}</span>
                                            <span className="ml-2 text-xs text-muted capitalize">{b.period}</span>
                                        </div>
                                        <div className="flex items-center gap-3">
                                            <span className={`text-sm font-medium ${overBudget ? 'text-negative' : 'text-positive'}`}>
                                                {money(b.spent)} / {money(b.amount)}
                                            </span>
                                            <button onClick={() => deleteBudget.mutate(b.id)} className="text-xs text-negative hover:underline">
                                                Remove
                                            </button>
                                        </div>
                                    </div>
                                    <div className="h-2 w-full rounded-full bg-surface">
                                        <div
                                            className={`h-2 rounded-full transition-all ${overBudget ? 'bg-[var(--negative)]' : 'bg-[var(--accent)]'}`}
                                            style={{ width: `${Math.min(b.percent, 100)}%` }}
                                        />
                                    </div>
                                    <div className="mt-1 flex justify-between text-xs text-muted">
                                        <span>{b.percent.toFixed(0)}% used</span>
                                        <span>{money(b.remaining)} remaining</span>
                                    </div>
                                </article>
                            )
                        })}
                    </div>
                )}
            </section>
        </PageShell>
    )
}
