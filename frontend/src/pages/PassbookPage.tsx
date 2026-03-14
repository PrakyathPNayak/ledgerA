import { Link, useParams } from 'react-router-dom'

import { PageShell } from '@/components/layout/PageShell'
import { SortableTable } from '@/components/shared/SortableTable'
import { TransactionRow } from '@/components/shared/TransactionRow'
import { useAccounts } from '@/hooks/useAccounts'
import { useCategories } from '@/hooks/useCategories'
import { useTransactions } from '@/hooks/useTransactions'
import type { Transaction } from '@/types'

const columns = [
    { key: 'transaction_date', label: 'Date', sortable: true },
    { key: 'name', label: 'Name', sortable: true },
    { key: 'account', label: 'Account' },
    { key: 'category', label: 'Category' },
    { key: 'amount', label: 'Amount', sortable: true, className: 'text-right' },
]

function money(value: number): string {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR' }).format(value)
}

export function PassbookPage() {
    const { id = '' } = useParams()
    const { data: accounts = [] } = useAccounts()
    const { data: categories = [] } = useCategories()
    const { data: transactions = [] } = useTransactions({ account_id: id, passbook_mode: true, per_page: 50, page: 1 })

    const account = accounts.find((item) => item.id === id)

    const accountMap = Object.fromEntries(accounts.map((item) => [item.id, item.name]))
    const categoryMap = Object.fromEntries(categories.map((item) => [item.id, item.name]))
    const subcategoryMap = Object.fromEntries(
        categories.flatMap((category) => (category.subcategories ?? []).map((sub) => [sub.id, sub.name])),
    )

    return (
        <PageShell title="Passbook">
            <section className="space-y-4 rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
                <div className="flex items-center justify-between">
                    <div>
                        <p className="text-sm text-slate-500">Account</p>
                        <h2 className="text-2xl font-bold text-slate-900">{account?.name ?? 'Unknown account'}</h2>
                    </div>
                    <div className="text-right">
                        <p className="text-xs uppercase tracking-wide text-slate-500">Current Balance</p>
                        <p className="text-2xl font-bold text-slate-900">{money(account?.current_balance ?? 0)}</p>
                    </div>
                </div>

                <SortableTable
                    columns={columns}
                    hasRows={transactions.length > 0}
                    emptyMessage="No passbook entries yet for this account."
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

                <Link to="/accounts" className="inline-flex text-sm font-medium text-slate-700 underline-offset-2 hover:underline">
                    Back to accounts
                </Link>
            </section>
        </PageShell>
    )
}
