import type { Transaction } from '@/types'

interface TransactionDetailModalProps {
    transaction: Transaction | null
    accountName?: string
    categoryName?: string
    subcategoryName?: string
    onClose: () => void
    onEdit?: (transaction: Transaction) => void
    onDelete?: (transaction: Transaction) => void
}

function currency(amount: number): string {
    return new Intl.NumberFormat('en-IN', {
        style: 'currency',
        currency: 'INR',
        maximumFractionDigits: 2,
    }).format(amount)
}

export function TransactionDetailModal({
    transaction,
    accountName,
    categoryName,
    subcategoryName,
    onClose,
    onEdit,
    onDelete,
}: TransactionDetailModalProps) {
    if (!transaction) {
        return null
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 p-4">
            <div className="w-full max-w-lg rounded-2xl border border-slate-200 bg-white p-6 shadow-xl">
                <div className="flex items-center justify-between">
                    <h2 className="text-xl font-semibold text-slate-900">{transaction.name}</h2>
                    <button
                        type="button"
                        onClick={onClose}
                        className="rounded-md px-3 py-1.5 text-sm text-slate-500 hover:bg-slate-100"
                    >
                        Close
                    </button>
                </div>

                <dl className="mt-5 space-y-3 text-sm">
                    <div className="flex justify-between gap-6">
                        <dt className="text-slate-500">Amount</dt>
                        <dd className="font-semibold text-slate-900">{currency(transaction.amount)}</dd>
                    </div>
                    <div className="flex justify-between gap-6">
                        <dt className="text-slate-500">Date</dt>
                        <dd className="text-slate-800">{transaction.transaction_date}</dd>
                    </div>
                    <div className="flex justify-between gap-6">
                        <dt className="text-slate-500">Account</dt>
                        <dd className="text-slate-800">{accountName ?? transaction.account_id}</dd>
                    </div>
                    <div className="flex justify-between gap-6">
                        <dt className="text-slate-500">Category</dt>
                        <dd className="text-slate-800">
                            {categoryName ?? transaction.category_id}
                            {subcategoryName ? ` / ${subcategoryName}` : ''}
                        </dd>
                    </div>
                    {transaction.notes ? (
                        <div className="space-y-1">
                            <dt className="text-slate-500">Notes</dt>
                            <dd className="rounded-lg bg-slate-50 p-3 text-slate-700">{transaction.notes}</dd>
                        </div>
                    ) : null}
                </dl>

                <div className="mt-6 flex justify-end gap-2">
                    <button
                        type="button"
                        onClick={() => onEdit?.(transaction)}
                        className="rounded-lg border border-slate-300 px-4 py-2 text-sm text-slate-700"
                    >
                        Edit
                    </button>
                    <button
                        type="button"
                        onClick={() => onDelete?.(transaction)}
                        className="rounded-lg bg-rose-600 px-4 py-2 text-sm font-medium text-white"
                    >
                        Delete
                    </button>
                </div>
            </div>
        </div>
    )
}
