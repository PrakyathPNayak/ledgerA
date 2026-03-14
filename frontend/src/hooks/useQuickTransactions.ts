import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import api from '@/lib/api'
import type { QuickTransaction } from '@/types'

export const quickTxQueryKeys = {
    all: ['quick-transactions'] as const,
}

/**
 * @description Fetches all quick transaction templates.
 */
export function useQuickTransactions() {
    return useQuery<QuickTransaction[]>({
        queryKey: quickTxQueryKeys.all,
        queryFn: () => api.get('/quick-transactions'),
        staleTime: 60 * 1000,
    })
}

/**
 * @description Creates a quick transaction template.
 */
export function useCreateQuickTransaction() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (payload: {
            label: string
            account_id?: string
            category_id?: string
            subcategory_id?: string
            amount?: number
            notes?: string
        }) => api.post('/quick-transactions', payload),
        onSuccess: () => queryClient.invalidateQueries({ queryKey: quickTxQueryKeys.all }),
    })
}

/**
 * @description Executes a quick transaction.
 */
export function useExecuteQuickTransaction() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (id: string) => api.post(`/quick-transactions/${id}/execute`),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: quickTxQueryKeys.all })
            queryClient.invalidateQueries({ queryKey: ['transactions'] })
        },
    })
}

/**
 * @description Reorders quick transactions by id list.
 */
export function useReorderQuickTransactions() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (ordered_ids: string[]) => api.patch('/quick-transactions/reorder', { ordered_ids }),
        onSuccess: () => queryClient.invalidateQueries({ queryKey: quickTxQueryKeys.all }),
    })
}
