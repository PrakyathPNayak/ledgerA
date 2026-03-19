import { useEffect, useState } from 'react'

import api from '@/lib/api'
import { useUIStore } from '@/store/uiStore'

/**
 * @description Settings sheet properties.
 */
interface SettingsSheetProps {
  open: boolean
  onClose: () => void
  displayName?: string
  currencyCode?: string
}

export function SettingsSheet({
  open,
  onClose,
  displayName,
  currencyCode,
}: SettingsSheetProps) {
  const { theme, setTheme } = useUIStore()
  const [name, setName] = useState(displayName ?? '')

  useEffect(() => {
    if (!open) {
      return
    }

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose()
      }
    }

    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [open, onClose])

  if (!open) {
    return null
  }

  const saveProfile = async () => {
    await api.patch('/users/me', { display_name: name })
    onClose()
  }

  return (
    <div className="fixed inset-0 z-50 bg-black/20" onClick={onClose}>
      <div
        className="absolute right-0 top-0 h-full w-full max-w-md border-l border-border bg-surface p-6 shadow-2xl"
        onClick={(event) => event.stopPropagation()}
      >
        <h2 className="text-lg font-semibold text-foreground">Settings</h2>
        <div className="mt-4 space-y-4">
          <label className="block text-sm">
            <span className="mb-1 block text-secondary">Display name</span>
            <input
              className="w-full rounded-lg border border-border px-3 py-2 border-border bg-elevated text-foreground"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </label>

          <div className="text-sm text-secondary">Currency: {currencyCode ?? 'Not set'}</div>

          <div>
            <span className="mb-2 block text-sm text-secondary">Theme</span>
            <div className="flex gap-2">
              {(['light', 'dark', 'system'] as const).map((value) => (
                <button
                  key={value}
                  className={`rounded-lg border px-3 py-1 text-sm ${theme === value
                      ? 'border-accent bg-accent text-white'
                      : 'border-border text-secondary border-border text-secondary'
                    }`}
                  onClick={() => setTheme(value)}
                >
                  {value}
                </button>
              ))}
            </div>
          </div>

          <div className="text-xs text-muted">Version 0.1.0</div>
        </div>

        <div className="mt-6 flex gap-2">
          <button className="rounded-lg bg-accent px-4 py-2 text-sm text-white" onClick={saveProfile}>
            Save
          </button>
          <button className="rounded-lg border border-border px-4 py-2 text-sm border-border text-foreground" onClick={onClose}>
            Close
          </button>
        </div>
      </div>
    </div>
  )
}
