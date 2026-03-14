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
            className="cursor-pointer transition-colors hover:bg-slate-50"
        >
            <td className="px-4 py-3 text-sm text-slate-900">{transaction.transaction_date}</td>
            <td className="px-4 py-3 text-sm font-medium text-slate-900">{transaction.name}</td>
            <td className="px-4 py-3 text-sm text-slate-600">{accountName ?? 'Unknown'}</td>
            <td className="px-4 py-3 text-sm text-slate-600">
                {categoryName ?? 'Unknown'}
                {subcategoryName ? <span className="text-slate-400"> / {subcategoryName}</span> : null}
            </td>
            <td className={`px-4 py-3 text-right text-sm font-semibold ${isIncome ? 'text-emerald-600' : 'text-rose-600'}`}>
                {currency(transaction.amount)}
            </td>
        </tr>
    )
}
