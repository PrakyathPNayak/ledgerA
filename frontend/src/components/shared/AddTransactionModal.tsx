import { useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'

import { useAccounts, useCreateAccount } from '@/hooks/useAccounts'
import { useCategories, useCreateCategory, useCreateSubcategory } from '@/hooks/useCategories'
import api from '@/lib/api'

const QUICK_AMOUNTS = [50, 100, 250, 500, 1000]

type TransactionKind = 'expense' | 'income'

interface AddTransactionModalProps {
    isOpen: boolean
    onClose: () => void
    onCreated?: () => void
}

interface CreateTransactionPayload {
    account_id: string
    category_id: string
    subcategory_id?: string
    name: string
    amount: number
    transaction_date: string
    notes?: string
    is_scheduled: boolean
}

function formatDate(d: Date): string {
    return d.toISOString().slice(0, 10)
}

function normalizeTransactionDate(input: string): string {
    if (/^\d{4}-\d{2}-\d{2}$/.test(input)) {
        return input
    }

    const ddmmyyyy = input.match(/^(\d{2})\/(\d{2})\/(\d{4})$/)
    if (ddmmyyyy) {
        const [, dd, mm, yyyy] = ddmmyyyy
        return `${yyyy}-${mm}-${dd}`
    }

    return input
}

function extractApiErrorMessage(error: unknown, fallback: string): string {
    const maybeMessage =
        (error as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message

    if (typeof maybeMessage === 'string' && maybeMessage.trim() !== '') {
        return maybeMessage
    }

    return fallback
}

export function AddTransactionModal({ isOpen, onClose, onCreated }: AddTransactionModalProps) {
    const queryClient = useQueryClient()
    const { data: accounts = [] } = useAccounts()
    const { data: categories = [] } = useCategories()

    const [kind, setKind] = useState<TransactionKind>('expense')
    const [accountId, setAccountId] = useState('')
    const [categoryId, setCategoryId] = useState('')
    const [subcategoryId, setSubcategoryId] = useState('')
    const [name, setName] = useState('')
    const [amount, setAmount] = useState('')
    const [transactionDate, setTransactionDate] = useState(formatDate(new Date()))
    const [notes, setNotes] = useState('')
    const [isScheduled, setIsScheduled] = useState(false)

    const [newAccountName, setNewAccountName] = useState('')
    const [newAccountOpeningBalance, setNewAccountOpeningBalance] = useState('0')
    const [newCategoryName, setNewCategoryName] = useState('')
    const [newSubcategoryName, setNewSubcategoryName] = useState('')
    const [inlineCreateError, setInlineCreateError] = useState<string | null>(null)
    const [transactionError, setTransactionError] = useState<string | null>(null)

    const createAccount = useCreateAccount()
    const createCategory = useCreateCategory()
    const createSubcategory = useCreateSubcategory(categoryId || '__none__')

    const selectedCategory = useMemo(
        () => categories.find((cat) => cat.id === categoryId),
        [categories, categoryId],
    )

    function resetForm() {
        setKind('expense')
        setAccountId('')
        setCategoryId('')
        setSubcategoryId('')
        setName('')
        setAmount('')
        setTransactionDate(formatDate(new Date()))
        setNotes('')
        setIsScheduled(false)
        setNewAccountName('')
        setNewAccountOpeningBalance('0')
        setNewCategoryName('')
        setNewSubcategoryName('')
        setInlineCreateError(null)
        setTransactionError(null)
    }

    function handleClose() {
        resetForm()
        onClose()
    }

    const createTransaction = useMutation({
        mutationFn: (payload: CreateTransactionPayload) => api.post('/transactions', payload),
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: ['transactions'] })
            setTransactionError(null)
            onCreated?.()
            handleClose()
        },
    })

    if (!isOpen) {
        return null
    }

    async function handleCreateAccount() {
        const trimmed = newAccountName.trim()
        if (!trimmed) return

        const parsedOpeningBalance = Number(newAccountOpeningBalance)
        if (!Number.isFinite(parsedOpeningBalance)) {
            setInlineCreateError('Opening balance must be a valid number.')
            return
        }

        try {
            setInlineCreateError(null)
            const created = (await createAccount.mutateAsync({
                name: trimmed,
                opening_balance: parsedOpeningBalance,
            })) as { id?: string }
            if (created?.id) {
                setAccountId(created.id)
            }
            setNewAccountName('')
            setNewAccountOpeningBalance('0')
        } catch (error) {
            setInlineCreateError(extractApiErrorMessage(error, 'Unable to add account right now.'))
        }
    }

    async function handleCreateCategory() {
        const trimmed = newCategoryName.trim()
        if (!trimmed) return

        try {
            setInlineCreateError(null)
            const created = (await createCategory.mutateAsync({ name: trimmed })) as { id?: string }
            if (created?.id) {
                setCategoryId(created.id)
                setSubcategoryId('')
            }
            setNewCategoryName('')
        } catch (error) {
            setInlineCreateError(extractApiErrorMessage(error, 'Unable to add category right now.'))
        }
    }

    async function handleCreateSubcategory() {
        const trimmed = newSubcategoryName.trim()
        if (!trimmed || !categoryId) return

        try {
            setInlineCreateError(null)
            const created = (await createSubcategory.mutateAsync({ name: trimmed })) as { id?: string }
            if (created?.id) {
                setSubcategoryId(created.id)
            }
            setNewSubcategoryName('')
        } catch (error) {
            setInlineCreateError(extractApiErrorMessage(error, 'Unable to add subcategory right now.'))
        }
    }

    async function handleSubmit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault()

        if (!accountId || !categoryId || !name.trim() || !amount) {
            return
        }

        const parsedAmount = Number(amount)
        if (Number.isNaN(parsedAmount) || parsedAmount <= 0) {
            return
        }

        const signedAmount = kind === 'expense' ? -Math.abs(parsedAmount) : Math.abs(parsedAmount)

        try {
            setTransactionError(null)
            await createTransaction.mutateAsync({
                account_id: accountId,
                category_id: categoryId,
                subcategory_id: subcategoryId || undefined,
                name: name.trim(),
                amount: signedAmount,
                transaction_date: normalizeTransactionDate(transactionDate),
                notes: notes.trim() || undefined,
                is_scheduled: isScheduled,
            })
        } catch (error) {
            setTransactionError(extractApiErrorMessage(error, 'Unable to create transaction.'))
        }
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
            <div className="w-full max-w-2xl rounded-2xl border border-border bg-surface shadow-xl">
                <div className="flex items-center justify-between border-b border-border px-6 py-4 border-border">
                    <h2 className="text-xl font-semibold text-foreground">Add Transaction</h2>
                    <button
                        type="button"
                        onClick={handleClose}
                        className="rounded-md px-3 py-1.5 text-sm text-muted hover:bg-surface-hover"
                    >
                        Close
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="space-y-5 px-6 py-5">
                    <div className="flex gap-2">
                        <button
                            type="button"
                            onClick={() => setKind('expense')}
                            className={`rounded-full px-4 py-2 text-sm font-medium ${kind === 'expense' ? 'bg-rose-100 text-negative' : 'bg-elevated text-secondary'
                                }`}
                        >
                            Expense
                        </button>
                        <button
                            type="button"
                            onClick={() => setKind('income')}
                            className={`rounded-full px-4 py-2 text-sm font-medium ${kind === 'income' ? 'bg-positive-muted text-positive' : 'bg-elevated text-secondary'
                                }`}
                        >
                            Income
                        </button>
                    </div>

                    <div className="grid gap-4 md:grid-cols-2">
                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Account</span>
                            <select
                                value={accountId}
                                onChange={(e) => setAccountId(e.target.value)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                                required
                            >
                                <option value="">Select account</option>
                                {accounts.map((item) => (
                                    <option key={item.id} value={item.id}>
                                        {item.name}
                                    </option>
                                ))}
                            </select>
                            <div className="flex gap-2">
                                <input
                                    value={newAccountName}
                                    onChange={(e) => setNewAccountName(e.target.value)}
                                    placeholder="New account"
                                    className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                                />
                                <input
                                    value={newAccountOpeningBalance}
                                    onChange={(e) => setNewAccountOpeningBalance(e.target.value)}
                                    placeholder="Opening"
                                    inputMode="decimal"
                                    className="w-32 rounded-lg border border-border px-3 py-2 text-sm"
                                />
                                <button
                                    type="button"
                                    onClick={handleCreateAccount}
                                    disabled={createAccount.isPending}
                                    className="rounded-lg border border-border px-3 text-sm"
                                >
                                    Add
                                </button>
                            </div>
                        </label>

                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Category</span>
                            <select
                                value={categoryId}
                                onChange={(e) => {
                                    setCategoryId(e.target.value)
                                    setSubcategoryId('')
                                }}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                                required
                            >
                                <option value="">Select category</option>
                                {categories.map((item) => (
                                    <option key={item.id} value={item.id}>
                                        {item.name}
                                    </option>
                                ))}
                            </select>
                            <div className="flex gap-2">
                                <input
                                    value={newCategoryName}
                                    onChange={(e) => setNewCategoryName(e.target.value)}
                                    placeholder="New category"
                                    className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                                />
                                <button
                                    type="button"
                                    onClick={handleCreateCategory}
                                    disabled={createCategory.isPending}
                                    className="rounded-lg border border-border px-3 text-sm"
                                >
                                    Add
                                </button>
                            </div>
                        </label>
                    </div>

                    <div className="grid gap-4 md:grid-cols-2">
                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Subcategory</span>
                            <select
                                value={subcategoryId}
                                onChange={(e) => setSubcategoryId(e.target.value)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                            >
                                <option value="">Auto / none</option>
                                {(selectedCategory?.subcategories ?? []).map((item) => (
                                    <option key={item.id} value={item.id}>
                                        {item.name}
                                    </option>
                                ))}
                            </select>
                            <div className="flex gap-2">
                                <input
                                    value={newSubcategoryName}
                                    onChange={(e) => setNewSubcategoryName(e.target.value)}
                                    placeholder="New subcategory"
                                    disabled={!categoryId}
                                    className="w-full rounded-lg border border-border px-3 py-2 text-sm disabled:bg-elevated"
                                />
                                <button
                                    type="button"
                                    onClick={handleCreateSubcategory}
                                    disabled={!categoryId || createSubcategory.isPending}
                                    className="rounded-lg border border-border px-3 text-sm disabled:opacity-60"
                                >
                                    Add
                                </button>
                            </div>
                        </label>

                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Name</span>
                            <input
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                                placeholder="e.g. Groceries"
                                required
                            />
                        </label>
                    </div>

                    <div className="grid gap-4 md:grid-cols-3">
                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Amount</span>
                            <input
                                value={amount}
                                onChange={(e) => setAmount(e.target.value)}
                                inputMode="decimal"
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                                placeholder="0.00"
                                required
                            />
                        </label>

                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Date</span>
                            <input
                                type="date"
                                value={transactionDate}
                                onChange={(e) => setTransactionDate(e.target.value)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                                required
                            />
                        </label>

                        <label className="flex items-end gap-2">
                            <input
                                type="checkbox"
                                checked={isScheduled}
                                onChange={(e) => setIsScheduled(e.target.checked)}
                            />
                            <span className="text-sm text-secondary">Scheduled</span>
                        </label>
                    </div>

                    <div className="space-y-2">
                        <p className="text-xs font-semibold uppercase tracking-wide text-muted">Quick Fill Amount</p>
                        <div className="flex flex-wrap gap-2">
                            {QUICK_AMOUNTS.map((value) => (
                                <button
                                    key={value}
                                    type="button"
                                    onClick={() => setAmount(String(value))}
                                    className="rounded-full border border-border px-3 py-1 text-sm text-secondary"
                                >
                                    {value}
                                </button>
                            ))}
                        </div>
                    </div>

                    <label className="space-y-1">
                        <span className="text-sm font-medium text-secondary">Notes</span>
                        <textarea
                            value={notes}
                            onChange={(e) => setNotes(e.target.value)}
                            rows={3}
                            className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                        />
                    </label>

                    {transactionError ? (
                        <p className="rounded-md bg-negative-muted px-3 py-2 text-sm text-negative">{transactionError}</p>
                    ) : null}

                    {inlineCreateError ? (
                        <p className="rounded-md bg-negative-muted px-3 py-2 text-sm text-negative">{inlineCreateError}</p>
                    ) : null}

                    <div className="flex justify-end gap-2 border-t border-border pt-4">
                        <button
                            type="button"
                            onClick={handleClose}
                            className="rounded-lg border border-border px-4 py-2 text-sm text-secondary"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={createTransaction.isPending}
                            className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white disabled:opacity-60"
                        >
                            {createTransaction.isPending ? 'Saving...' : 'Save Transaction'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    )
}
