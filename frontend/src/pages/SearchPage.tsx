import { useState } from 'react'
import { useInfiniteQuery } from '@tanstack/react-query'

import { PageShell } from '@/components/layout/PageShell'
import { SortableTable } from '@/components/shared/SortableTable'
import { TransactionRow } from '@/components/shared/TransactionRow'
import { useAccounts } from '@/hooks/useAccounts'
import { useCategories } from '@/hooks/useCategories'
import api from '@/lib/api'
import type { Transaction } from '@/types'

const columns = [
    { key: 'transaction_date', label: 'Date' },
    { key: 'name', label: 'Name' },
    { key: 'account', label: 'Account' },
    { key: 'category', label: 'Category' },
    { key: 'amount', label: 'Amount', className: 'text-right' },
]

export function SearchPage() {
    const [searchText, setSearchText] = useState('')
    const [accountId, setAccountId] = useState('')
    const [categoryId, setCategoryId] = useState('')
    const [type, setType] = useState<'income' | 'expense' | 'all'>('all')

    const { data: accounts = [] } = useAccounts()
    const { data: categories = [] } = useCategories()

    const filters = {
        search: searchText || undefined,
        account_id: accountId || undefined,
        category_id: categoryId || undefined,
        type,
        sort_by: 'transaction_date',
        sort_dir: 'desc',
    } as const

    const query = useInfiniteQuery({
        queryKey: ['transactions-search', filters],
        initialPageParam: 1,
        queryFn: ({ pageParam }) =>
            api.get('/transactions', {
                params: {
                    ...filters,
                    page: pageParam,
                    per_page: 25,
                },
            }) as Promise<Transaction[]>,
        getNextPageParam: (lastPage, pages) => (lastPage.length >= 25 ? pages.length + 1 : undefined),
    })

    const items = query.data?.pages.flat() ?? []

    const accountMap = Object.fromEntries(accounts.map((item) => [item.id, item.name]))
    const categoryMap = Object.fromEntries(categories.map((item) => [item.id, item.name]))
    const subcategoryMap = Object.fromEntries(
        categories.flatMap((category) => (category.subcategories ?? []).map((sub) => [sub.id, sub.name])),
    )

    return (
        <PageShell title="Search">
            <div className="space-y-4">
                <section className="rounded-2xl border border-border bg-surface p-4 shadow-sm border-border bg-surface">
                    <div className="grid gap-3 md:grid-cols-4">
                        <input
                            value={searchText}
                            onChange={(e) => setSearchText(e.target.value)}
                            placeholder="Search name / notes"
                            className="rounded-lg border border-border px-3 py-2 text-sm"
                        />

                        <select
                            value={accountId}
                            onChange={(e) => setAccountId(e.target.value)}
                            className="rounded-lg border border-border px-3 py-2 text-sm"
                        >
                            <option value="">All accounts</option>
                            {accounts.map((account) => (
                                <option key={account.id} value={account.id}>
                                    {account.name}
                                </option>
                            ))}
                        </select>

                        <select
                            value={categoryId}
                            onChange={(e) => setCategoryId(e.target.value)}
                            className="rounded-lg border border-border px-3 py-2 text-sm"
                        >
                            <option value="">All categories</option>
                            {categories.map((category) => (
                                <option key={category.id} value={category.id}>
                                    {category.name}
                                </option>
                            ))}
                        </select>

                        <select
                            value={type}
                            onChange={(e) => setType(e.target.value as 'income' | 'expense' | 'all')}
                            className="rounded-lg border border-border px-3 py-2 text-sm"
                        >
                            <option value="all">All</option>
                            <option value="income">Income</option>
                            <option value="expense">Expense</option>
                        </select>
                    </div>
                </section>

                <section className="space-y-3 rounded-2xl border border-border bg-surface p-4 shadow-sm border-border bg-surface">
                    <SortableTable
                        columns={columns}
                        hasRows={items.length > 0}
                        emptyMessage="No matching transactions found."
                    >
                        {items.map((transaction) => (
                            <TransactionRow
                                key={transaction.id}
                                transaction={transaction}
                                accountName={accountMap[transaction.account_id]}
                                categoryName={categoryMap[transaction.category_id]}
                                subcategoryName={subcategoryMap[transaction.subcategory_id]}
                            />
                        ))}
                    </SortableTable>

                    <div className="flex justify-center">
                        <button
                            type="button"
                            onClick={() => query.fetchNextPage()}
                            disabled={!query.hasNextPage || query.isFetchingNextPage}
                            className="rounded-lg border border-border px-4 py-2 text-sm text-secondary disabled:opacity-60"
                        >
                            {query.isFetchingNextPage ? 'Loading...' : query.hasNextPage ? 'Load More' : 'No More Results'}
                        </button>
                    </div>
                </section>
            </div>
        </PageShell>
    )
}
