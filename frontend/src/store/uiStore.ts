import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type ThemeMode = 'light' | 'dark' | 'system'

type UIState = {
    sidebarCollapsed: boolean
    theme: ThemeMode
    toggleSidebar: () => void
    setTheme: (theme: ThemeMode) => void
    applyTheme: () => void
}

function applyThemeToDom(theme: ThemeMode) {
    if (typeof window === 'undefined') {
        return
    }

    const root = document.documentElement
    root.classList.remove('light', 'dark')

    if (theme === 'system') {
        const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
        root.classList.add(systemTheme)
        return
    }

    root.classList.add(theme)
}

export const useUIStore = create<UIState>()(
    persist(
        (set, get) => ({
            sidebarCollapsed: false,
            theme: 'system',
            toggleSidebar: () => set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
            setTheme: (theme) => {
                set({ theme })
                applyThemeToDom(theme)
            },
            applyTheme: () => applyThemeToDom(get().theme),
        }),
        {
            name: 'ledgera-ui-store',
            onRehydrateStorage: () => (state) => {
                if (state) {
                    state.applyTheme()
                }
            },
        },
    ),
)
