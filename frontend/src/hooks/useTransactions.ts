import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import api from '@/lib/api'
import type { Transaction, TransactionFilters } from '@/types'

export const transactionQueryKeys = {
    all: ['transactions'] as const,
    list: (filters: TransactionFilters) => ['transactions', filters] as const,
}

/**
 * @description Fetches transactions with filter support.
 */
export function useTransactions(filters: TransactionFilters = {}) {
    return useQuery<Transaction[]>({
        queryKey: transactionQueryKeys.list(filters),
        queryFn: () => api.get('/transactions', { params: filters }),
        staleTime: 60 * 1000,
    })
}

/**
 * @description Creates a transaction and refreshes transaction data.
 */
export function useCreateTransaction() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (payload: {
            account_id: string
            category_id: string
            subcategory_id?: string
            name: string
            amount: number
            transaction_date: string
            notes?: string
            is_scheduled?: boolean
        }) => api.post('/transactions', payload),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: transactionQueryKeys.all })
            queryClient.invalidateQueries({ queryKey: ['accounts'] })
            queryClient.invalidateQueries({ queryKey: ['stats'] })
        },
    })
}

/**
 * @description Updates a transaction and refreshes all related caches.
 */
export function useUpdateTransaction() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ id, ...payload }: {
            id: string
            account_id?: string
            category_id?: string
            subcategory_id?: string
            name?: string
            amount?: number
            transaction_date?: string
            notes?: string
            is_scheduled?: boolean
        }) => api.patch(`/transactions/${id}`, payload),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: transactionQueryKeys.all })
            queryClient.invalidateQueries({ queryKey: ['accounts'] })
            queryClient.invalidateQueries({ queryKey: ['stats'] })
        },
    })
}

/**
 * @description Deletes a transaction and refreshes list cache.
 */
export function useDeleteTransaction() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (id: string) => api.delete(`/transactions/${id}`),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: transactionQueryKeys.all })
            queryClient.invalidateQueries({ queryKey: ['accounts'] })
            queryClient.invalidateQueries({ queryKey: ['stats'] })
        },
    })
}
