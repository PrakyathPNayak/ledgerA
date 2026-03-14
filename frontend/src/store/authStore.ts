import type { User } from 'firebase/auth'
import { create } from 'zustand'

import { onAuthReady } from '@/lib/firebase'

type AuthState = {
    user: User | null
    isLoading: boolean
    isFirstTime: boolean
    setUser: (user: User | null) => void
    setLoading: (isLoading: boolean) => void
    setFirstTime: (isFirstTime: boolean) => void
    initialize: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
    user: null,
    isLoading: true,
    isFirstTime: false,
    setUser: (user) => set({ user }),
    setLoading: (isLoading) => set({ isLoading }),
    setFirstTime: (isFirstTime) => set({ isFirstTime }),
    initialize: () => {
        onAuthReady((user) => {
            set({ user, isLoading: false })
        })
    },
}))
