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
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
            <div className="w-full max-w-lg rounded-2xl border border-border bg-surface p-6 shadow-xl">
                <div className="flex items-center justify-between">
                    <h2 className="text-xl font-semibold text-foreground">{transaction.name}</h2>
                    <button
                        type="button"
                        onClick={onClose}
                        className="rounded-md px-3 py-1.5 text-sm text-muted hover:bg-surface-hover"
                    >
                        Close
                    </button>
                </div>

                <dl className="mt-5 space-y-3 text-sm">
                    <div className="flex justify-between gap-6">
                        <dt className="text-muted">Amount</dt>
                        <dd className="font-semibold text-foreground">{currency(transaction.amount)}</dd>
                    </div>
                    <div className="flex justify-between gap-6">
                        <dt className="text-muted">Date</dt>
                        <dd className="text-foreground">{transaction.transaction_date}</dd>
                    </div>
                    <div className="flex justify-between gap-6">
                        <dt className="text-muted">Account</dt>
                        <dd className="text-foreground">{accountName ?? transaction.account_id}</dd>
                    </div>
                    <div className="flex justify-between gap-6">
                        <dt className="text-muted">Category</dt>
                        <dd className="text-foreground">
                            {categoryName ?? transaction.category_id}
                            {subcategoryName ? ` / ${subcategoryName}` : ''}
                        </dd>
                    </div>
                    {transaction.notes ? (
                        <div className="space-y-1">
                            <dt className="text-muted">Notes</dt>
                            <dd className="rounded-lg bg-elevated p-3 text-secondary">{transaction.notes}</dd>
                        </div>
                    ) : null}
                </dl>

                <div className="mt-6 flex justify-end gap-2">
                    <button
                        type="button"
                        onClick={() => onEdit?.(transaction)}
                        className="rounded-lg border border-border px-4 py-2 text-sm text-secondary"
                    >
                        Edit
                    </button>
                    <button
                        type="button"
                        onClick={() => onDelete?.(transaction)}
                        className="rounded-lg bg-negative px-4 py-2 text-sm font-medium text-white"
                    >
                        Delete
                    </button>
                </div>
            </div>
        </div>
    )
}
