import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import api from '@/lib/api'
import type { User } from '@/types'

export const authQueryKeys = {
  me: ['auth', 'me'] as const,
}

/**
 * @description Fetches current user profile.
 * @returns User query result.
 */
export function useCurrentUser() {
  return useQuery<User>({
    queryKey: authQueryKeys.me,
    queryFn: () => api.get('/users/me'),
    staleTime: 5 * 60 * 1000,
  })
}

/**
 * @description Synchronizes authenticated firebase user with backend profile.
 * @returns Mutation for auth sync.
 */
export function useSyncUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: { firebase_token: string; display_name: string; email: string; currency_code?: string }) =>
      api.post('/auth/sync', payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: authQueryKeys.me })
    },
  })
}
