import { useState } from 'react'
import type { FormEvent } from 'react'
import { DndContext, PointerSensor, closestCenter, useSensor, useSensors } from '@dnd-kit/core'
import type { DragEndEvent } from '@dnd-kit/core'
import { SortableContext, arrayMove, useSortable, verticalListSortingStrategy } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'

import { PageShell } from '@/components/layout/PageShell'
import { useAccounts } from '@/hooks/useAccounts'
import { useCategories } from '@/hooks/useCategories'
import {
    useCreateQuickTransaction,
    useExecuteQuickTransaction,
    useQuickTransactions,
    useReorderQuickTransactions,
} from '@/hooks/useQuickTransactions'
import type { QuickTransaction } from '@/types'

export function QuickTransactionsPage() {
    const { data: quickTransactions = [] } = useQuickTransactions()
    const { data: accounts = [] } = useAccounts()
    const { data: categories = [] } = useCategories()

    const [ordered, setOrdered] = useState<QuickTransaction[]>([])
    const [label, setLabel] = useState('')
    const [amount, setAmount] = useState('')
    const [accountId, setAccountId] = useState('')
    const [categoryId, setCategoryId] = useState('')

    const createQuick = useCreateQuickTransaction()
    const executeQuick = useExecuteQuickTransaction()
    const reorderQuick = useReorderQuickTransactions()

    const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 6 } }))

    const list = ordered.length > 0 ? ordered : [...quickTransactions].sort((a, b) => a.sort_order - b.sort_order)

    async function handleCreate(e: FormEvent<HTMLFormElement>) {
        e.preventDefault()
        const trimmed = label.trim()
        if (!trimmed) {
            return
        }

        await createQuick.mutateAsync({
            label: trimmed,
            amount: amount ? Number(amount) : undefined,
            account_id: accountId || undefined,
            category_id: categoryId || undefined,
        })

        setLabel('')
        setAmount('')
        setAccountId('')
        setCategoryId('')
    }

    async function handleDragEnd(event: DragEndEvent) {
        const { active, over } = event
        if (!over || active.id === over.id) {
            return
        }

        const oldIndex = list.findIndex((item) => item.id === active.id)
        const newIndex = list.findIndex((item) => item.id === over.id)
        if (oldIndex < 0 || newIndex < 0) {
            return
        }

        const moved = arrayMove(list, oldIndex, newIndex)
        setOrdered(moved)
        await reorderQuick.mutateAsync(moved.map((item) => item.id))
    }

    return (
        <PageShell title="Quick Transactions">
            <div className="grid gap-4 lg:grid-cols-[380px,1fr]">
                <section className="rounded-2xl border border-border bg-surface p-4 shadow-sm">
                    <h2 className="text-base font-semibold text-foreground">New Template</h2>
                    <form className="mt-3 space-y-3" onSubmit={handleCreate}>
                        <input
                            value={label}
                            onChange={(e) => setLabel(e.target.value)}
                            placeholder="Label"
                            className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                            required
                        />
                        <input
                            value={amount}
                            onChange={(e) => setAmount(e.target.value)}
                            inputMode="decimal"
                            placeholder="Amount"
                            className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                        />
                        <select
                            value={accountId}
                            onChange={(e) => setAccountId(e.target.value)}
                            className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                        >
                            <option value="">Any account</option>
                            {accounts.map((account) => (
                                <option key={account.id} value={account.id}>
                                    {account.name}
                                </option>
                            ))}
                        </select>
                        <select
                            value={categoryId}
                            onChange={(e) => setCategoryId(e.target.value)}
                            className="w-full rounded-lg border border-border px-3 py-2 text-sm"
                        >
                            <option value="">Any category</option>
                            {categories.map((category) => (
                                <option key={category.id} value={category.id}>
                                    {category.name}
                                </option>
                            ))}
                        </select>
                        <button
                            type="submit"
                            disabled={createQuick.isPending}
                            className="w-full rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white disabled:opacity-60"
                        >
                            {createQuick.isPending ? 'Saving...' : 'Create Template'}
                        </button>
                    </form>
                </section>

                <section className="rounded-2xl border border-border bg-surface p-4 shadow-sm">
                    <div className="mb-3 flex items-center justify-between">
                        <h2 className="text-base font-semibold text-foreground">Templates</h2>
                        <p className="text-xs text-muted">Drag to reorder</p>
                    </div>

                    <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
                        <SortableContext items={list.map((item) => item.id)} strategy={verticalListSortingStrategy}>
                            <div className="space-y-2">
                                {list.map((item) => (
                                    <QuickTransactionCard
                                        key={item.id}
                                        item={item}
                                        accountName={accounts.find((a) => a.id === item.account_id)?.name}
                                        categoryName={categories.find((c) => c.id === item.category_id)?.name}
                                        onExecute={() => executeQuick.mutate(item.id)}
                                        executing={executeQuick.isPending}
                                    />
                                ))}
                            </div>
                        </SortableContext>
                    </DndContext>
                </section>
            </div>
        </PageShell>
    )
}

function QuickTransactionCard({
    item,
    accountName,
    categoryName,
    onExecute,
    executing,
}: {
    item: QuickTransaction
    accountName?: string
    categoryName?: string
    onExecute: () => void
    executing: boolean
}) {
    const { attributes, listeners, setNodeRef, transform, transition } = useSortable({ id: item.id })

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
    }

    return (
        <article
            ref={setNodeRef}
            style={style}
            className="rounded-xl border border-border bg-elevated p-3"
        >
            <div className="flex items-start justify-between gap-3">
                <button
                    type="button"
                    className="cursor-grab rounded-md border border-border px-2 py-1 text-xs text-secondary"
                    {...attributes}
                    {...listeners}
                >
                    Drag
                </button>
                <div className="flex-1">
                    <p className="text-sm font-semibold text-foreground">{item.label}</p>
                    <p className="mt-1 text-xs text-muted">
                        {accountName ?? 'Any account'} • {categoryName ?? 'Any category'}
                    </p>
                </div>
                <button
                    type="button"
                    onClick={onExecute}
                    disabled={executing}
                    className="rounded-md bg-accent px-3 py-1.5 text-xs font-medium text-white disabled:opacity-60"
                >
                    Run
                </button>
            </div>
        </article>
    )
}
