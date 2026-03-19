import { useMutation, useQueryClient } from '@tanstack/react-query'

import api from '@/lib/api'
import { transactionQueryKeys } from './useTransactions'

export function useTransfer() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (payload: {
            from_account_id: string
            to_account_id: string
            category_id: string
            subcategory_id?: string
            amount: number
            transaction_date: string
            name: string
            notes?: string
        }) => api.post('/transfers', payload),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: transactionQueryKeys.all })
            queryClient.invalidateQueries({ queryKey: ['accounts'] })
            queryClient.invalidateQueries({ queryKey: ['stats'] })
        },
    })
}
