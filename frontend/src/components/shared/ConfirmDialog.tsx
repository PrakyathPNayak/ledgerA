interface ConfirmDialogProps {
    isOpen: boolean
    title: string
    description?: string
    confirmLabel?: string
    cancelLabel?: string
    intent?: 'default' | 'danger'
    isLoading?: boolean
    onConfirm: () => void
    onCancel: () => void
}

export function ConfirmDialog({
    isOpen,
    title,
    description,
    confirmLabel = 'Confirm',
    cancelLabel = 'Cancel',
    intent = 'default',
    isLoading = false,
    onConfirm,
    onCancel,
}: ConfirmDialogProps) {
    if (!isOpen) {
        return null
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 p-4">
            <div className="w-full max-w-md rounded-2xl border border-slate-200 bg-white p-6 shadow-xl">
                <h2 className="text-lg font-semibold text-slate-900">{title}</h2>
                {description ? <p className="mt-2 text-sm text-slate-600">{description}</p> : null}

                <div className="mt-6 flex justify-end gap-2">
                    <button
                        type="button"
                        onClick={onCancel}
                        className="rounded-lg border border-slate-300 px-4 py-2 text-sm text-slate-700"
                    >
                        {cancelLabel}
                    </button>
                    <button
                        type="button"
                        onClick={onConfirm}
                        disabled={isLoading}
                        className={`rounded-lg px-4 py-2 text-sm font-medium text-white disabled:opacity-60 ${intent === 'danger' ? 'bg-rose-600' : 'bg-slate-900'
                            }`}
                    >
                        {isLoading ? 'Working...' : confirmLabel}
                    </button>
                </div>
            </div>
        </div>
    )
}
