import type { Transaction } from '@/types'

interface TransactionRowProps {
    transaction: Transaction
    accountName?: string
    categoryName?: string
    subcategoryName?: string
    onClick?: (transaction: Transaction) => void
}

function currency(amount: number): string {
    return new Intl.NumberFormat('en-IN', {
        style: 'currency',
        currency: 'INR',
        maximumFractionDigits: 2,
    }).format(amount)
}

export function TransactionRow({ transaction, accountName, categoryName, subcategoryName, onClick }: TransactionRowProps) {
    const isIncome = transaction.amount > 0

    return (
        <tr
            onClick={() => onClick?.(transaction)}
            className="cursor-pointer transition-colors hover:bg-elevated"
        >
            <td className="px-4 py-3 text-sm text-foreground">{transaction.transaction_date}</td>
            <td className="px-4 py-3 text-sm font-medium text-foreground">{transaction.name}</td>
            <td className="px-4 py-3 text-sm text-secondary">{accountName ?? 'Unknown'}</td>
            <td className="px-4 py-3 text-sm text-secondary">
                {categoryName ?? 'Unknown'}
                {subcategoryName ? <span className="text-muted"> / {subcategoryName}</span> : null}
            </td>
            <td className={`px-4 py-3 text-right text-sm font-semibold ${isIncome ? 'text-positive' : 'text-negative'}`}>
                {currency(transaction.amount)}
            </td>
        </tr>
    )
}
