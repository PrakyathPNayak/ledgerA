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
        className="absolute right-0 top-0 h-full w-full max-w-md border-l border-slate-200 bg-white p-6 shadow-2xl dark:border-slate-700 dark:bg-slate-900"
        onClick={(event) => event.stopPropagation()}
      >
        <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">Settings</h2>
        <div className="mt-4 space-y-4">
          <label className="block text-sm">
            <span className="mb-1 block text-slate-600 dark:text-slate-300">Display name</span>
            <input
              className="w-full rounded-lg border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </label>

          <div className="text-sm text-slate-600 dark:text-slate-300">Currency: {currencyCode ?? 'Not set'}</div>

          <div>
            <span className="mb-2 block text-sm text-slate-600 dark:text-slate-300">Theme</span>
            <div className="flex gap-2">
              {(['light', 'dark', 'system'] as const).map((value) => (
                <button
                  key={value}
                  className={`rounded-lg border px-3 py-1 text-sm ${theme === value
                      ? 'border-slate-900 bg-slate-900 text-white'
                      : 'border-slate-300 text-slate-600 dark:border-slate-700 dark:text-slate-300'
                    }`}
                  onClick={() => setTheme(value)}
                >
                  {value}
                </button>
              ))}
            </div>
          </div>

          <div className="text-xs text-slate-400 dark:text-slate-500">Version 0.1.0</div>
        </div>

        <div className="mt-6 flex gap-2">
          <button className="rounded-lg bg-slate-900 px-4 py-2 text-sm text-white dark:bg-slate-100 dark:text-slate-900" onClick={saveProfile}>
            Save
          </button>
          <button className="rounded-lg border border-slate-300 px-4 py-2 text-sm dark:border-slate-700 dark:text-slate-200" onClick={onClose}>
            Close
          </button>
        </div>
      </div>
    </div>
  )
}
