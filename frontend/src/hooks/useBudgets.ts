import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Budget, BudgetProgress } from '@/types'

export const budgetQueryKeys = {
    all: ['budgets'] as const,
    progress: ['budgets', 'progress'] as const,
}

export function useBudgets() {
    return useQuery<{ data: Budget[]; meta: { total: number } }>({
        queryKey: budgetQueryKeys.all,
        queryFn: () => api.get('/budgets'),
    })
}

export function useBudgetProgress() {
    return useQuery<BudgetProgress[]>({
        queryKey: budgetQueryKeys.progress,
        queryFn: () => api.get('/budgets/progress'),
    })
}

export function useCreateBudget() {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (payload: { category_id: string; amount: number; period: string }) =>
            api.post('/budgets', payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: budgetQueryKeys.all })
            qc.invalidateQueries({ queryKey: budgetQueryKeys.progress })
        },
    })
}

export function useUpdateBudget() {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: ({ id, ...payload }: { id: string; amount?: number; is_active?: boolean }) =>
            api.patch(`/budgets/${id}`, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: budgetQueryKeys.all })
            qc.invalidateQueries({ queryKey: budgetQueryKeys.progress })
        },
    })
}

export function useDeleteBudget() {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (id: string) => api.delete(`/budgets/${id}`),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: budgetQueryKeys.all })
            qc.invalidateQueries({ queryKey: budgetQueryKeys.progress })
        },
    })
}
