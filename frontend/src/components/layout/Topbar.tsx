import { Settings } from 'lucide-react'

/**
 * @description Topbar properties.
 */
interface TopbarProps {
  title: string
  onOpenSettings: () => void
  userInitial?: string
}

export function Topbar({ title, onOpenSettings, userInitial }: TopbarProps) {
  return (
    <header className="flex items-center justify-between border-b border-border bg-surface px-6 py-3">
      <div className="flex items-center gap-3">
        <div className="grid h-9 w-9 place-items-center rounded-full bg-accent text-sm font-bold text-white">
          {userInitial ?? 'U'}
        </div>
        <h1 className="text-lg font-semibold text-foreground">{title}</h1>
      </div>
      <button
        className="rounded-lg border border-border p-2 text-secondary hover:bg-surface-hover border-border text-secondary hover:bg-surface-hover"
        onClick={onOpenSettings}
      >
        <Settings size={16} />
      </button>
    </header>
  )
}
