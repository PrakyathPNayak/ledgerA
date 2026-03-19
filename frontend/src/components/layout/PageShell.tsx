import { useEffect, useMemo, useState } from 'react'

import { Topbar } from '@/components/layout/Topbar'
import { Sidebar } from '@/components/layout/Sidebar'
import { SettingsSheet } from '@/components/layout/SettingsSheet'
import { Alert } from '@/components/ui/alert'
import { useAuthStore } from '@/store/authStore'
import { useUIStore } from '@/store/uiStore'
import { useCurrentUser } from '@/hooks/useAuth'
import api from '@/lib/api'

const CURRENCIES = ['INR', 'USD', 'EUR', 'GBP', 'JPY', 'AUD', 'CAD', 'CHF', 'SGD', 'AED']

/**
 * @description Page shell properties.
 */
interface PageShellProps {
  title: string
  children: React.ReactNode
}

export function PageShell({ title, children }: PageShellProps) {
  const [online, setOnline] = useState(navigator.onLine)
  const [settingsOpen, setSettingsOpen] = useState(false)
  const [selectedCurrency, setSelectedCurrency] = useState('INR')
  const [savingCurrency, setSavingCurrency] = useState(false)
  const { user, isFirstTime, setFirstTime } = useAuthStore()
  const { theme, applyTheme } = useUIStore()
  const { data: profile } = useCurrentUser()

  useEffect(() => {
    const onOnline = () => setOnline(true)
    const onOffline = () => setOnline(false)
    window.addEventListener('online', onOnline)
    window.addEventListener('offline', onOffline)
    return () => {
      window.removeEventListener('online', onOnline)
      window.removeEventListener('offline', onOffline)
    }
  }, [])

  useEffect(() => {
    applyTheme()

    if (theme !== 'system') {
      return
    }

    const media = window.matchMedia('(prefers-color-scheme: dark)')
    const onChange = () => applyTheme()
    media.addEventListener('change', onChange)
    return () => media.removeEventListener('change', onChange)
  }, [theme, applyTheme])

  const greetingName = useMemo(() => user?.displayName ?? user?.email ?? 'User', [user])
  const userInitial = useMemo(() => {
    const name = user?.displayName ?? user?.email ?? ''
    return name.charAt(0).toUpperCase() || 'U'
  }, [user])

  async function handleSaveCurrency() {
    setSavingCurrency(true)
    try {
      await api.patch('/users/me', { currency_code: selectedCurrency })
      setFirstTime(false)
    } finally {
      setSavingCurrency(false)
    }
  }

  return (
    <div className="flex min-h-screen bg-elevated text-foreground bg-app text-foreground">
      <Sidebar />
      <div className="flex min-h-screen flex-1 flex-col">
        {!online && (
          <div className="p-3">
            <Alert title="Offline Mode">
              You're offline. Showing cached data - some actions may not be available.
            </Alert>
          </div>
        )}
        <Topbar title={title} onOpenSettings={() => setSettingsOpen(true)} userInitial={userInitial} />
        <div className="flex-1 p-6">{children}</div>
      </div>

      <SettingsSheet
        open={settingsOpen}
        onClose={() => setSettingsOpen(false)}
        displayName={greetingName}
        currencyCode={profile?.currency_code}
      />

      {isFirstTime && (
        <div className="fixed inset-0 z-40 grid place-items-center bg-black/30">
          <div className="w-full max-w-md rounded-2xl bg-surface p-6 shadow-xl bg-surface">
            <h3 className="text-lg font-semibold text-foreground">Set Base Currency</h3>
            <p className="mt-2 text-sm text-muted">
              Choose your base currency once. It cannot be changed later.
            </p>
            <select
              value={selectedCurrency}
              onChange={(e) => setSelectedCurrency(e.target.value)}
              className="mt-3 w-full rounded-lg border border-border px-3 py-2 text-sm border-border bg-elevated text-foreground"
            >
              {CURRENCIES.map((code) => (
                <option key={code} value={code}>{code}</option>
              ))}
            </select>
            <button
              className="mt-4 rounded-lg bg-accent px-4 py-2 text-sm text-white disabled:opacity-60"
              onClick={handleSaveCurrency}
              disabled={savingCurrency}
            >
              {savingCurrency ? 'Saving...' : 'Continue'}
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
