import { Link } from 'react-router-dom'
import { useState } from 'react'
import type { FormEvent } from 'react'

import { PageShell } from '@/components/layout/PageShell'
import { useAccounts, useCreateAccount } from '@/hooks/useAccounts'

function money(value: number): string {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR' }).format(value)
}

export function AccountsPage() {
    const { data: accounts = [], isLoading } = useAccounts()
    const createAccount = useCreateAccount()

    const [showForm, setShowForm] = useState(false)
    const [name, setName] = useState('')
    const [openingBalance, setOpeningBalance] = useState('0')
    const [error, setError] = useState<string | null>(null)

    async function handleCreateAccount(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()

        const trimmedName = name.trim()
        if (!trimmedName) {
            setError('Account name is required.')
            return
        }

        const parsedOpeningBalance = Number(openingBalance)
        if (!Number.isFinite(parsedOpeningBalance)) {
            setError('Opening balance must be a valid number.')
            return
        }

        setError(null)
        try {
            await createAccount.mutateAsync({
                name: trimmedName,
                opening_balance: parsedOpeningBalance,
            })

            setShowForm(false)
            setName('')
            setOpeningBalance('0')
        } catch {
            setError('Unable to create account. Please try again.')
        }
    }

    return (
        <PageShell title="Accounts">
            <section className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm dark:border-slate-700 dark:bg-slate-900">
                <div className="mb-4 flex items-center justify-between">
                    <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">All Accounts</h2>
                    <div className="flex items-center gap-3">
                        <p className="text-sm text-slate-500">{accounts.length} accounts</p>
                        <button
                            type="button"
                            onClick={() => {
                                setShowForm((prev) => !prev)
                                setError(null)
                            }}
                            className="rounded-lg bg-slate-900 px-3 py-2 text-sm font-medium text-white"
                        >
                            {showForm ? 'Cancel' : 'Add Account'}
                        </button>
                    </div>
                </div>

                {showForm ? (
                    <form onSubmit={handleCreateAccount} className="mb-4 grid gap-3 rounded-xl border border-slate-200 bg-slate-50 p-3 md:grid-cols-3">
                        <label className="space-y-1 md:col-span-2">
                            <span className="text-sm font-medium text-slate-700">Account name</span>
                            <input
                                value={name}
                                onChange={(event) => setName(event.target.value)}
                                className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm"
                                placeholder="e.g. HDFC Savings"
                                required
                            />
                        </label>
                        <label className="space-y-1">
                            <span className="text-sm font-medium text-slate-700">Opening balance</span>
                            <input
                                value={openingBalance}
                                onChange={(event) => setOpeningBalance(event.target.value)}
                                inputMode="decimal"
                                className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm"
                                placeholder="0"
                            />
                        </label>

                        {error ? <p className="text-sm text-rose-700 md:col-span-3">{error}</p> : null}

                        <div className="md:col-span-3">
                            <button
                                type="submit"
                                disabled={createAccount.isPending}
                                className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-medium text-white disabled:opacity-60"
                            >
                                {createAccount.isPending ? 'Creating...' : 'Create Account'}
                            </button>
                        </div>
                    </form>
                ) : null}

                {isLoading ? <p className="text-sm text-slate-500">Loading accounts...</p> : null}

                <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
                    {accounts.map((account) => (
                        <Link
                            key={account.id}
                            to={`/accounts/${account.id}`}
                            className="rounded-xl border border-slate-200 p-4 transition-colors hover:border-slate-300 hover:bg-slate-50 dark:border-slate-700 dark:hover:border-slate-600 dark:hover:bg-slate-800"
                        >
                            <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">{account.name}</p>
                            <p className="mt-1 text-xs uppercase tracking-wide text-slate-500">{account.account_type || 'cash'}</p>
                            <p className="mt-4 text-lg font-bold text-slate-900 dark:text-slate-100">{money(account.current_balance)}</p>
                            <p className="mt-1 text-xs text-slate-500">Opening: {money(account.opening_balance)}</p>
                        </Link>
                    ))}
                </div>
            </section>
        </PageShell>
    )
}
