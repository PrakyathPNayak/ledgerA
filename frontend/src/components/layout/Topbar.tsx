import { Settings } from 'lucide-react'

/**
 * @description Topbar properties.
 */
interface TopbarProps {
  title: string
  onOpenSettings: () => void
}

export function Topbar({ title, onOpenSettings }: TopbarProps) {
  return (
    <header className="flex items-center justify-between border-b border-slate-200 bg-white px-6 py-3 dark:border-slate-800 dark:bg-slate-900">
      <div className="flex items-center gap-3">
        <div className="grid h-9 w-9 place-items-center rounded-full bg-slate-900 text-sm font-bold text-white dark:bg-slate-100 dark:text-slate-900">
          9
        </div>
        <h1 className="text-lg font-semibold text-slate-900 dark:text-slate-100">{title}</h1>
      </div>
      <button
        className="rounded-lg border border-slate-300 p-2 text-slate-600 hover:bg-slate-100 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
        onClick={onOpenSettings}
      >
        <Settings size={16} />
      </button>
    </header>
  )
}
