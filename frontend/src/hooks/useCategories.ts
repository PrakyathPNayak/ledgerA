import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import api from '@/lib/api'
import type { Category } from '@/types'

export const categoryQueryKeys = {
  all: ['categories'] as const,
}

/**
 * @description Fetches categories with subcategories.
 * @returns Categories query result.
 */
export function useCategories() {
  return useQuery<Category[]>({
    queryKey: categoryQueryKeys.all,
    queryFn: () => api.get('/categories'),
    staleTime: 5 * 60 * 1000,
  })
}

/**
 * @description Creates a category.
 * @returns Create category mutation.
 */
export function useCreateCategory() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: { name: string }) => api.post('/categories', payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: categoryQueryKeys.all }),
  })
}

/**
 * @description Creates a subcategory under category.
 * @param categoryId Parent category id.
 * @returns Create subcategory mutation.
 */
export function useCreateSubcategory(categoryId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: { name: string }) => api.post(`/categories/${categoryId}/subcategories`, payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: categoryQueryKeys.all }),
  })
}
