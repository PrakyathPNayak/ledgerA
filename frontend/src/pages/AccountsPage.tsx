import { Link } from 'react-router-dom'
import { useState } from 'react'
import type { FormEvent } from 'react'

import { PageShell } from '@/components/layout/PageShell'
import { TransferModal } from '@/components/shared/TransferModal'
import { useAccounts, useCreateAccount, useArchiveAccount } from '@/hooks/useAccounts'

function money(value: number): string {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR' }).format(value)
}

export function AccountsPage() {
    const { data: accounts = [], isLoading } = useAccounts()
    const createAccount = useCreateAccount()
    const archiveAccount = useArchiveAccount()

    const [showForm, setShowForm] = useState(false)
    const [showTransfer, setShowTransfer] = useState(false)
    const [showArchived, setShowArchived] = useState(false)
    const [name, setName] = useState('')
    const [openingBalance, setOpeningBalance] = useState('0')
    const [error, setError] = useState<string | null>(null)

    const activeAccounts = accounts.filter(a => !a.is_archived)
    const archivedAccounts = accounts.filter(a => a.is_archived)
    const displayedAccounts = showArchived ? archivedAccounts : activeAccounts

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
            <section className="rounded-2xl border border-border bg-surface p-4 shadow-sm">
                <div className="mb-4 flex items-center justify-between">
                    <h2 className="text-lg font-semibold text-foreground">All Accounts</h2>
                    <div className="flex items-center gap-3">
                        {archivedAccounts.length > 0 && (
                            <button type="button" onClick={() => setShowArchived(!showArchived)} className="text-xs text-muted hover:underline">
                                {showArchived ? 'Show Active' : `Archived (${archivedAccounts.length})`}
                            </button>
                        )}
                        <p className="text-sm text-muted">{activeAccounts.length} active</p>
                        <button
                            type="button"
                            onClick={() => setShowTransfer(true)}
                            className="rounded-lg border border-border px-3 py-2 text-sm font-medium text-secondary"
                        >
                            Transfer
                        </button>
                        <button
                            type="button"
                            onClick={() => {
                                setShowForm((prev) => !prev)
                                setError(null)
                            }}
                            className="rounded-lg bg-accent px-3 py-2 text-sm font-medium text-white"
                        >
                            {showForm ? 'Cancel' : 'Add Account'}
                        </button>
                    </div>
                </div>

                {showForm ? (
                    <form onSubmit={handleCreateAccount} className="mb-4 grid gap-3 rounded-xl border border-border bg-elevated p-3 md:grid-cols-3">
                        <label className="space-y-1 md:col-span-2">
                            <span className="text-sm font-medium text-secondary">Account name</span>
                            <input
                                value={name}
                                onChange={(event) => setName(event.target.value)}
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                                placeholder="e.g. HDFC Savings"
                                required
                            />
                        </label>
                        <label className="space-y-1">
                            <span className="text-sm font-medium text-secondary">Opening balance</span>
                            <input
                                value={openingBalance}
                                onChange={(event) => setOpeningBalance(event.target.value)}
                                inputMode="decimal"
                                className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                                placeholder="0"
                            />
                        </label>

                        {error ? <p className="text-sm text-negative md:col-span-3">{error}</p> : null}

                        <div className="md:col-span-3">
                            <button
                                type="submit"
                                disabled={createAccount.isPending}
                                className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white disabled:opacity-60"
                            >
                                {createAccount.isPending ? 'Creating...' : 'Create Account'}
                            </button>
                        </div>
                    </form>
                ) : null}

                {isLoading ? <p className="text-sm text-muted">Loading accounts...</p> : null}

                <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
                    {displayedAccounts.map((account) => (
                        <div key={account.id} className="relative rounded-xl border border-border p-4 transition-colors hover:border-accent/30 hover:bg-surface-hover">
                            <Link to={`/accounts/${account.id}`} className="block">
                                <p className="text-sm font-semibold text-foreground">{account.name}
                                    {account.is_archived && <span className="ml-2 text-xs text-muted">(archived)</span>}
                                </p>
                                <p className="mt-1 text-xs uppercase tracking-wide text-muted">{account.account_type || 'cash'}</p>
                                <p className="mt-4 text-lg font-bold text-foreground">{money(account.current_balance)}</p>
                                <p className="mt-1 text-xs text-muted">Opening: {money(account.opening_balance)}</p>
                            </Link>
                            <button
                                onClick={() => archiveAccount.mutate({ id: account.id, is_archived: !account.is_archived })}
                                className="absolute right-3 top-3 text-xs text-muted hover:text-secondary"
                                title={account.is_archived ? 'Unarchive' : 'Archive'}
                            >
                                {account.is_archived ? '↩' : '📦'}
                            </button>
                        </div>
                    ))}
                </div>
            </section>

            <TransferModal isOpen={showTransfer} onClose={() => setShowTransfer(false)} />
        </PageShell>
    )
}
