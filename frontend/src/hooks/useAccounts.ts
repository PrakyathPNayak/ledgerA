import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import api from '@/lib/api'
import type { Account } from '@/types'

export const accountQueryKeys = {
  all: ['accounts'] as const,
  detail: (id: string) => ['accounts', id] as const,
}

/**
 * @description Fetches all accounts.
 * @returns Accounts query result.
 */
export function useAccounts() {
  return useQuery<Account[]>({
    queryKey: accountQueryKeys.all,
    queryFn: () => api.get('/accounts'),
    staleTime: 5 * 60 * 1000,
  })
}

/**
 * @description Creates a new account.
 * @returns Account creation mutation.
 */
export function useCreateAccount() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: { name: string; opening_balance: number }) => api.post('/accounts', payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: accountQueryKeys.all }),
  })
}

/**
 * @description Updates an account name.
 * @returns Account update mutation.
 */
export function useUpdateAccount() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: { id: string; name: string }) => api.patch(`/accounts/${payload.id}`, { name: payload.name }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: accountQueryKeys.all }),
  })
}

/**
 * @description Soft deletes an account.
 * @returns Account delete mutation.
 */
export function useDeleteAccount() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.delete(`/accounts/${id}`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: accountQueryKeys.all }),
  })
}
