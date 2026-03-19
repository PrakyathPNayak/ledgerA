import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { RecurringTransaction } from '@/types'

export const recurringQueryKeys = {
    all: ['recurring'] as const,
}

export function useRecurringTransactions() {
    return useQuery<{ data: RecurringTransaction[]; meta: { total: number } }>({
        queryKey: recurringQueryKeys.all,
        queryFn: () => api.get('/recurring'),
    })
}

export function useCreateRecurring() {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (payload: {
            account_id: string
            category_id: string
            subcategory_id?: string
            name: string
            amount: number
            notes?: string
            frequency: string
            start_date: string
            end_date?: string
        }) => api.post('/recurring', payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: recurringQueryKeys.all })
        },
    })
}

export function useUpdateRecurring() {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: ({ id, ...payload }: { id: string;[key: string]: unknown }) =>
            api.patch(`/recurring/${id}`, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: recurringQueryKeys.all })
        },
    })
}

export function useDeleteRecurring() {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (id: string) => api.delete(`/recurring/${id}`),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: recurringQueryKeys.all })
        },
    })
}
