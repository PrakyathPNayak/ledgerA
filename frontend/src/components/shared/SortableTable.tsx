import type { ReactNode } from 'react'

export interface SortColumn {
    key: string
    label: string
    sortable?: boolean
    className?: string
}

interface SortableTableProps {
    columns: SortColumn[]
    sortBy?: string
    sortDir?: 'asc' | 'desc'
    onSort?: (columnKey: string) => void
    children: ReactNode
    emptyMessage?: string
    hasRows: boolean
}

export function SortableTable({
    columns,
    sortBy,
    sortDir = 'desc',
    onSort,
    children,
    emptyMessage = 'No rows found.',
    hasRows,
}: SortableTableProps) {
    return (
        <div className="overflow-hidden rounded-2xl border border-border bg-surface shadow-sm">
            <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-border">
                    <thead className="bg-elevated">
                        <tr>
                            {columns.map((column) => {
                                const active = sortBy === column.key
                                return (
                                    <th key={column.key} className={`px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-muted ${column.className ?? ''}`}>
                                        {column.sortable && onSort ? (
                                            <button
                                                type="button"
                                                onClick={() => onSort(column.key)}
                                                className="inline-flex items-center gap-1 hover:text-foreground"
                                            >
                                                <span>{column.label}</span>
                                                {active ? <span>{sortDir === 'asc' ? '↑' : '↓'}</span> : null}
                                            </button>
                                        ) : (
                                            column.label
                                        )}
                                    </th>
                                )
                            })}
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-border-subtle bg-surface">{hasRows ? children : null}</tbody>
                </table>
            </div>
            {!hasRows ? <p className="px-4 py-8 text-center text-sm text-muted">{emptyMessage}</p> : null}
        </div>
    )
}
